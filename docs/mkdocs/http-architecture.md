# Arquitetura HTTP - Inspirada no Grafana Loki

## Visão Geral

O Monetics implementa uma arquitetura modular inspirada no **Grafana Loki**, permitindo que os módulos sejam executados de forma independente ou como parte de um monólito, com comunicação HTTP entre serviços.

## Módulos e Serviços

### Módulos Disponíveis

- **auth**: Autenticação, autorização, usuários, roles e permissions
- **budget**: Contas, categorias, transações, orçamentos e relatórios

### Modos de Execução

```bash
# Monólito - todos os módulos
./bin/monetics --module=all

# Microservice - apenas auth
./bin/monetics --module=auth

# Microservice - apenas budget
./bin/monetics --module=budget

# Múltiplos módulos
./bin/monetics --module=auth,budget
```

## Comunicação HTTP

### Configuração

A comunicação entre módulos é controlada via `config.yaml`:

```yaml
modules:
  auth:
    url: ""  # Local (monólito)
    # url: "http://localhost:8081"  # HTTP remoto
```

### Client HTTP (`pkg/httpclient`)

Cliente HTTP reutilizável com recursos avançados de resiliência:

- **Connection pooling**: 100 conexões idle por padrão
- **Timeouts**: 10s por requisição (configurável)
- **Idle timeout**: 90s para conexões idle
- **Retry logic**: Retry automático com backoff exponencial
  - 3 tentativas por padrão
  - Backoff inicial: 100ms
  - Backoff máximo: 5s
  - Multiplicador: 2.0
  - Jitter aleatório (0-25%) para evitar thundering herd
- **Circuit breaker**: Proteção contra cascata de falhas
  - Estados: Closed → Open → Half-Open
  - 5 falhas consecutivas abrem o circuito
  - 30s timeout antes de tentar Half-Open
  - 3 requisições de teste no estado Half-Open

#### Uso do Client

```go
import "github.com/edalferes/monetics/pkg/httpclient"

// Criar client
cfg := httpclient.DefaultConfig("http://auth-service:8080")
client := httpclient.NewClient(cfg)

// GET request
resp, err := client.Get(ctx, "/v1/auth/users/1")
if err != nil {
    return err
}

// POST request
body := map[string]interface{}{"name": "John"}
resp, err := client.Post(ctx, "/v1/auth/users", body)

// PUT request
resp, err := client.Put(ctx, "/v1/auth/users/1", body)

// DELETE request
resp, err := client.Delete(ctx, "/v1/auth/users/1")
```

### UserServiceHTTP

Exemplo de implementação HTTP para comunicação entre módulos:

```go
// internal/modules/auth/service/user_service_http.go
type UserServiceHTTP struct {
    client *httpclient.Client
}

func NewUserServiceHTTP(baseURL string) contracts.UserService {
    cfg := httpclient.DefaultConfig(baseURL)
    return &UserServiceHTTP{
        client: httpclient.NewClient(cfg),
    }
}

func (s *UserServiceHTTP) GetUserByID(ctx context.Context, userID uint) (*contracts.UserInfo, error) {
    path := fmt.Sprintf("/v1/auth/users/%d", userID)
    resp, err := s.client.Get(ctx, path)
    if err != nil {
        return nil, err
    }
    
    var userInfo contracts.UserInfo
    json.Unmarshal(resp.Data, &userInfo)
    return &userInfo, nil
}
```

## Health Checks

Endpoints inspirados no Loki para Kubernetes:

### `/health`

Retorna status geral do serviço e módulos ativos:

```json
{
  "status": "healthy",
  "service": "monetics",
  "modules": ["auth", "budget"],
  "details": {
    "database": "connected"
  }
}
```

### `/ready`

Readiness probe - verifica se o serviço está pronto para receber tráfego:

```json
{
  "status": "ready"
}
```

Retorna 503 se:
- Banco de dados não estiver conectado
- Ping do banco falhar

### `/live`

Liveness probe - verifica se o processo está vivo:

```json
{
  "status": "alive"
}
```

### Uso no Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: monetics-auth
spec:
  template:
    spec:
      containers:
      - name: auth
        image: monetics:latest
        args: ["--module=auth"]
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /live
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

## Deployment Patterns

### Monólito (Desenvolvimento)

```yaml
# docker-compose.yml
services:
  monetics:
    build: .
    command: ["--module=all"]
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://...
```

### Microservices (Produção)

```yaml
# docker-compose.yml
services:
  auth:
    build: .
    command: ["--module=auth"]
    ports:
      - "8081:8080"
    
  budget:
    build: .
    command: ["--module=budget"]
    ports:
      - "8082:8080"
    environment:
      - AUTH_SERVICE_URL=http://auth:8080
```

### Kubernetes

```yaml
# auth-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: auth
        image: monetics:latest
        args: ["--module=auth"]

---
# budget-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: budget
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: budget
        image: monetics:latest
        args: ["--module=budget"]
        env:
        - name: AUTH_SERVICE_URL
          value: "http://auth-service:8080"
```

## Comparação com Loki

| Aspecto | Loki | Monetics |
|---------|------|----------|
| **Module Selection** | `-target=all,distributor,ingester` | `--module=all,auth,budget` |
| **Communication** | gRPC + HTTP | HTTP (com plano para gRPC) |
| **Health Checks** | `/ready`, `/health` | `/ready`, `/health`, `/live` |
| **Service Discovery** | Ring (Memberlist/Consul) | Config-based URLs |
| **Deployment** | Kubernetes, Docker | Docker, Kubernetes |

## Próximos Passos

### Implementado

1. ✅ **Retry Logic**: Retry automático com backoff exponencial e jitter
2. ✅ **Circuit Breaker**: Proteção contra cascata de falhas com estados (Closed, Open, Half-Open)
3. ✅ **Connection Pooling**: Gerenciamento otimizado de conexões HTTP
4. ✅ **Health Checks**: Endpoints `/health`, `/ready`, `/live` para Kubernetes

### Planejado

1. **gRPC Support**: Adicionar gRPC para performance em comunicação entre serviços
2. **Distributed Tracing**: OpenTelemetry para rastreamento distribuído
3. **Service Mesh**: Integração com Istio/Linkerd
4. **Rate Limiting**: Proteção contra sobrecarga de requisições

### Quando usar gRPC vs HTTP

**Use HTTP quando:**
- Desenvolvimento local
- Debugging e troubleshooting
- Integrações externas
- APIs públicas

**Use gRPC quando:**
- Alta performance necessária
- Volume alto de chamadas
- Necessidade de streaming
- Comunicação interna entre serviços

## Referências

- [Grafana Loki Architecture](https://grafana.com/docs/loki/latest/architecture/)
- [Kubernetes Health Checks](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [Go HTTP Client Best Practices](https://golang.org/pkg/net/http/)
