# ![Logz Banner](assets/top_banner.png)

---

**Uma ferramenta de gerenciamento de logs e métricas avançada, com suporte nativo à integração com Prometheus, notificações dinâmicas e uma CLI poderosa.**

---

## **Índice**

1. [Sobre o Projeto](#sobre-o-projeto)
2. [Destaques](#destaques)
3. [Instalação](#instalação)
4. [Uso](#uso)
    - [Interface de Linha de Comando](#interface-de-linha-de-comando)
    - [Configuração](#configuração)
5. [Integração com Prometheus](#integração-com-prometheus)
6. [Roadmap](#roadmap)
7. [Contribuições](#contribuições)
8. [Contato](#contato)

---

## **Sobre o Projeto**

O Logz é uma solução poderosa e flexível para gerenciar logs e métricas em sistemas modernos. Construído com **Go**, oferece suporte extensivo a múltiplos métodos de notificação, incluindo **HTTP Webhooks**, **streaming em tempo real via WebSocket**, e **integração nativa com Prometheus**, além de integração fluida com sistemas de monitoramento para observabilidade avançada.

O objetivo é fornecer uma ferramenta robusta, altamente configurável e escalável para desenvolvedores, equipes DevOps e arquitetos de software que precisam de uma abordagem centralizada para gerenciar logs e métricas.

**O que torna o Logz especial?**

- 🚀 **Testado em Produção**: Lida com sucesso com 500+ operações de logging simultâneas
- 🌊 **Streaming em Tempo Real**: Suporte WebSocket para monitoramento de logs ao vivo
- 📊 **Prometheus Nativo**: Coleta de métricas integrada e exposição de endpoint HTTP
- 🎯 **Zero Mocking**: Todos os testes usam servidores HTTP reais e conexões WebSocket
- 🧪 **Testes Abrangentes**: 11+ testes de integração cobrindo todos os recursos
- 🔗 **Integração Simples**: API simples que funciona com aplicações Go existentes

**Principais Benefícios:**

- 💡 **Fácil de usar**: Configure e gerencie logs sem esforço com APIs intuitivas
- 🌐 **Integração Perfeita**: Substituto direto para logging padrão com recursos avançados
- 🔧 **Extensível**: Adicione novos notificadores e serviços conforme necessário
- 🏭 **Pronto para Enterprise**: Logging thread-safe e de alta performance para cargas de produção

---

## **Destaques**

✨ **Sistema de Logging Avançado**:

- **Múltiplos Níveis de Log**: DEBUG, INFO, WARN, ERROR com formatação emoji bonita
- **Logging Concorrente Thread-Safe**: Lida com milhares de operações de log simultâneas
- **Padrão Builder LogEntry**: Construção flexível e intuitiva de entradas de log
- **Múltiplos Formatadores**: Logs estruturados JSON e saída de texto colorida

🌐 **Notificações em Tempo Real**:

- **HTTP Webhooks**: Envie logs para serviços externos via HTTP POST
- **Integração WebSocket**: Streaming de logs em tempo real para clientes conectados
- **Gerenciamento Dinâmico de Notificadores**: Adicione/remova notificadores dinamicamente
- **Suporte de Autenticação**: Entrega segura de webhook com autenticação baseada em token

📊 **Integração com Prometheus**:

- **Suporte Nativo de Métricas**: Tipos Counter, Gauge e métricas customizadas
- **Endpoint HTTP de Métricas**: Endpoint padrão `/metrics` para scraping do Prometheus
- **Persistência de Métricas**: Salvamento/carregamento automático de métricas para/do disco
- **Validação de Métricas**: Aplica convenções de nomenclatura do Prometheus
- **Suporte Whitelist**: Controle quais métricas são exportadas

� **Experiência do Desenvolvedor**:

- **Zero Mocking Necessário**: Todos os testes usam servidores HTTP reais e conexões WebSocket
- **Suíte de Testes Abrangente**: 11+ testes de integração cobrindo todos os recursos
- **API Simples**: Interfaces fáceis de usar para integração rápida
- **Logging Consciente de Contexto**: Suporte rico de metadata para melhor debugging

� **Pronto para Produção**:

- **Alta Performance**: Testado com 500+ operações simultâneas
- **Eficiente em Memória**: Uso otimizado de mutex e gerenciamento de recursos
- **Tratamento de Erros**: Mecanismos robustos de tratamento de erros e recuperação
- **Configurável**: Opções de configuração flexíveis para diferentes ambientes

---

## **Instalação**

Requisitos:

- **Go** versão 1.23 ou superior.
- Prometheus (opcional para monitoramento avançado).

```bash
# Clone este repositório
git clone https://github.com/kubex-ecosystem/logz.git

# Navegue até o diretório do projeto
cd logz

# Construa o binário usando make
make build

# Instale o binário usando make
make install

# (Opcional) Adicione o binário ao PATH para usá-lo globalmente
export PATH=$PATH:$(pwd)
```

---

## **Uso**

## Exemplos de Logging Básico

### Logging Simples

```go
package main

import "github.com/kubex-ecosystem/logz/logger"

func main() {
    // Criar uma nova instância do logger
    log := logger.NewLogger("meu-app")

    // Logging básico com formatação emoji
    log.InfoCtx("Aplicação iniciada com sucesso", nil)
    log.WarnCtx("Esta é uma mensagem de aviso", nil)
    log.ErrorCtx("Algo deu errado", nil)
}
```

**Saída:**

```plaintext
 [INFO]  ℹ️  - Aplicação iniciada com sucesso
 [WARN]  ⚠️  - Esta é uma mensagem de aviso
 [ERROR] ❌  - Algo deu errado
```

### Logging com Metadata

```go
package main

import "github.com/kubex-ecosystem/logz/logger"

func main() {
    log := logger.NewLogger("meu-servico")

    // Definir metadata global
    log.SetMetadata("service", "user-api")
    log.SetMetadata("version", "1.2.3")

    // Log com contexto adicional
    log.InfoCtx("Login de usuário bem-sucedido", map[string]interface{}{
        "user_id":    12345,
        "ip_address": "192.168.1.100",
        "duration":   "250ms",
    })
}
```

### Usando o Padrão Builder LogEntry

```go
package main

import (
    "github.com/kubex-ecosystem/logz/internal/core"
)

func main() {
    // Criar entradas de log estruturadas
    entry := core.NewLogEntry().
        WithLevel(core.INFO).
        WithMessage("Pagamento processado com sucesso").
        AddMetadata("transaction_id", "txn_123456").
        AddMetadata("amount", 99.99).
        AddMetadata("currency", "BRL").
        SetSeverity(1)

    // Usar a entrada com formatadores
    formatter := core.NewJSONFormatter()
    output := formatter.Format(entry)
    fmt.Println(output)
}
```

## Recursos Avançados

### Notificações HTTP Webhook

```go
package main

import (
    "github.com/kubex-ecosystem/logz/logger"
    "github.com/kubex-ecosystem/logz/internal/core"
)

func main() {
    // Criar logger com notificador HTTP
    log := logger.NewLogger("webhook-app")

    // Adicionar notificador webhook HTTP
    httpNotifier := core.NewHTTPNotifier(
        "https://hooks.slack.com/services/SEU/WEBHOOK/URL",
        "seu-token-auth",
    )

    // Criar entrada de log
    entry := core.NewLogEntry().
        WithLevel(core.ERROR).
        WithMessage("Erro crítico do sistema detectado").
        AddMetadata("severity", "high").
        AddMetadata("component", "database")

    // Enviar notificação
    err := httpNotifier.Notify(entry)
    if err != nil {
        log.ErrorCtx("Falha ao enviar webhook", map[string]interface{}{
            "error": err.Error(),
        })
    }
}
```

### Logging WebSocket em Tempo Real

```go
package main

import (
    "net/http"
    "github.com/kubex-ecosystem/logz/internal/core"
    "github.com/gorilla/websocket"
)

func main() {
    // Criar notificador WebSocket
    wsNotifier := core.NewWebSocketNotifier("ws://localhost:8080/logs", nil)

    // Criar e enviar entrada de log
    entry := core.NewLogEntry().
        WithLevel(core.INFO).
        WithMessage("Atualização de log em tempo real").
        AddMetadata("timestamp", time.Now().Unix()).
        AddMetadata("event", "user_action")

    err := wsNotifier.Notify(entry)
    if err != nil {
        fmt.Printf("Notificação WebSocket falhou: %v\n", err)
    }
}
```

### Integração de Métricas Prometheus

```go
package main

import "github.com/kubex-ecosystem/logz/internal/core"

func main() {
    // Obter instância do gerenciador Prometheus
    prometheus := core.GetPrometheusManager()

    // Adicionar várias métricas
    prometheus.AddMetric("http_requests_total", 100, map[string]string{
        "method": "GET",
        "status": "200",
    })

    prometheus.AddMetric("response_time_seconds", 0.045, map[string]string{
        "endpoint": "/api/users",
    })

    // Incrementar contador
    prometheus.IncrementMetric("api_calls_total", 1)

    // Iniciar servidor HTTP para expor endpoint /metrics
    prometheus.StartHTTPServer(":2112")
}
```

### Logging Concorrente (Pronto para Produção)

```go
package main

import (
    "sync"
    "github.com/kubex-ecosystem/logz/logger"
)

func main() {
    log := logger.NewLogger("concurrent-app")
    var wg sync.WaitGroup

    // Simular logging de alto tráfego
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            log.InfoCtx("Processando requisição", map[string]interface{}{
                "request_id": id,
                "worker":     "goroutine",
                "status":     "processing",
            })
        }(i)
    }

    wg.Wait()
    log.InfoCtx("Todas as requisições processadas", nil)
}
```

## Exemplos de Uso em Produção

### API Web com Logging e Métricas

```go
package main

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/kubex-ecosystem/logz/logger"
    "github.com/kubex-ecosystem/logz/internal/core"
)

func main() {
    // Inicializar logger e métricas
    log := logger.NewLogger("api-server")
    prometheus := core.GetPrometheusManager()

    // Definir metadata global
    log.SetMetadata("service", "user-api")
    log.SetMetadata("version", "1.0.0")

    // Iniciar servidor de métricas
    go prometheus.StartHTTPServer(":2112")

    // Configurar roteador Gin
    r := gin.Default()

    // Middleware para logging e métricas
    r.Use(func(c *gin.Context) {
        start := time.Now()

        // Processar requisição
        c.Next()

        // Fazer log dos detalhes da requisição
        duration := time.Since(start)
        log.InfoCtx("Requisição HTTP", map[string]interface{}{
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "status":     c.Writer.Status(),
            "duration":   duration.String(),
            "client_ip":  c.ClientIP(),
            "user_agent": c.Request.UserAgent(),
        })

        // Atualizar métricas
        prometheus.AddMetric("http_requests_total", 1, map[string]string{
            "method": c.Request.Method,
            "status": fmt.Sprintf("%d", c.Writer.Status()),
        })

        prometheus.AddMetric("http_request_duration_seconds",
            duration.Seconds(), map[string]string{
            "endpoint": c.Request.URL.Path,
        })
    })

    // Endpoints da API
    r.GET("/users/:id", getUserHandler(log))
    r.POST("/users", createUserHandler(log))

    log.InfoCtx("Servidor iniciando na porta :8080", nil)
    r.Run(":8080")
}
```

## Performance e Benchmarks

O Logz foi extensivamente testado para uso em produção:

- ✅ **Operações Concorrentes**: Lida com sucesso com 500+ operações de logging simultâneas
- ✅ **Zero Race Conditions**: Design thread-safe com uso otimizado de mutex
- ✅ **Eficiente em Memória**: Overhead mínimo de memória com gerenciamento inteligente de recursos
- ✅ **Resiliente a Rede**: Tratamento robusto de erros para falhas HTTP e WebSocket
- ✅ **Pronto para Prometheus**: Coleta de métricas com impacto mínimo na performance

### Uso no Mundo Real

Atualmente implantado e testado em batalha em sistemas de produção incluindo:

- **Sistema Backend GOBE**: Backend completo com suporte MCP (Model Context Protocol)
- **APIs de Alto Tráfego**: APIs REST com milhares de requisições por minuto
- **Arquiteturas de Microsserviços**: Sistemas distribuídos com monitoramento em tempo real
- **Pipelines CI/CD**: Workflows automatizados de deployment e monitoramento

## Interface de Linha de Comando

Aqui estão alguns exemplos de comandos que podem ser executados com a CLI:

```bash
# Registrar logs de diferentes níveis
logz info --msg "Iniciando a aplicação."
logz error --msg "Erro ao se conectar ao banco de dados."

# Iniciar o serviço destacado
logz start

# Parar o serviço destacado
logz stop

# Monitorar logs em tempo real
logz watch
```

### Configuração

O Logz utiliza um arquivo de configuração JSON ou YAML para centralizar sua configuração. O arquivo será automaticamente gerado no primeiro uso ou pode ser configurado manualmente em:
`~/.kubex/logz/config.json`.

Exemplo de configuração:

```json
{
  "port": "2112",
  "bindAddress": "0.0.0.0",
  "logLevel": "info",
  "notifiers": {
    "webhook1": {
      "type": "http",
      "webhookURL": "https://example.com/webhook",
      "authToken": "seu-token-aqui"
    }
  }
}
```

---

## **Integração com Prometheus**

Uma vez iniciado, o Logz expõe métricas no endpoint:

```plaintext
http://localhost:2112/metrics
```

### Exemplo de configuração no Prometheus

```yaml
scrape_configs:
  - job_name: 'logz'
    static_configs:
      - targets: ['localhost:2112']
```

---

## **Roadmap**

🔜 **Próximos Recursos**:

- Suporte a novos tipos de notificadores (como Slack, Discord e e-mails).
- Painel de monitoramento integrado.
- Configuração avançada com validações automáticas.

---

## **Contribuições**

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests. Confira o [guia de contribuições](CONTRIBUTING.md) para mais detalhes.

---

## **Contato**

Rafael Mori

- 🌐 [Portfolio](https://rafa-mori.dev)
- 🔗 [LinkedIn](https://www.linkedin.com/in/rafa-mori/)
- 📧 [Email](mailto:faelmori@gmail.com)
- 💼 Follow me on GitHub:
  - [faelmori](https://github.com/faelmori)
  - [rafa-mori](https://github.com/kubex-ecosystem)

Adoraria ouvir sobre novas oportunidades de trabalho ou colaborações. Se você gostou desse projeto, não hesite em entrar em contato comigo!
