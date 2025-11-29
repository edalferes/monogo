# API de Health Check

## Visão Geral

A API de Health Check fornece endpoints para verificar a saúde da aplicação e suas dependências.

## Endpoints

### GET /health

Retorna o status de saúde básico da aplicação.

**Resposta de Sucesso (200 OK)**:

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Resposta de Falha (503 Service Unavailable)**:

```json
{
  "status": "unhealthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "details": {
    "database": "connection failed"
  }
}
```

### GET /health/ready

Verifica se a aplicação está pronta para receber tráfego (readiness probe).

**Resposta de Sucesso (200 OK)**:

```json
{
  "status": "ready",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Resposta de Falha (503 Service Unavailable)**:

```json
{
  "status": "not_ready",
  "timestamp": "2024-01-15T10:30:00Z",
  "reason": "database connection not established"
}
```

### GET /health/live

Verifica se a aplicação está viva (liveness probe).

**Resposta de Sucesso (200 OK)**:

```json
{
  "status": "alive",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Health Checks Implementados

### 1. Database Health Check

Verifica conectividade com o banco de dados.

**Critérios**:
- Conexão estabelecida
- Ping bem-sucedido
- Pool de conexões disponível

**Falha se**:
- Não consegue conectar
- Timeout na query
- Pool esgotado

### 2. Module Health Check

Verifica se todos os módulos habilitados estão funcionando.

**Critérios**:
- Todos os módulos inicializados
- Dependências resolvidas
- Serviços registrados

**Falha se**:
- Módulo não inicializou
- Dependência faltando
- Falha na inicialização

## Configuração

### Timeouts

```yaml
health:
  timeout: 5s
  database_timeout: 3s
```

### Intervalos de Verificação (Kubernetes)

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

## Uso em Diferentes Ambientes

### Desenvolvimento Local

```bash
# Verificar saúde
curl http://localhost:8080/health

# Verificar readiness
curl http://localhost:8080/health/ready

# Verificar liveness
curl http://localhost:8080/health/live
```

### Docker

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
```

### Docker Compose

```yaml
services:
  api:
    image: monetics:latest
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"\]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 40s
```

## Monitoramento

### Métricas Recomendadas

- **health_check_duration**: Tempo de execução do health check
- **health_check_failures**: Número de falhas consecutivas
- **dependency_status**: Status de cada dependência

### Alertas

Configure alertas para:

- Health check falhando por mais de 2 minutos
- Readiness probe falhando
- Database health check falhando

## Boas Práticas

1. **Use Readiness para Rolling Updates**: Previne tráfego antes da aplicação estar pronta
2. **Use Liveness para Detectar Deadlocks**: Reinicia aplicação se ela parar de responder
3. **Configure Timeouts Apropriados**: Evite falsos positivos
4. **Monitore Métricas**: Track tendências de saúde ao longo do tempo

## Troubleshooting

### Health Check Falhando

1. Verifique logs da aplicação
2. Teste conexão com banco de dados manualmente
3. Verifique status dos módulos
4. Valide configuração

### Readiness Probe Falhando

1. Verifique se aplicação terminou inicialização
2. Valide todas as dependências estão disponíveis
3. Verifique logs de startup

## Documentação Relacionada

- [Visão Geral de Arquitetura](../architecture/overview.md)
- [Padrões de Comunicação](../architecture/communication.md)
