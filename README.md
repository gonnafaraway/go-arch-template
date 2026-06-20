# go-arch-template

Best optimized and effective template to develop services and apps with Go.

## Architectural Schema

<img width="674" height="693" alt="{7C048F03-FDD9-43A0-A643-07AE5D8BAD0A}" src="https://github.com/user-attachments/assets/9b6abf9c-335d-472a-a2a9-9ca2f39bd953" />

## System Architecture Overview

This project follows **Clean Architecture** (also known as Hexagonal Architecture or Ports & Adapters) principles, combined with **Domain-Driven Design (DDD)** concepts. The architecture is organized in layers with clear separation of concerns and dependency inversion.

### Architecture Layers

#### 1. **Domain Layer** (`internal/api/domain/`)
The innermost layer containing business logic and domain entities. This layer has no dependencies on external frameworks or infrastructure.

- **Entities**: Core business objects (Company, User, Order)
- **Domain Validators**: Business rule validation logic
- **Domain Errors**: Domain-specific error types

**Principles Applied:**
- **Domain-Driven Design (DDD)**: Rich domain models with business logic
- **Single Responsibility Principle**: Each entity has a clear, focused purpose
- **Dependency Inversion**: Domain layer has no external dependencies

#### 2. **Use Case Layer** (`internal/api/usecase/`)
Contains application-specific business logic and orchestrates domain entities. This layer implements the application's use cases.

- **Use Cases**: Application workflows (CreateCompany, CreateOrder, etc.)
- **Request/Response DTOs**: Data transfer objects for use case boundaries
- **Orchestration**: Coordinates between repositories, domain entities, and integrations

**Principles Applied:**
- **Use Case Pattern**: Each use case represents a single business operation
- **Command Pattern**: Use cases encapsulate operations as commands
- **Dependency Injection**: Use cases depend on abstractions (interfaces)

#### 3. **Repository Layer** (`internal/api/repository/`)
Abstracts data access logic. Provides interfaces for data persistence and implements the Repository pattern.

- **Repository Interfaces**: Define data access contracts
- **Multiple Implementations**: PostgreSQL, MongoDB, and Mock repositories
- **Strategy Pattern**: Different storage strategies can be swapped

**Principles Applied:**
- **Repository Pattern**: Abstracts data access from business logic
- **Strategy Pattern**: Multiple repository implementations
- **Interface Segregation**: Focused, specific repository interfaces
- **Dependency Inversion**: Business logic depends on repository interfaces, not implementations

#### 4. **Infrastructure Layer** (`internal/api/infrastructure/`)
Contains technical implementations and external service integrations.

**Sub-layers:**

- **Storage** (`storage/`): Database clients (PostgreSQL, MongoDB, Redis, Kafka, S3, Nexus)
- **External Services** (`external/`): Clients for external APIs (Billing, Users Service, Prometheus, Sentry)
- **Local Services** (`local/`): Local infrastructure (Logging, Tracing, OAuth)

**Principles Applied:**
- **Adapter Pattern**: Adapts external services to internal interfaces
- **Factory Pattern**: Creates infrastructure components
- **Dependency Injection**: Infrastructure is injected into upper layers

#### 5. **Integration Layer** (`internal/api/integration/`)
Orchestrates interactions with external services and systems.

- **Integration Interfaces**: Define contracts for external service interactions
- **Integration Implementations**: Concrete implementations for billing, company sync, OAuth
- **Anti-Corruption Layer**: Protects domain from external service changes

**Principles Applied:**
- **Adapter Pattern**: Adapts external services to application needs
- **Anti-Corruption Layer**: Isolates domain from external dependencies
- **Interface Segregation**: Focused integration interfaces

#### 6. **Transport Layer** (`internal/api/transport/`)
Handles communication protocols (HTTP, gRPC, RPC).

- **HTTP Transport**: REST API implementation
- **HTTP Middleware**: Cross-cutting concerns (Logging, Tracing, Metrics, Sentry)
- **gRPC Transport**: gRPC server implementation
- **Response Handling**: Standardized response formatting

**Principles Applied:**
- **Middleware Pattern**: Chain of responsibility for cross-cutting concerns
- **Decorator Pattern**: Middleware decorates handlers
- **Separation of Concerns**: Transport logic separated from business logic

