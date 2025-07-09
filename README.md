# AuthGate - gRPC Authentication Service

A robust and scalable authentication service built with Go and gRPC, providing secure user authentication, registration, and token management capabilities.

## 🚀 Features

- **User Registration & Authentication** - Support for multiple identifier types (Email, CPF, CNPJ, Phone)
- **JWT Token Management** - Access and refresh token handling with configurable expiration
- **Token Verification** - Secure token validation for protected resources
- **Password Encryption** - BCrypt password hashing for security
- **Clean Architecture** - Well-structured codebase following clean architecture principles
- **Database Integration** - PostgreSQL integration with GORM
- **Docker Support** - Containerized deployment with Docker and Docker Compose
- **Dependency Injection** - Using Uber's FX framework for dependency management

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Docker Deployment](#docker-deployment)
- [Architecture](#architecture)
- [Contributing](#contributing)

## 🛠 Prerequisites

- Go 1.24.0 or higher
- PostgreSQL database
- Docker and Docker Compose (for containerized deployment)
- Protocol Buffers compiler (protoc)

## 📦 Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Gabriel-Schiestl/authgate.git
   cd authgate
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Generate protobuf files (if needed):**
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative \
          --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          proto/auth.proto
   ```

## ⚙️ Configuration

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=authgate
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION_HOURS=24
REFRESH_TOKEN_EXPIRATION_HOURS=168

# Server Configuration
GRPC_PORT=50051

# Optional Security Settings
MAX_WRONG_ATTEMPTS=5
MAX_TOKEN_AGE_SECONDS=86400
```

## 🚀 Usage

### Running the Service

1. **Local Development:**

   ```bash
   go run cmd/main.go
   ```

2. **Using Docker Compose:**
   ```bash
   docker-compose up -d
   ```

The gRPC server will start on port `50051` by default.

### Client Connection

Connect to the service using any gRPC client:

```go
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := authpb.NewAuthServiceClient(conn)
```

## 📖 API Documentation

### Service Methods

#### 1. Register

Register a new user with various identifier types.

```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse);
```

**Request:**

```protobuf
message RegisterRequest {
    IdentifierType identifier_type = 1;
    string identifier_value = 2;
    string password = 3;
    UserInfo user_info = 4;
    bool encrypt_token = 5;
    optional int32 max_wrong_attempts = 6;
    optional int32 max_token_age_seconds = 7;
}
```

#### 2. Login

Authenticate user and retrieve access/refresh tokens.

```protobuf
rpc Login(LoginRequest) returns (LoginResponse);
```

**Request:**

```protobuf
message LoginRequest {
    IdentifierType identifier_type = 1;
    string identifier_value = 2;
    string password = 3;
}
```

#### 3. VerifyToken

Validate an access token and retrieve user information.

```protobuf
rpc VerifyToken(VerifyTokenRequest) returns (VerifyTokenResponse);
```

#### 4. RefreshToken

Generate a new access token using a refresh token.

```protobuf
rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
```

#### 5. DeleteAuth

Remove user authentication data.

```protobuf
rpc DeleteAuth(DeleteAuthRequest) returns (DeleteAuthResponse);
```

### Supported Identifier Types

- `IDENTIFIER_TYPE_EMAIL` - Email address
- `IDENTIFIER_TYPE_CPF` - Brazilian CPF
- `IDENTIFIER_TYPE_CNPJ` - Brazilian CNPJ
- `IDENTIFIER_TYPE_PHONE` - Phone number

## 🛠 Development

### Project Structure

```
authgate/
├── cmd/                    # Application entrypoints
│   └── main.go
├── internal/src/
│   ├── application/        # Application layer
│   │   ├── dtos/          # Data Transfer Objects
│   │   └── usecases/      # Business logic
│   ├── config/            # Configuration
│   ├── controller/        # Controllers
│   ├── domain/            # Domain layer
│   │   ├── models/        # Domain models
│   │   ├── repositories/  # Repository interfaces
│   │   └── services/      # Service interfaces
│   ├── infra/             # Infrastructure layer
│   │   ├── adapters/      # Service implementations
│   │   ├── database/      # Database implementations
│   │   └── entities/      # Database entities
│   ├── module/            # Dependency injection
│   └── server/            # gRPC server
├── proto/                 # Protocol buffer definitions
├── authpb/               # Generated protobuf files
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile           # Docker image definition
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building the Application

```bash
go build -o authgate cmd/main.go
```

## 🐳 Docker Deployment

### Building Docker Image

```bash
docker build -t authgate:latest .
```

### Using Docker Compose

1. **Start the services:**

   ```bash
   docker-compose up -d
   ```

2. **View logs:**

   ```bash
   docker-compose logs -f app
   ```

3. **Stop the services:**
   ```bash
   docker-compose down
   ```

## 🏗 Architecture

This project follows **Clean Architecture** principles:

- **Domain Layer** - Core business logic and entities
- **Application Layer** - Use cases and DTOs
- **Infrastructure Layer** - External dependencies (database, encryption)
- **Presentation Layer** - gRPC server and controllers

### Key Components

- **gRPC Server** - Handles client requests
- **Use Cases** - Business logic implementation
- **Repositories** - Data access abstraction
- **Services** - External service interfaces (JWT, encryption)
- **Entities** - Database models
- **Mappers** - Data transformation between layers

## 🔧 Dependencies

### Core Dependencies

- **gRPC** - Remote procedure call framework
- **GORM** - ORM for database operations
- **JWT** - JSON Web Token implementation
- **BCrypt** - Password hashing
- **FX** - Dependency injection framework

### Development Dependencies

- **Protocol Buffers** - Interface definition language
- **PostgreSQL Driver** - Database connectivity
- **GoDotEnv** - Environment variable loading

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 Links

- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📧 Contact

Gabriel Schiestl - [GitHub](https://github.com/Gabriel-Schiestl)

Project Link: [https://github.com/Gabriel-Schiestl/authgate](https://github.com/Gabriel-Schiestl/authgate)
