# Padrões de Comunicação

## Visão Geral

O Monetics suporta dois modos de comunicação entre módulos:

1. **Comunicação Local**: Módulos compartilhando o mesmo processo
2. **Comunicação HTTP**: Módulos executando como serviços separados

Este design permite flexibilidade total de deployment - você pode executar tudo como um monólito ou dividir em microserviços independentes.

## Comunicação Local (Monólito)

### Quando Usar

- Desenvolvimento local
- Deployments pequenos a médios
- Prioridade em baixa latência
- Configuração simplificada
- Uso eficiente de recursos

### Como Funciona

```go
// Módulo Auth registra serviços localmente
container.Register("UserService", userService)

// Módulo Budget recupera serviços do mesmo container
userService := container.Get("UserService").(UserServiceInterface)
```

**Benefícios**:
- Sem latência de rede
- Sem serialização/desserialização
- Tipagem forte em tempo de compilação
- Transações compartilhadas possíveis
- Debugging simplificado

**Trade-offs**:
- Escalabilidade acoplada (todos os módulos escalam juntos)
- Falhas podem afetar todos os módulos
- Deploy de uma única unidade

## Comunicação HTTP (Microserviços)

### Quando Usar

- Escalabilidade independente de módulos
- Equipes distribuídas
- Deploys independentes necessários
- Isolamento de falhas crítico
- Sistemas de grande escala

### Como Funciona

```go
// Módulo Budget detecta que Auth está remoto
if container.Has("UserService") {
    // Local: usar serviço direto
    userService = container.Get("UserService")
} else {
    // Remoto: usar cliente HTTP
    userService = NewHTTPUserServiceClient(authServiceURL)
}
```

**Benefícios**:
- Escalabilidade independente
- Deploys independentes
- Isolamento de falhas
- Escolha flexível de tecnologia
- Limites claros entre serviços

**Trade-offs**:
- Latência de rede adicionada
- Necessita de retry logic e circuit breakers
- Gerenciamento de configuração mais complexo
- Tratamento de erros distribuídos

## Camada de Abstração de Comunicação

O sistema usa interfaces para abstrair o meio de comunicação:

```go
// Interface única para ambos os modos
type UserServiceInterface interface {
    GetUserByID(ctx context.Context, id string) (*User, error)
    ValidateUser(ctx context.Context, token string) (*User, error)
}

// Implementação local (em memória)
type LocalUserService struct {
    repo UserRepository
}

// Implementação remota (HTTP)
type HTTPUserServiceClient struct {
    baseURL    string
    httpClient *http.Client
    retrier    *Retrier
    breaker    *CircuitBreaker
}
```

### Benefícios da Abstração

1. **Transparência**: O código do módulo não precisa saber se está falando local ou remotamente
2. **Testabilidade**: Fácil criar mocks e stubs
3. **Flexibilidade**: Mudar de local para remoto não requer mudanças no código
4. **Migração Gradual**: Começar como monólito, migrar para microserviços incrementalmente

## Resiliência em Comunicação HTTP

### Retry Logic

```go
type Retrier struct {
    maxRetries     int
    initialBackoff time.Duration
    maxBackoff     time.Duration
}

func (r *Retrier) Do(operation func() error) error {
    backoff := r.initialBackoff
    for i := 0; i < r.maxRetries; i++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        if !isRetriable(err) {
            return err
        }
        
        time.Sleep(backoff)
        backoff = min(backoff*2, r.maxBackoff)
    }
    return fmt.Errorf("max retries exceeded")
}
```

**Configuração**:
- Max Retries: 3
- Initial Backoff: 100ms
- Max Backoff: 2s
- Estratégia: Exponential backoff

### Circuit Breaker

```go
type CircuitBreaker struct {
    failureThreshold int
    resetTimeout     time.Duration
    state            State // Closed, Open, HalfOpen
}

func (cb *CircuitBreaker) Call(operation func() error) error {
    if cb.state == Open {
        if time.Since(cb.lastFailure) > cb.resetTimeout {
            cb.state = HalfOpen
        } else {
            return ErrCircuitOpen
        }
    }
    
    err := operation()
    if err != nil {
        cb.recordFailure()
        return err
    }
    
    cb.recordSuccess()
    return nil
}
```

**Configuração**:
- Failure Threshold: 5 falhas consecutivas
- Reset Timeout: 30s
- Estados: Closed → Open → HalfOpen → Closed

### Timeout Configuration

```go
// Cliente HTTP com timeouts
httpClient := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   5 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout:   5 * time.Second,
        ResponseHeaderTimeout: 5 * time.Second,
    },
}
```

## Configuração

### Modo Local (config.yaml)