#### 7. **Handler Layer** (`internal/api/handlers/`)
HTTP request handlers that translate HTTP requests to use case calls.

- **Request Handlers**: Handle HTTP requests and responses
- **Route Binding**: Maps HTTP routes to handlers
- **Request Validation**: Input validation before use case execution

**Principles Applied:**
- **Controller Pattern**: Handlers act as controllers
- **Thin Controllers**: Minimal logic, delegates to use cases
- **Single Responsibility**: Each handler focuses on one resource

#### 8. **Application Layer** (`internal/api/app/`)
Application bootstrap and composition root. Assembles all components using Dependency Injection.

- **Component Assembly**: Wires up all layers
- **Service Orchestration**: Manages service lifecycle
- **Dependency Resolution**: Resolves dependencies between components

**Principles Applied:**
- **Composition Root**: Single place for dependency resolution
- **Dependency Injection**: All dependencies injected
- **Service Locator Anti-Pattern Avoided**: Explicit dependency injection

#### 9. **Service Layer** (`internal/api/service/`)
Long-running services and background jobs.

- **API Service**: HTTP API server
- **Jobs Service**: Background job processing
- **CDC Service**: Change Data Capture service
- **Service Controller**: Manages multiple services lifecycle

**Principles Applied:**
- **Service Pattern**: Encapsulates service lifecycle
- **Concurrent Execution**: Services run concurrently
- **Graceful Shutdown**: Proper resource cleanup

### Design Patterns Used

1. **Repository Pattern**: Abstracts data access
2. **Factory Pattern**: Component creation (`Prepare*` functions)
3. **Strategy Pattern**: Multiple repository/storage implementations
4. **Adapter Pattern**: External service integration
5. **Dependency Injection**: Loose coupling through interfaces
6. **Middleware Pattern**: Cross-cutting concerns in HTTP layer
7. **Use Case Pattern**: Application business logic encapsulation
8. **Command Pattern**: Use cases as commands
9. **DTO Pattern**: Data transfer objects for layer boundaries
10. **Anti-Corruption Layer**: Protection from external dependencies

### Architectural Principles

1. **Clean Architecture**: Dependency rule - dependencies point inward
2. **SOLID Principles**:
   - **S**ingle Responsibility: Each component has one reason to change
   - **O**pen/Closed: Open for extension, closed for modification
   - **L**iskov Substitution: Implementations are substitutable
   - **I**nterface Segregation: Focused, specific interfaces
   - **D**ependency Inversion: Depend on abstractions, not concretions

3. **Domain-Driven Design (DDD)**:
   - Rich domain models
   - Domain validation
   - Ubiquitous language

4. **Separation of Concerns**: Clear boundaries between layers
5. **Dependency Inversion**: High-level modules don't depend on low-level modules
6. **Interface-Based Design**: Program to interfaces, not implementations

### Data Flow

```
HTTP Request
    ↓
Transport Layer (HTTP Middleware)
    ↓
Handler Layer (Request/Response Translation)
    ↓
Use Case Layer (Business Logic Orchestration)
    ↓
Domain Layer (Business Rules & Validation)
    ↓
Repository Layer (Data Access)
    ↓
Infrastructure Layer (Database/Storage)
```

### Observability

The architecture includes comprehensive observability:

- **Logging**: Structured logging with Zap
- **Tracing**: OpenTelemetry with Jaeger/OTLP support
- **Metrics**: Prometheus metrics middleware
- **Error Tracking**: Sentry integration
- **Fallback Mechanisms**: Graceful degradation when observability tools unavailable

### Storage Support

The architecture supports multiple storage backends:

- **PostgreSQL**: Relational database
- **MongoDB**: Document database
- **Redis**: Caching layer
- **Kafka**: Message queue
- **S3**: Object storage
- **Nexus**: Artifact repository

Repository pattern allows switching between storage implementations without changing business logic.

### Testing Strategy

- **Mock Repositories**: For testing without database
- **Interface-Based Testing**: Easy to mock dependencies
- **Domain Testing**: Test business logic in isolation
- **Integration Testing**: Test component interactions

## Data Architecture Overview

Describes common patterns for working with data models when using clean architecture in Go. 
The goal is to unify the design approach and avoid mixing responsibilities.

## Patterns

### 1. Entities (Domain Models)

**Purpose:** represent key business concepts of the application (e.g., `User`, `Order`). Contain core business logic and invariants.

