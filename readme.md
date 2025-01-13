## Prerequisites Installation

To set up the development environment for this project, ensure the following tools are installed:

### 1. Install Protocol Buffers (protobuf)

Protocol Buffers (protobuf) is required for generating gRPC and microservice stubs.

#### Installation:

- **macOS**:
  ```bash
  brew install protobuf
  ```

- **Linux**:
  ```bash
  sudo apt-get install protobuf-compiler
  ```

### Installing Protocol Buffers generators (protoc-gen-go, protoc-gen-go-grpc, and protoc-gen-micro):

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/micro/go-micro/cmd/protoc-gen-micro@latest
```

### Installing migration tool dbmate:
```shell
npm i -g dbmate // NPM
brew install dbmate // MacOS

// for Linux users
sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
sudo chmod +x /usr/local/bin/dbmate
```

Usage:
```shell
make migration-new-[migration_name] // create new migration
make migration-up                   // run migration
make migration-down                 // rollback migration
```

Click for [more info](https://github.com/amacneil/dbmate)

### **Setup dependencies from Go Modules**
https://docs.gitlab.com/ee/user/project/use_project_as_go_package.html#authenticate-go-requests-to-private-projects

## Project Structure
```shell
.
├── cmd                                     # Main entry points for different commands
│   ├── apis                                # API-related commands
│   │   ├── api-gateway                     # API Gateway service
│   │   │   ├── config                      # Configuration files specific to the API Gateway
│   │   │   └── internal                    # Internal logic for the API Gateway, such as authentication and routing
│   │   └── user-api                        # User API service
│   │       ├── config                      # Configuration files specific to the User API
│   │       └── handlers                    # Handlers for user-related requests in the User API
│   ├── consumers                           # Contains consumer services, typically for processing messages from a queue
│   ├── schedulers                          # Services for scheduled tasks or cron jobs
│   └── services                            # Microservices within the application
│       └── user-service                    # User-related service
│           ├── config                      # Configuration files specific to the User Service
│           └── handlers                    # Handlers for user-related logic in the User Service
├── internal                                # Internal package for shared logic, accessible only within the module
├── pkg                                     # Core application components and utilities
│   ├── app                                 # Core application logic, including initialization functions for APIs and services
│   ├── common                              # Common utilities and helpers shared across the application
│   ├── config                              # Configuration management, such as loading from files or environment variables
│   ├── const                               # Constants used across the application
│   ├── dbtool                              # Database-related utilities, such as connection management or migrations
│   ├── logger                              # Logging setup and utilities for structured logging
│   ├── transport                           # Transport layer handling HTTP or gRPC requests and responses
│   └── utils                               # Utility functions and libraries shared across the application
│       └── codec                           # Encoding and decoding utilities
│           └── bejson                      # Custom JSON encoding/decoding utilities
├── proto                                   # Protocol Buffer definitions and generated files
│   ├── exmsg                               # External message-related definitions
│   │   ├── models                          # Generated Protocol Buffers models for external messages
│   │   └── services                        # Generated gRPC services for external message handling
│   ├── models                              # Protocol Buffers definitions for shared data models
│   └── services                            # Protocol Buffers definitions for gRPC services
├── scripts                                 # Helper scripts for tasks such as deployment, setup, or migrations
└── tests                                   # Test configurations and scripts
    └── config                              # Configuration files and setup for testing environments

```
## Getting Started

### 1. Set up infrastructure

Before starting the development, you need to set up the necessary infrastructure, such as databases, message queues, and
caching systems
Or you can use Docker for testing and deploying the application.

#### 1. Native Infrastructure

- Install Consul for service discovery and configuration management.
- Install RabbitMQ for message queues.
- Install Redis for caching and storing session data.
- Install Posgres for database.

#### 2. Docker for Testing and Deploying

- Run docker-compose to set up the infrastructure by using command:
```
make start-services
```

### 3. Build and run the microservices (for cicd)
```shell
DOCKER_BUILDKIT=1 docker build --build-arg BIN=cmd/consumers/sample-consumer -f scripts/Dockerfile --ssh default=~/.ssh/id_rsa . -t sample-consumer

docker run -e CONSUL_ADDRESS=host.docker.internal:8500 -p 30000:30000 user-api
```


## Generate proto file
- cd to root directory of the project and run the following command:
- install atcli tools
 ```shell
  go install pkg/toolkit/atcli.go
 ```
### for models/entities
```shell
 atcli gen proto proto/models/jwt.proto #path to proto file from root dir 
```
or
```shell
protoc \
  --proto_path=proto/models \
  --go_out=paths=source_relative:proto/exmsg/models \
  --go-grpc_out=paths=source_relative:proto/exmsg/models \
  user.proto

```

### for grpc services
```shell
 atcli gen proto proto/services/users.proto #path to proto file from root dir 
```
or 
```shell
protoc \
  --plugin=protoc-gen-micro=$GOPATH/bin/protoc-gen-micro \
  --proto_path=proto \
  --go_out=paths=source_relative:proto/exmsg \
  --go-grpc_out=paths=source_relative:proto/exmsg \
  --micro_out=paths=source_relative:proto/exmsg \
  services/user.proto
```

---

# Tech Stack Overview

This project leverages a robust tech stack to deliver high-performance, scalable, and resilient services. Here’s a
breakdown of the core technologies used:

## Golang

- **Description**: The core programming language used to build our microservices, known for its simplicity, efficiency,
  and strong support for concurrency.

## PostgreSQL

- **Description**: A powerful, open-source relational database system, used as the primary data storage solution for
  structured data.

## Consul

- **Purpose**: Configuration management and service discovery.
- **Description**: A service mesh solution providing configuration management, service discovery, and health checking to
  ensure the system's resilience and scalability.

## Go Micro

- **Purpose**: Web Framework.
- **Description**: A framework for building and managing microservices in Go, offering out-of-the-box solutions for RPC
  communication, load balancing, and more.

## RabbitMQ

- **Purpose**: Message Queue.
- **Description**: A messaging broker that facilitates asynchronous communication between services, enhancing
  reliability and enabling message-based workflows.

## gRPC with Protocol Buffers

- **Purpose**: Service Transport.
- **Description**: A high-performance RPC framework that uses Protocol Buffers as the interface definition language,
  enabling efficient, language-agnostic communication between services.