```yaml
modules:
  enabled:
    - auth
    - budget
  
# Sem configuração de URLs de serviços = modo local
```

### Modo Microserviços (config.yaml)

```yaml
modules:
  enabled:
    - budget  # Apenas budget habilitado neste serviço
  
  service_urls:
    auth: "http://auth-service:8081"
  
  http_client:
    timeout: 10s
    max_retries: 3
    circuit_breaker:
      failure_threshold: 5
      reset_timeout: 30s
```

### Variáveis de Ambiente

```bash
# Serviço Auth
export MODULE=auth
export PORT=8081

# Serviço Budget
export MODULE=budget
export PORT=8082
export AUTH_SERVICE_URL=http://localhost:8081
```

## Exemplos Práticos

### Exemplo 1: Desenvolvimento Local (Monólito)

```bash
# Iniciar todos os módulos em um processo
./bin/monetics --module=all

# Ou via Makefile
make run
```

**Resultado**:
- Auth e Budget no mesmo processo
- Comunicação local (sem HTTP)
- Porta única: 8080

### Exemplo 2: Microserviços Separados

```bash
# Terminal 1: Iniciar Auth
./bin/monetics --module=auth --port=8081

# Terminal 2: Iniciar Budget
export AUTH_SERVICE_URL=http://localhost:8081
./bin/monetics --module=budget --port=8082
```

**Resultado**:
- Auth na porta 8081
- Budget na porta 8082
- Budget → Auth via HTTP

### Exemplo 3: Docker Compose

```yaml
version: '3.8'

services:
  auth:
    image: monetics:latest
    command: ["--module=auth"]
    ports:
      - "8081:8080"
    environment:
      - MODULE=auth
  
  budget:
    image: monetics:latest
    command: ["--module=budget"]
    ports:
      - "8082:8080"
    environment:
      - MODULE=budget
      - AUTH_SERVICE_URL=http://auth:8080
    depends_on:
      - auth
```

## Monitoramento e Observabilidade

### Logs

```go
// Comunicação local
log.Info().
    Str("module", "budget").
    Str("dependency", "auth").
    Str("method", "GetUserByID").
    Str("mode", "local").
    Msg("calling user service")

// Comunicação HTTP
log.Info().
    Str("module", "budget").
    Str("dependency", "auth").
    Str("method", "GetUserByID").
    Str("mode", "http").
    Str("url", authServiceURL).
    Int("attempt", retry).
    Msg("calling user service")
```

### Métricas Recomendadas

- **Latência**: Tempo de resposta (p50, p95, p99)
- **Taxa de Erro**: Percentual de requisições com falha
- **Taxa de Retry**: Quantas requisições precisaram de retry
- **Circuit Breaker**: Estado e mudanças de estado
- **Throughput**: Requisições por segundo

## Boas Práticas

### 1. Comece Simples

- Inicie com modo monólito (local)
- Meça e profile antes de dividir
- Adicione complexidade apenas quando necessário

### 2. Design para Ambos os Modos

- Use interfaces desde o início
- Evite assumir comunicação local
- Sempre considere latência e falhas de rede

### 3. Implemente Resiliência

- Use timeouts em todas as chamadas
- Implemente retry com backoff exponencial
- Adicione circuit breakers para falhas persistentes
- Log todas as interações de serviços

### 4. Teste Ambos os Modos

- Testes unitários com mocks
- Testes de integração locais
- Testes de integração distribuídos
- Testes de caos (chaos engineering)

### 5. Monitore Tudo

- Log cada chamada de serviço
- Colete métricas de performance
- Configure alertas para circuit breakers
- Track taxa de retry

## Troubleshooting

### Problema: "Service not available"

**Causa**: Módulo dependente não está inicializado ou não está acessível

**Solução**:
1. Verifique se o módulo está habilitado no config
2. Para modo HTTP, verifique `service_urls` na configuração
3. Verifique logs de inicialização
4. Teste conectividade de rede (se distribuído)

### Problema: "Circuit breaker open"

**Causa**: Muitas falhas consecutivas para um serviço

**Solução**:
1. Verifique saúde do serviço dependente
2. Revise logs de erro do serviço
3. Aguarde reset timeout (30s por padrão)
4. Investigue causa raiz das falhas
5. 

### Problema: Alto tempo de resposta

**Causa**: Pode ser latência de rede ou processamento lento

**Solução**:
1. Verifique métricas de latência de rede
2. Profile o serviço dependente
3. Considere caching
4. Revise timeouts e retry logic

## Documentação Relacionada

- [Dependências entre Módulos](../module-dependencies.md)
- [Testando Dependências](../testing-dependencies.md)
- [Grafo de Dependências](../dependency-graph.md)
