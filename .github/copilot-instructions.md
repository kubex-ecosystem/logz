# AI Coding Instructions for GoBE Backend

## Project Overview

**GoBE** is a modular and secure HTTP backend server built with Go and Gin framework. The project emphasizes zero-configuration security, interface segregation, and clean architecture patterns with integrated Docker-based database management via the gdbase factory system.

### Core Technologies
- **Framework**: Gin (HTTP routing and middleware)
- **Database**: GORM with PostgreSQL/SQLite via gdbase factory
- **Architecture**: Clean Architecture with Repository/Service patterns
- **Security**: Certificate-based authentication, keyring storage, CORS
- **CLI**: Cobra-based command interface
- **Build**: Make-based with cross-platform GitHub Actions

## Architectural Patterns

### 1. Router System (`internal/routes/`)

The router system follows an interface-based architecture with centralized route management:

```go
// Core pattern for route creation
func NewEntityRoutes(rtr *ar.IRouter) map[string]ar.IRoute {
    if rtr == nil {
        gl.Log("error", "Router is nil for EntityRoute")
        return nil
    }
    rtl := *rtr
    
    dbService := rtl.GetDatabaseService()
    dbGorm, err := dbService.GetDB()
    if err != nil {
        gl.Log("error", "Failed to get DB from service", err)
        return nil
    }
    
    entityController := entity.NewController(dbGorm)
    routesMap := make(map[string]ar.IRoute)
    middlewaresMap := make(map[string]any)
    
    routesMap["RouteKey"] = NewRoute(
        http.MethodGet, 
        "/path", 
        "application/json", 
        entityController.Handler, 
        middlewaresMap, 
        dbService,
    )
    
    return routesMap
}
```

**Key Points:**
- Always validate router is not nil
- Extract database service from router interface
- Create controller with GORM database connection
- Use descriptive route keys for debugging
- Include middleware maps even if empty

### 2. MCP (Model Context Protocol) Architecture

The MCP system follows a strict tripé pattern: Model → Repository → Service → Controller

#### Factory Pattern (`gdbase/factory/models/mcp/`)
```go
// Factory functions for dependency injection
func NewEntityService(entityRepo EntityRepo) EntityService {
    return m.NewEntityService(entityRepo)
}

func NewEntityRepo(db *gorm.DB) EntityRepo {
    return m.NewEntityRepo(db)
}

func NewEntityModel(field1 string, field2 JsonB) EntityModel {
    return &m.EntityModel{
        ID:        uuid.New().String(),
        Field1:    field1,
        Field2:    field2,
        CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
        UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
        CreatedBy: "admin", // Should be set by service layer
    }
}
```

#### Repository Layer (`gdbase/internal/models/mcp/*/repo.go`)
```go
type IEntityRepo interface {
    TableName() string
    Create(e IEntityModel) (IEntityModel, error)
    FindOne(where ...interface{}) (IEntityModel, error)
    FindAll(where ...interface{}) ([]IEntityModel, error)
    Update(e IEntityModel) (IEntityModel, error)
    Delete(id string) error
    Close() error
    List(where ...interface{}) (xtt.TableDataHandler, error)
    GetContextDBService() t.IDBService
}

// Always validate nil models
func (er *EntityRepo) Create(e IEntityModel) (IEntityModel, error) {
    if e == nil {
        return nil, fmt.Errorf("EntityModel repository: EntityModel is nil")
    }

    eModel, ok := e.(*EntityModel)
    if !ok {
        return nil, fmt.Errorf("EntityModel repository: model is not of type *EntityModel")
    }

    // UUID validation and generation
    if eModel.GetID() != "" {
        if _, err := uuid.Parse(eModel.GetID()); err != nil {
            return nil, fmt.Errorf("EntityModel repository: invalid UUID: %w", err)
        }
    } else {
        eModel.SetID(uuid.New().String())
    }

    // Always validate before database operations
    if err := eModel.Validate(); err != nil {
        return nil, fmt.Errorf("EntityModel repository: validation failed: %w", err)
    }

    // Always sanitize before saving
    eModel.Sanitize()

    err := er.g.Create(eModel).Error
    if err != nil {
        return nil, fmt.Errorf("EntityModel repository: failed to create: %w", err)
    }
    return eModel, nil
}
```

