graph TD
    CLI[CLI (cobra)] -->|calls| Factory[Factory Layer]
    Factory -->|SetLogger, Config| InternalConfig[internal Config]
    Factory -->|SetNotifier| InternalNotifier[internal Notifier]
    InternalConfig -->|Log, Debug, Info| CoreLogger[Core Logger Logic]
    InternalNotifier -->|Alert, Notify| NotifierMgr[Notifier Manager]
    CoreLogger -->|Write JSON/Text| Writer[log writer & format]
    NotifierMgr -->|Send HTTP/Prometheus| PromHTTP[Prometheus HTTP endpoint]
    CLI -->|can import as lib| CoreLogger
    Writer -->|Outputs| Output[Console, File, etc.]
    PromHTTP -->|/metrics| Output