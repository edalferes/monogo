# Testando Dependências entre Módulos

Este guia demonstra como testar o sistema de injeção de dependências com diferentes configurações de módulos.

## Cenários de Teste

### Cenário 1: Modo Monólito (Todos os Módulos Locais)

**Configuração**: Todos os módulos rodam em um único processo

```bash
./bin/monetics --module=all
```

**Comportamento Esperado**:

- Módulo Auth inicializa primeiro
- Módulo Budget inicializa em segundo (espera pelo auth)
- Budget usa UserService do Auth **local** (em memória)
- Sem chamadas HTTP entre módulos

**Verificação**:

```bash
# Verificar nos logs por:
# "Initializing Auth module..."
# "Auth module started successfully"
# "Initializing Budget module..."
# "Using local Auth service" (ou similar)
```

### Cenário 2: Serviço Apenas Auth

**Configuração**: Executar apenas o módulo Auth

```bash
export PORT=8081
./bin/monetics --module=auth
```

**Comportamento Esperado**:

- Apenas módulo Auth inicializa
- Endpoints: `/v1/auth/*`
- Escuta na porta 8081

**Verificação**:

```bash
curl http://localhost:8081/health
# Deve retornar health check apenas para o módulo auth
```

### Cenário 3: Serviço Budget com Auth Remoto

**Configuração**: Módulo Budget usando HTTP para comunicar com Auth

**Terminal 1 - Serviço Auth**:

```bash
export PORT=8081
./bin/monetics --module=auth
```

**Terminal 2 - Serviço Budget**:

```bash
export PORT=8082
export AUTH_SERVICE_URL=http://localhost:8081
./bin/monetics --module=budget
```

**Comportamento Esperado**:

- Módulo Budget inicia
- Detecta que Auth NÃO está local
- Lê `AUTH_SERVICE_URL` da configuração
- Usa cliente HTTP para comunicar com Auth
- Cliente HTTP usa retry + circuit breaker

**Verificação**:

```bash
# No Terminal 2, verificar nos logs:
# "Using remote Auth service"
# "auth_url": "http://localhost:8081"

# Testar endpoints do budget (deve funcionar)
curl http://localhost:8082/v1/budget/accounts

# Parar Terminal 1 (matar serviço Auth)
# Tentar novamente:
curl http://localhost:8082/v1/budget/accounts
# Deve ver tentativas de retry e ativação do circuit breaker
```

### Cenário 4: Erro de Dependência Ausente

**Configuração**: Tentar rodar Budget sem serviço Auth disponível

```bash
# SEM AUTH_SERVICE_URL configurada
./bin/monetics --module=budget
```

**Comportamento Esperado**:

- Módulo Budget tenta inicializar
- Detecta que Auth NÃO está local
- `AUTH_SERVICE_URL` está vazio
- **Deve falhar com erro**: "budget requires auth service"

**Verificação**:

```bash
# Deve ver erro nos logs e aplicação encerra
```

### Cenário 5: Múltiplos Serviços (Auth + Budget)

**Configuração**: Executar Auth e Budget localmente (mas separadamente)

**Terminal 1 - Auth**:

```bash
export PORT=8081
./bin/monetics --module=auth
```

**Terminal 2 - Budget (com Auth habilitado localmente também)**:

```bash
export PORT=8082
./bin/monetics --module=auth,budget
```

**Comportamento Esperado**:

- Terminal 2 executa Auth E Budget
- Budget usa Auth **local** (em memória)
- Terminal 1 e Terminal 2 têm instâncias independentes do Auth

**Nota**: Este cenário mostra como você pode executar o mesmo módulo em múltiplos processos para balanceamento de carga.

## Testes de Integração

### Testar Integração Budget → Auth

**Endpoint**: Criar uma transação de orçamento (requer autenticação de usuário)

**Script de Teste**:

```bash
#!/bin/bash

# 1. Iniciar serviço Auth
export PORT=8081
./bin/monetics --module=auth &
AUTH_PID=$!
sleep 2

# 2. Iniciar serviço Budget com Auth remoto
export PORT=8082
export AUTH_SERVICE_URL=http://localhost:8081
./bin/monetics --module=budget &
BUDGET_PID=$!
sleep 2

# 3. Registrar um usuário (serviço Auth)
curl -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123!@#"
  }'

# 4. Login para obter token
TOKEN=$(curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123!@#"
  }' | jq -r '.token')

echo "Token: $TOKEN"

# 5. Criar conta de orçamento (serviço Budget, autenticado via Auth)
curl -X POST http://localhost:8082/v1/budget/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Conta Corrente",
    "type": "checking",
    "currency": "BRL"
  }'

# 6. Limpeza
kill $AUTH_PID $BUDGET_PID
```

**Esperado**: Serviço Budget valida token chamando o serviço Auth (HTTP)

## Testes de Performance

### Medir Overhead: Local vs HTTP

**Teste**: Comparar latência de operações do Budget com Auth local vs remoto

**Configuração**:

