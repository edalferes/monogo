# Gerenciamento de Dependências entre Módulos

## Visão Geral

O projeto Monetics usa um sistema de **Container de Injeção de Dependências** e **Registro de Módulos** para gerenciar dependências entre módulos. Isso garante que os módulos sejam inicializados na ordem correta e possam se comunicar localmente (em memória) ou remotamente (via HTTP).

## Arquitetura

### Componentes

1. **ModuleContainer** (`container.go`)
   - Gerencia registro e recuperação de serviços
   - Armazena configuração do módulo, banco de dados, logger, etc.
   - Registro de serviços thread-safe

2. **ModuleRegistry** (`registry.go`)
   - Gerencia registro de módulos com dependências
   - Resolve ordem de dependências automaticamente
   - Inicializa módulos na sequência correta

## Grafo Atual de Dependências dos Módulos

```
auth (sem dependências)
  ↓
budget (depende de auth)
  
test (sem dependências)
```

## Como Funciona

### Modo Monólito (`--module=all`)

Quando todos os módulos rodam juntos:

1. Módulo Auth inicializa primeiro
2. UserService do Auth é registrado no container
3. Módulo Budget inicializa e usa o UserService local
4. **Comunicação**: Em memória (sem HTTP)

### Modo Microserviços

#### Serviço Auth (`--module=auth`)

```bash
./bin/monetics --module=auth
# Roda na porta 8080 (ou porta configurada)
```

#### Serviço Budget (`--module=budget`)

```bash
export AUTH_SERVICE_URL=http://localhost:8080
./bin/monetics --module=budget
# Usa HTTP para comunicar com o serviço Auth
```

## Adicionando Novos Módulos

### Exemplo: Adicionando um Módulo "Notificações"

Digamos que você queira adicionar um módulo de notificações que depende tanto de Auth quanto de Budget.

#### Passo 1: Criar Estrutura do Módulo

```bash
internal/modules/notifications/
├── module.go
├── entities.go
├── seed.go
├── domain/
├── usecase/
├── adapters/
│   ├── http/
│   └── repository/
```

#### Passo 2: Implementar Funções do Módulo

```go
// internal/modules/notifications/module.go
package notifications

import (
"github.com/edalferes/monetics/pkg/logger"
"github.com/labstack/echo/v4"
"gorm.io/gorm"
)

// WireUp inicializa o módulo de notificações com serviços locais
func WireUp(group *echo.Group, db *gorm.DB, jwtSecret string, log logger.Logger) {
	// Inicializar com dependências locais
}

// WireUpWithHTTP inicializa o módulo de notificações com serviços HTTP
func WireUpWithHTTP(group *echo.Group, db *gorm.DB, jwtSecret string, log logger.Logger, 
authURL string, budgetURL string) {
	// Inicializar com dependências via cliente HTTP
}

// Entities retorna as entidades do banco de dados
func Entities() []interface{} {
	return []interface{}{
		// Suas entidades aqui
	}
}

// Seed popula dados iniciais
func Seed(db *gorm.DB) error {
	return nil
}
```

#### Passo 3: Registrar no Module Registry

Editar `internal/applications/api/registry.go`:

```go
func (r *ModuleRegistry) RegisterBuiltInModules() {
	// ... registro existente de auth e budget ...

	// Módulo Notifications - depende de Auth e Budget
	r.Register("notifications", []string{"auth", "budget"}, func(c *ModuleContainer) error {
cfg := c.GetConfig()
		db := c.GetDB()
		log := c.GetLogger()
		group := c.GetEchoGroup()

		authLocal := c.IsModuleEnabled("auth")
		budgetLocal := c.IsModuleEnabled("budget")

		if authLocal && budgetLocal {
			// Ambos os serviços locais
			log.Info().Msg("Usando serviços locais para módulo Notifications")
			notifications.WireUp(group, db, cfg.JWT.Secret, log)
		} else {
			// Usar HTTP para serviços remotos
			authURL := cfg.Modules.Auth.URL
			budgetURL := cfg.Modules.Budget.URL
			
			log.Info().
				Str("auth_url", authURL).
				Str("budget_url", budgetURL).
				Msg("Usando serviços remotos para módulo Notifications")
			
			notifications.WireUpWithHTTP(group, db, cfg.JWT.Secret, log, authURL, budgetURL)
		}

		return nil
	})
}
```