**Characteristics:**
* independent of external details (DB, frameworks, HTTP);
* located in the `domain` layer of clean architecture;
* may contain methods for validation and business rules.

**Example (Go):**
```go
package domain

import "errors"

type User struct {
    ID    string
    Name  string
    Email string
}

func NewUser(name, email string) (*User, error) {
    if name == "" {
        return nil, errors.New("name is required")
    }
    // ... other business validation
    return &User{Name: name, Email: email}, nil
}
```

### 2. Value Objects
**Purpose:** describe characteristics or composite values without their own identifier (e.g., Email, Money, Address).

**Characteristics:**
* immutable
* compared by value, not by reference
* often used as fields in entities

```go
package domain

import (
    "errors"
    "strings"
)

type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    if !strings.Contains(value, "@") {
        return nil, errors.New("invalid email format")
    }
    return &Email{value: value}, nil
}

func (e Email) Value() string {
    return e.value
}
```

### 3. DTO (Data Transfer Object)
**Purpose:** transfer data between application layers or via API (e.g., for HTTP requests/responses).

**Characteristics:**
* simple structure, often without methods
* may differ from the domain model (e.g., contain fields for serialization)
* located in the interface adapters or transport layer

#### DTO example for HTTP response (Go):

```go
package transport

type UserResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### 4. Data Mapper
**Purpose:** convert data between formats (e.g., from domain entity to DB row and back).

**Characteristics:**
* separates domain logic from persistence details
* implements mapping functions between domain models and data storage formats
* often used with Repository pattern

#### Example (Go):
```go
package infrastructure

import "your-app/domain"

type UserMapper struct{}

func (m UserMapper) ToDB(user *domain.User) map[string]interface{} {
    return map[string]interface{}{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
    }
}

func (m UserMapper) FromDB(row map[string]interface{}) *domain.User {
    return &domain.User{
        ID:    row["id"].(string),
        Name:  row["name"].(string),
        Email: row["email"].(string),
    }
}
```

### 5. Repository
**Purpose:** act as an intermediary between the domain model and the data source (DB, external API).

**Characteristics:**
* hides storage details from business logic
* provides a clean interface for data operations
* allows easy switching between different storage implementations
#### Example (Go):
```go
package repository

import "your-app/domain"

type UserRepository interface {
    Create(user *domain.User) error
    FindByID(id string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id string) error
}
```

### 6. View Model
**Purpose:** adapt data specifically for a particular interface (e.g., frontend framework requirements).
**Characteristics:**
* similar to DTO but focused on UI needs
* may include additional computed fields for display
* located in presentation layer

### Full Data Flow in Clean Architecture
**Typical data flow when processing an HTTP request:**
* Incoming HTTP request → DTO (transport layer)
* DTO → mapped to Domain Entity (interface adapters layer)
* Business logic operates on the Entity (use cases layer)
* Entity passed to Repository → Data Mapper converts to DB format
* Response: domain data mapped to DTO/View Model → serialized to JSON/XML

## Best Practices
* Separate responsibilities: DTO for transfer, Entity for logic, Data Mapper for conversion.
* Prefer immutability: make models immutable where possible for thread safety.
* Use mapping libraries: consider using libraries like mapstructure to reduce boilerplate.
* Plan versioning: design DTO/View Model versioning early for API evolution.
* Validate early: perform input validation at the transport layer.
* Keep domain pure: avoid importing infrastructure packages in domain models.

## Directory Structure Example
```go
your-app/
├── domain/                  # Entities and Value Objects
│   ├── user.go
│   └── email.go
├── use-cases/             # Business logic
├── interface-adapters/
│   ├── transport/         # DTOs and HTTP handlers
│   │   └── user_dto.go
│   └── repository/        # Repository interfaces
├── infrastructure/        # Data Mappers and DB implementations
│   ├── mappers/
│   └── db/
└── main.go
```
---

## Contributing

This project is open to development and welcomes proactive pull requests! We encourage contributions that:

- Improve code quality and maintainability
- Add new features following the established architecture
- Enhance documentation
- Fix bugs and improve error handling
- Optimize performance
- Add new integrations or storage backends
- Improve observability and monitoring

Please ensure that any contributions follow the architectural principles and patterns established in this template. We're looking forward to your contributions!