```bash
# Teste 1: Monólito (local)
./bin/monetics --module=all
ab -n 1000 -c 10 http://localhost:8080/v1/budget/accounts

# Teste 2: Microserviços (HTTP)
# Terminal 1:
./bin/monetics --module=auth
# Terminal 2:
export AUTH_SERVICE_URL=http://localhost:8080
./bin/monetics --module=budget
ab -n 1000 -c 10 http://localhost:8080/v1/budget/accounts
```

**Métricas para Comparar**:

- Tempo médio de resposta
- Latência do percentil 95
- Requisições por segundo
- Uso de memória

**Esperado**: Modo monólito deve ser mais rápido (sem overhead de rede)

## Testes do Circuit Breaker

### Teste 1: Circuit Abre Após Falhas

**Configuração**:

```bash
# Terminal 1: Iniciar Auth
./bin/monetics --module=auth --port=8081

# Terminal 2: Iniciar Budget apontando para Auth
export AUTH_SERVICE_URL=http://localhost:8081
./bin/monetics --module=budget --port=8082
```

**Passos do Teste**:

1. Fazer requisições bem-sucedidas para verificar que funciona
2. **Parar serviço Auth** (matar Terminal 1)
3. Fazer 5+ requisições para endpoints do Budget que requerem Auth
4. **Verificar logs**: Circuit breaker deve abrir após máximo de falhas (5)
5. Continuar fazendo requisições → Deve falhar rápido (sem tentativas de retry)
6. Reiniciar serviço Auth (Terminal 1)
7. Aguardar 30 segundos (timeout do circuit breaker)
8. Fazer requisição → Circuit breaker vai para half-open
9. Requisições bem-sucedidas → Circuit breaker fecha

**Logs de Verificação**:

```
[WARN] Chamada ao serviço Auth falhou, tentando novamente (tentativa 1/3)
[WARN] Chamada ao serviço Auth falhou, tentando novamente (tentativa 2/3)
[WARN] Chamada ao serviço Auth falhou, tentando novamente (tentativa 3/3)
[ERROR] Circuit breaker aberto para serviço auth
[INFO] Circuit breaker half-open, testando conexão
[INFO] Circuit breaker fechado, serviço recuperado
```

### Teste 2: Backoff Exponencial

**Configuração**: Mesmo do Teste 1

**Observação**: Observar timing de retry nos logs

**Padrão Esperado**:

- Tentativa 1: Imediato
- Tentativa 2: ~100ms de delay
- Tentativa 3: ~200ms de delay (com jitter)

**Verificação**:

```bash
# Verificar timestamps nos logs:
# 2024-01-01T10:00:00.000Z [WARN] Tentativa de retry 1
# 2024-01-01T10:00:00.100Z [WARN] Tentativa de retry 2  # ~100ms depois
# 2024-01-01T10:00:00.300Z [WARN] Tentativa de retry 3  # ~200ms depois
```

## Dicas de Depuração

### Habilitar Logging Verboso

```bash
export LOG_LEVEL=debug
./bin/monetics --module=all
```

### Verificar Registro de Módulos

Adicionar logging em `registry.go`:

```go
func (r *ModuleRegistry) Initialize() error {
    log.Debug().Interface("enabled_modules", r.container.GetModules()).Msg("Iniciando inicialização")
    // ... resto do código
}
```

### Verificar Registro de Serviços

Adicionar em `container.go`:

```go
func (c *ModuleContainer) Register(name string, service interface{}) {
    c.logger.Debug().Str("service", name).Msg("Registrando serviço")
    // ... resto do código
}
```

### Inspecionar Requisições HTTP

Usar middleware para registrar todas as chamadas HTTP entre serviços:

```bash
# Em httpclient, adicionar logging de requisições
log.Debug().
    Str("method", req.Method).
    Str("url", req.URL.String()).
    Msg("Requisição HTTP saindo")
```

## Problemas Comuns

### "panic: runtime error: invalid memory address"

**Causa**: Tentando usar um serviço do container que não foi registrado

**Correção**: Verificar ordem de inicialização de módulos e garantir que todos os serviços estão registrados

### Erros "connection refused"

**Causa**: URL do serviço remoto aponta para host/porta errados

**Correção**: Verificar variáveis de ambiente e configuração:

```bash
echo $AUTH_SERVICE_URL
# Deve corresponder ao endereço real do serviço Auth
```

### Budget inicia mas não consegue autenticar usuários

**Causa**: Descompasso no JWT secret entre serviços Auth e Budget

**Correção**: Garantir que ambos os serviços usam o mesmo `config.yaml` ou variável de ambiente

## Próximos Passos

1. ✅ Documentar sistema de dependência de módulos
2. ✅ Criar guia de testes
3. ⬜ Adicionar testes de integração automatizados
4. ⬜ Implementar endpoint de métricas do circuit breaker
5. ⬜ Adicionar ferramenta de visualização do grafo de dependências
6. ⬜ Criar benchmarks de performance