#### Service Layer (`gdbase/internal/models/mcp/*/service.go`)
```go
type IEntityService interface {
    CreateEntity(entity IEntityModel) (IEntityModel, error)
    GetEntityByID(id string) (IEntityModel, error)
    UpdateEntity(entity IEntityModel) (IEntityModel, error)
    DeleteEntity(id string) error
    ListEntities() ([]IEntityModel, error)
}

func (es *EntityService) CreateEntity(entity IEntityModel) (IEntityModel, error) {
    // Business validation at service level
    if entity.GetRequiredField() == "" {
        return nil, errors.New("missing required field: required_field is required")
    }

    // Delegate model validation to repository
    if err := entity.Validate(); err != nil {
        return nil, fmt.Errorf("validation error: %w", err)
    }

    createdEntity, err := es.repo.Create(entity)
    if err != nil {
        return nil, fmt.Errorf("error creating entity: %w", err)
    }
    return createdEntity, nil
}
```

#### Controller Layer (`gobe/internal/controllers/mcp/*/controller.go`)
```go
type Controller struct {
    service svc.Service
}

func NewController(db *gorm.DB) *Controller {
    return &Controller{
        service: svc.NewService(models.NewRepo(db)),
    }
}

func (c *Controller) CreateEntity(ctx *gin.Context) {
    var entity EntityCreateRequest
    if err := ctx.ShouldBindJSON(&entity); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    createdEntity, err := c.service.CreateEntity(&entity)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, createdEntity)
}
```

### 3. Security and Middleware System

The security system is initialized in `internal/routes/utils.go`:

```go
func SecureServerInit(r *gin.Engine, fullBindAddress string) error {
    // Trusted proxy validation
    trustedProxies, err := getTrustedProxies()
    if err != nil {
        return err
    }
    r.SetTrustedProxies(trustedProxies)

    // Global security middleware
    r.Use(func(c *gin.Context) {
        // Host validation
        if !validateExpectedHosts(fullBindAddress, c) {
            c.Abort()
            return
        }

        // CORS headers (comprehensive)
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        // Security headers
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("X-Content-Type-Options", "nosniff")

        // Handle OPTIONS preflight
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }

        c.Next()
    })

    return nil
}
```

**Middleware Registration Pattern:**
```go
func (rtr *Router) RegisterMiddleware(name string, middleware gin.HandlerFunc, global bool) {
    if global {
        rtr.engine.Use(middleware)
    } else {
        rtr.middlewares[name] = middleware
    }
}

// Route-specific middleware application
if route.Secure() {
    if authMdw, ok := rtr.middlewares["authentication"]; ok {
        middlewaresStack = append(middlewaresStack, authMdw)
    }
}
```

### 4. Database Integration Patterns

#### Database Service Access
```go
// Standard pattern for getting database from router
dbService := rtl.GetDatabaseService()
dbGorm, err := dbService.GetDB()
if err != nil {
    gl.Log("error", "Failed to get DB from service", err)
    return nil
}
```

#### JSONB Support Pattern
```go
// For models with JSONB fields
type EntityModel struct {
    ID        string    `json:"id" gorm:"primary_key"`
    Config    t.JsonB   `json:"config" gorm:"type:jsonb"`
    // ... other fields
}

// Validation includes JSONB checks
func (e *EntityModel) Validate() error {
    if e.Config == nil {
        return errors.New("config cannot be nil")
    }
    return nil
}
```

### 5. Error Handling Patterns

#### Repository Level
```go
func (r *Repo) Operation(param Type) (Type, error) {
    if param == nil {
        return nil, fmt.Errorf("Repository: parameter is nil")
    }

    // Type assertion with proper error
    concrete, ok := param.(*ConcreteType)
    if !ok {
        return nil, fmt.Errorf("Repository: parameter is not of expected type")
    }

    // Validation before database operation
    if err := concrete.Validate(); err != nil {
        return nil, fmt.Errorf("Repository: validation failed: %w", err)
    }

    // Database operation with context
    err := r.g.Operation(concrete).Error
    if err != nil {
        return nil, fmt.Errorf("Repository: database operation failed: %w", err)
    }

    return concrete, nil
}
```