#### Passo 4: Adicionar ao parseModules

Editar `internal/applications/api/app.go`:

```go
func parseModules(input string) []string {
	if input == "" || input == "all" {
		return []string{"auth", "budget", "notifications"}  // Adicionar novo módulo
	}
	// ... resto da função
}
```

#### Passo 5: Adicionar Configuração

Editar `config.yaml`:

```yaml
modules:
  auth:
    url: "${AUTH_SERVICE_URL:}"
  budget:
    url: "${BUDGET_SERVICE_URL:}"
  notifications:
    url: "${NOTIFICATIONS_SERVICE_URL:}"
```

## Exemplos de Execução

### Desenvolvimento (Monólito)

```bash
./bin/monetics --module=all
# Todos os módulos usam comunicação local
```

### Produção (Microserviços)

Terminal 1 - Serviço Auth:

```bash
export PORT=8081
./bin/monetics --module=auth
```

Terminal 2 - Serviço Budget:

```bash
export PORT=8082
export AUTH_SERVICE_URL=http://localhost:8081
./bin/monetics --module=budget
```

Terminal 3 - Serviço Notifications:

```bash
export PORT=8083
export AUTH_SERVICE_URL=http://localhost:8081
export BUDGET_SERVICE_URL=http://localhost:8082
./bin/monetics --module=notifications
```

## Resolução de Dependências

O sistema automaticamente:

1. ✅ Detecta dependências circulares (retorna erro)
2. ✅ Inicializa módulos na ordem correta
3. ✅ Escolhe comunicação local vs HTTP automaticamente
4. ✅ Valida que dependências necessárias estão disponíveis
5. ✅ Fornece mensagens de erro claras para dependências ausentes

## Boas Práticas

### 1. Manter Dependências Mínimas

```go
// ✅ Bom: Dependências específicas
r.Register("notifications", []string{"auth"}, ...)

// ❌ Ruim: Muitas dependências
r.Register("notifications", []string{"auth", "budget", "user", "payment"}, ...)
```

### 2. Usar Interfaces para Serviços

```go
// Definir interface no pacote de contratos
type NotificationService interface {
    SendEmail(ctx context.Context, email string) error
}

// Registrar no container
c.Register("notifications.NotificationService", service)
```

### 3. Falhar Rápido em Dependências Ausentes

```go
if !c.IsModuleEnabled("auth") && cfg.Modules.Auth.URL == "" {
    return fmt.Errorf("notifications requer serviço auth")
}
```

### 4. Registrar Tipo de Comunicação em Log

```go
if useLocal {
    log.Info().Msg("Usando serviço Auth local")
} else {
    log.Info().Str("url", authURL).Msg("Usando serviço Auth remoto")
}
```

## Resolução de Problemas

### "Module X not registered"

**Causa**: Módulo não adicionado ao `RegisterBuiltInModules()`  
**Solução**: Adicionar registro do módulo em `registry.go`

### "Dependency module Y not enabled"

**Causa**: Rodando módulo sem sua dependência e sem URL remota configurada  
**Solução**: Habilitar módulo de dependência ou configurar URL remota

### Dependências Circulares

**Causa**: Módulo A depende de B, e B depende de A  
**Solução**: Refatorar para remover dependência circular ou criar serviço compartilhado

## Melhorias Futuras

- [ ] Hot-reload de módulos sem restart
- [ ] Health checks para dependências remotas
- [ ] Service discovery automático (Consul/etcd)
- [ ] Visualização do grafo de dependências
- [ ] Suporte a versionamento de módulos
