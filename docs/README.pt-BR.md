![Logz Banner](assets/top_banner.png)

---

**Uma ferramenta de gerenciamento de logs e métricas avançada, com suporte nativo à integração com Prometheus, notificações dinâmicas e uma CLI poderosa.**

---

## **Índice**
1. [Sobre o Projeto](#sobre-o-projeto)
2. [Destaques](#destaques)
3. [Instalação](#instalação)
4. [Uso](#uso)
    - [CLI](#cli)
    - [Configuração](#configuração)
5. [Integração com Prometheus](#integração-com-prometheus)
6. [Roadmap](#roadmap)
7. [Contribuições](#contribuições)
8. [Contato](#contato)

---

## **Sobre o Projeto**
O Logz é uma solução poderosa e flexível para gerenciar logs e métricas em sistemas modernos. Construído com **Go**, oferece suporte extensivo a múltiplos métodos de notificação, incluindo **HTTP Webhooks**, **ZeroMQ**, e **DBus**, além de integração fluida com o **Prometheus** para monitoramento avançado.

O objetivo é fornecer uma ferramenta robusta, altamente configurável e escalável para desenvolvedores, equipes DevOps e arquitetos de software que precisam de uma abordagem centralizada para gerenciar logs e métricas.

**Principais Benefícios:**
- 💡 Fácil de configurar e usar.
- 🌐 Integração direta com Prometheus e outros sistemas.
- 🔧 Extensível com novos notifiers e serviços.

---

## **Destaques**
✨ **Notificadores Dinâmicos**:
- Suporte a múltiplos notifiers simultaneamente.
- Configuração centralizada e flexível via JSON ou YAML.

📊 **Monitoramento e Métricas**:
- Exposição de métricas compatíveis com Prometheus.
- Gerenciamento dinâmico de métricas com suporte a persistência.

💻 **CLI Poderoso**:
- Comandos simples e objetivos para gerenciar logs e serviços.
- Extensível para novos fluxos de trabalho.

🔒 **Resiliente e Seguro**:
- Validações em conformidade com as regras do Prometheus.
- Modos de operação distintos para processos destacados e convencionais.

---

## **Instalação**
Requisitos:
- **Go** versão 1.19 ou superior.
- Prometheus (opcional para monitoramento avançado).

```bash
# Clone este repositório
git clone https://github.com/faelmori/logz.git

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

### CLI
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
```
http://localhost:2112/metrics
```

### Exemplo de configuração no Prometheus:
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
💌 **Desenvolvedor:**  
[Seu Nome](mailto:seu-email@dominio.com)  
🌐 [Seu LinkedIn](https://linkedin.com/in/seu-perfil)  
💼 [Portfólio](https://seu-portfolio.com)

Adoraria ouvir sobre novas oportunidades de trabalho ou colaborações. Se você gostou desse projeto, não hesite em entrar em contato comigo!