#### Controller Level
```go
func (c *Controller) Handler(ctx *gin.Context) {
    // Standard error response format
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Success response
    ctx.JSON(http.StatusOK, result)
}
```

## Development Guidelines

### 1. When Adding New MCP Controllers

1. **Create Factory Functions** in `gdbase/factory/models/mcp/entity.go`
2. **Implement Model** in `gdbase/internal/models/mcp/entity/entity_model.go`
3. **Implement Repository** in `gdbase/internal/models/mcp/entity/entity_repo.go`
4. **Implement Service** in `gdbase/internal/models/mcp/entity/entity_service.go`
5. **Create Controller** in `gobe/internal/controllers/mcp/entity/entity_controller.go`
6. **Add Routes** in `gobe/internal/routes/mcp_entity.go`
7. **Register Routes** in main router system

### 2. Route Registration Pattern

```go
// File: internal/routes/mcp_entity.go
func NewMCPEntityRoutes(rtr *ar.IRouter) map[string]ar.IRoute {
    // Standard validation and setup
    if rtr == nil {
        l.ErrorCtx("Router is nil for MCPEntityRoute", nil)
        return nil
    }
    
    rtl := *rtr
    dbService := rtl.GetDatabaseService()
    dbGorm, err := dbService.GetDB()
    if err != nil {
        gl.Log("error", "Failed to get DB from service", err)
        return nil
    }
    
    entityController := entity_controller.NewEntityController(dbGorm)
    routesMap := make(map[string]ar.IRoute)
    middlewaresMap := make(map[string]any)

    // Standard CRUD routes
    routesMap["ListEntities"] = NewRoute(http.MethodGet, "/mcp/entities", "application/json", entityController.GetAllEntities, middlewaresMap, dbService)
    routesMap["GetEntity"] = NewRoute(http.MethodGet, "/mcp/entities/:id", "application/json", entityController.GetEntityByID, middlewaresMap, dbService)
    routesMap["CreateEntity"] = NewRoute(http.MethodPost, "/mcp/entities", "application/json", entityController.CreateEntity, middlewaresMap, dbService)
    routesMap["UpdateEntity"] = NewRoute(http.MethodPut, "/mcp/entities/:id", "application/json", entityController.UpdateEntity, middlewaresMap, dbService)
    routesMap["DeleteEntity"] = NewRoute(http.MethodDelete, "/mcp/entities/:id", "application/json", entityController.DeleteEntity, middlewaresMap, dbService)

    return routesMap
}
```

### 3. Interface Pattern Usage

Always use interfaces for dependency injection and testing:

```go
// Repository interface
type IEntityRepo interface {
    // Core CRUD
    Create(e IEntityModel) (IEntityModel, error)
    FindOne(where ...interface{}) (IEntityModel, error)
    FindAll(where ...interface{}) ([]IEntityModel, error)
    Update(e IEntityModel) (IEntityModel, error)
    Delete(id string) error
    
    // Utility methods
    Close() error
    TableName() string
    List(where ...interface{}) (xtt.TableDataHandler, error)
    GetContextDBService() t.IDBService
}

// Service interface
type IEntityService interface {
    CreateEntity(entity IEntityModel) (IEntityModel, error)
    GetEntityByID(id string) (IEntityModel, error)
    UpdateEntity(entity IEntityModel) (IEntityModel, error)
    DeleteEntity(id string) error
    ListEntities() ([]IEntityModel, error)
}
```

### 4. Logging Patterns

Use the global logger with appropriate context:

```go
import gl "github.com/rafa-mori/gobe/logger"

// Error logging
gl.Log("error", "Operation failed", err)
gl.Log("error", fmt.Sprintf("Failed to process %s", entityType), err)

// Info logging
gl.Log("info", fmt.Sprintf("Route registered: [%s] %s", method, path))

// Warning logging
gl.Log("warn", "Middleware not found for route")
```

### 5. Testing Patterns

Follow the established testing patterns for new components:

```bash
# Compile check for MCP controllers
go build -v ./internal/controllers/mcp/
go build -v ./internal/routes/

# Run tests
go test ./internal/controllers/mcp/entity/...
go test ./internal/routes/...
```

## Common Pitfalls to Avoid

1. **Router Nil Checks**: Always validate router pointer before dereferencing
2. **Database Service Errors**: Handle database service extraction errors properly
3. **Type Assertions**: Always check type assertion success in repositories
4. **UUID Validation**: Validate UUIDs before database operations
5. **Model Validation**: Call `Validate()` and `Sanitize()` before database saves
6. **Interface Consistency**: Ensure all repo/service interfaces follow the established pattern
7. **Error Wrapping**: Always wrap errors with context using `fmt.Errorf`
8. **CORS Headers**: Use the established CORS pattern in `SecureServerInit`
9. **Middleware Stack**: Preserve middleware order and use `uniqueMiddlewareStack`
10. **Route Registration**: Use descriptive route keys for debugging

## Project Structure Reference

```
gobe/
├── cmd/                          # CLI commands and main entry
├── internal/
│   ├── controllers/              # HTTP handlers
│   │   ├── mcp/                 # MCP protocol controllers
│   │   └── users/               # User management
│   ├── routes/                   # Route definitions and registration
│   ├── interfaces/               # Interface definitions
│   └── types/                    # Type definitions
├── factory/                      # Dependency injection
├── logger/                       # Centralized logging
├── services/                     # Business logic services
└── utils/                        # Utility functions

gdbase/                           # Database layer (separate module)
├── factory/models/               # Model factory functions
├── internal/models/              # Model implementations
│   ├── mcp/                     # MCP models
│   ├── orders/                  # Order models
│   └── products/                # Product models
└── types/                        # Database types and interfaces
```

## Environment and Configuration

- Use Viper for configuration management
- Environment variables follow `GOBE_*` prefix pattern
- Database connections managed through gdbase factory
- Security settings configured via `SecureServerInit`
- CORS settings are comprehensive and permissive for development

## Documentation Standards

- Document all public interfaces and methods
- Include usage examples in complex functions
- Maintain the MCP_CONTROLLERS_SUMMARY.md for new controllers
- Update this file when adding new architectural patterns

Remember: This project prioritizes security, modularity, and clean interfaces. Always follow the established patterns for consistency and maintainability.


## Golang Craftsmanship Standards

Use Go Modules for dependency management. Keep `go.mod` and `go.sum` clean and minimal. Avoid indirect dependencies when possible.

Organize projects using idiomatic structure: `cmd/`, `cmd/cli/`, `internal/`, `internal/types`, `internal/interfaces`, `api/`, `support/`, `support/instructions`, `tests/`.  

Place the main CLI entrypoint in `cmd/main.go` and the library entrypoint in the root with the package at the same name of project.

Every package must contain a comment describing its purpose before the package declaration in one line. Use `// Package <name> ...` format.

Write **table-driven tests** with the standard `testing` package. For complex assertions, use `testify`. Coverage should be high on business logic, especially for error paths.

Mock dependencies via interfaces — never via `globals` or side effects. Benchmark performance-sensitive functions. Keep tests fast and deterministic.

Naming: `CamelCase` for exported, `camelCase` for internal. Avoid stutter in package names (e.g., `user.User` is wrong).

Functions must be small and cohesive. Return early. Nesting is a code smell. Handle errors explicitly. Don’t ignore them — even temporarily.

Favor **composition over inheritance**. Accept interfaces, return concrete structs. Document behavior at interface boundaries.

Always use `context.Context` for cancellation, timeouts, and tracing. Pass it explicitly — do not store it in structs.

Exported types, functions, and packages MUST include **godoc-compatible comments**. Start with the function/type name. Include usage examples when applicable.

README must be clear, technical and up to date. Include build instructions, feature summary, and example usage. If possible, add architecture diagrams and CLI reference.

Be consistent. Be happy. Be fast. Be safe.