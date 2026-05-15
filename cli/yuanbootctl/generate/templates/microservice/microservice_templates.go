package microservice

const Main_Tel = `package main

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/pkg/servicediscovery/nacos"
	"github.com/liangboceo/yuanboot/web"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/endpoints"
	"github.com/liangboceo/yuanboot/web/mvc"
	"github.com/liangboceo/yuanboot/web/router"
	"{{.ModelName}}/internal/controller"
	"{{.ModelName}}/internal/middleware"
	"{{.ModelName}}/internal/service"
)

func main() {
	host := CreateHostBuilder().Build()
	host.Run()
}

func CreateHostBuilder() *abstractions.HostBuilder {
	config := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddYamlFile("config/config").Build()

	return web.NewWebHostBuilder().
		UseConfiguration(config).
		Configure(func(app *web.ApplicationBuilder) {
			app.UseMiddleware(middleware.NewRecovery())
			app.UseMiddleware(middleware.NewLogger())
			app.UseMiddleware(middleware.NewCORS())
			app.UseEndpoints(registerEndpoints)
			app.UseMvc(func(builder *mvc.ControllerBuilder) {
				builder.AddViewsByConfig()
				builder.AddController(controller.NewUserController)
			})
		}).
		ConfigureServices(func(sc *dependencyinjection.ServiceCollection) {
			// Register services
			sc.AddTransientByImplements(service.NewUserService, new(service.IUserService))
			sc.AddTransientByImplements(service.NewUserRepository, new(service.IUserRepository))

			// Enable service discovery (Nacos/Eureka/Consul)
			nacos.UseServiceDiscovery(sc)
		}).
		OnApplicationLifeEvent(func(life *abstractions.ApplicationLife) {
			go func() {
				for {
					select {
					case ev := <-life.ApplicationStarted:
						xlog.GetXLogger("Application").Info("Application Started: %v", ev.Data)
					case ev := <-life.ApplicationStopped:
						xlog.GetXLogger("Application").Info("Application Stopped: %v", ev.Data)
					}
				}
			}()
		})
}

func registerEndpoints(rb router.IRouterBuilder) {
	// Health check
	endpoints.UseHealth(rb)
	endpoints.UseViz(rb)
	endpoints.UsePrometheus(rb)
	endpoints.UsePprof(rb)
	endpoints.UseReadiness(rb)
	endpoints.UseLiveness(rb)

	// JWT endpoint
	endpoints.UseJwt(rb)
}
`

const UserController_Tel = `package controller

import (
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/mvc"
	"{{.ModelName}}/internal/service"
)

type UserController struct {
	mvc.ApiController
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		ApiController: *mvc.NewApiController(),
		userService:   userService,
	}
}

// GetUserRequest request object
type GetUserRequest struct {
	` + "`mvc.RequestGET route:\"/api/users/:id\"`" + `
	Id int64 ` + "`param:\"id\"`" + `
}

// CreateUserRequest request object
type CreateUserRequest struct {
	` + "`mvc.RequestBody route:\"/api/users\" doc:\"Create User\"`" + `
	UserName string ` + "`json:\"username\" doc:\"Username\"`" + `
	Password string ` + "`json:\"password\" doc:\"Password\"`" + `
	Email    string ` + "`json:\"email\" doc:\"Email\"`" + `
	Phone    string ` + "`json:\"phone\" doc:\"Phone\"`" + `
}

// UpdateUserRequest request object
type UpdateUserRequest struct {
	` + "`mvc.RequestPUT route:\"/api/users/:id\"`" + `
	Id       int64  ` + "`param:\"id\"`" + `
	UserName string ` + "`json:\"username\" doc:\"Username\"`" + `
	Email    string ` + "`json:\"email\" doc:\"Email\"`" + `
	Phone    string ` + "`json:\"phone\" doc:\"Phone\"`" + `
}

// LoginRequest request object
type LoginRequest struct {
	` + "`mvc.RequestPOST route:\"/api/users/login\" doc:\"User Login\"`" + `
	UserName string ` + "`json:\"username\" doc:\"Username\"`" + `
	Password string ` + "`json:\"password\" doc:\"Password\"`" + `
}

// GetUser gets user by ID
func (c *UserController) GetUser(ctx *context.HttpContext, req *GetUserRequest) mvc.ApiResult {
	user, err := c.userService.GetById(req.Id)
	if err != nil {
		return c.Error(err)
	}
	return c.OK(user)
}

// CreateUser creates a new user
func (c *UserController) CreateUser(req *CreateUserRequest) mvc.ApiResult {
	user, err := c.userService.Create(req)
	if err != nil {
		return c.Error(err)
	}
	return c.Created(user)
}

// UpdateUser updates user info
func (c *UserController) UpdateUser(req *UpdateUserRequest) mvc.ApiResult {
	user, err := c.userService.Update(req)
	if err != nil {
		return c.Error(err)
	}
	return c.OK(user)
}

// DeleteUser deletes a user
func (c *UserController) DeleteUser(ctx *context.HttpContext, req *GetUserRequest) mvc.ApiResult {
	err := c.userService.Delete(req.Id)
	if err != nil {
		return c.Error(err)
	}
	return c.OK(nil)
}

// ListUsers lists all users
func (c *UserController) ListUsers(ctx *context.HttpContext) mvc.ApiResult {
	users, err := c.userService.ListAll()
	if err != nil {
		return c.Error(err)
	}
	return c.OK(users)
}

// Login user login
func (c *UserController) Login(req *LoginRequest) mvc.ApiResult {
	token, err := c.userService.Login(req.UserName, req.Password)
	if err != nil {
		return c.Error(err)
	}
	return c.OK(map[string]interface{}{
		"token": token,
	})
}
`

const UserService_Tel = `package service

import (
	"errors"
	"{{.ModelName}}/internal/model"
	"{{.ModelName}}/internal/repository"
	"github.com/liangboceo/yuanboot/utils/jwt"
	"time"
)

// IUserService user service interface
type IUserService interface {
	GetById(id int64) (*model.User, error)
	Create(req *CreateUserRequest) (*model.User, error)
	Update(req *UpdateUserRequest) (*model.User, error)
	Delete(id int64) error
	ListAll() ([]*model.User, error)
	Login(username, password string) (string, error)
}

// UserService user service implementation
type UserService struct {
	userRepo repository.IUserRepository
}

// NewUserService creates user service
func NewUserService(userRepo repository.IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUserRequest create user request
type CreateUserRequest struct {
	UserName string
	Password string
	Email    string
	Phone    string
}

// UpdateUserRequest update user request
type UpdateUserRequest struct {
	Id       int64
	UserName string
	Email    string
	Phone    string
}

func (s *UserService) GetById(id int64) (*model.User, error) {
	return s.userRepo.FindById(id)
}

func (s *UserService) Create(req *CreateUserRequest) (*model.User, error) {
	// Check if username exists
	existing, _ := s.userRepo.FindByUsername(req.UserName)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	user := &model.User{
		UserName: req.UserName,
		Password: req.Password, // Should be encrypted in production
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	return s.userRepo.Create(user)
}

func (s *UserService) Update(req *UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindById(req.Id)
	if err != nil {
		return nil, err
	}

	user.UserName = req.UserName
	user.Email = req.Email
	user.Phone = req.Phone

	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id int64) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) ListAll() ([]*model.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if user.Password != password {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT Token
	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.UserName,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token, err := jwt.NewToken(claims, "your-secret-key").SignedString()
	if err != nil {
		return "", err
	}

	return token, nil
}

// IUserRepository user repository interface
type IUserRepository interface {
	FindById(id int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll() ([]*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id int64) error
}

// NewUserRepository creates user repository
func NewUserRepository() IUserRepository {
	// Can get datasource through dependency injection
	return &UserRepository{}
}

// UserRepository user repository implementation
type UserRepository struct{}

func (r *UserRepository) FindById(id int64) (*model.User, error) {
	// TODO: Implement database query
	return &model.User{
		Id:       id,
		UserName: "test",
		Email:    "test@example.com",
	}, nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	// TODO: Implement database query
	return nil, errors.New("user not found")
}

func (r *UserRepository) FindAll() ([]*model.User, error) {
	// TODO: Implement database query
	return []*model.User{}, nil
}

func (r *UserRepository) Create(user *model.User) (*model.User, error) {
	// TODO: Implement database insert
	user.Id = time.Now().Unix()
	return user, nil
}

func (r *UserRepository) Update(user *model.User) (*model.User, error) {
	// TODO: Implement database update
	return user, nil
}

func (r *UserRepository) Delete(id int64) error {
	// TODO: Implement database delete
	return nil
}
`

const UserRepository_Tel = `package repository

import (
	"{{.ModelName}}/internal/model"
)

// IUserRepository user repository interface
type IUserRepository interface {
	FindById(id int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindAll() ([]*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id int64) error
}
`

const UserModel_Tel = `package model

import "time"

// User user model
type User struct {
	Id        int64     ` + "`" + `json:"id"` + "`" + `
	UserName  string    ` + "`" + `json:"username"` + "`" + `
	Password  string    ` + "`" + `json:"-"` + "`" + `
	Email     string    ` + "`" + `json:"email"` + "`" + `
	Phone     string    ` + "`" + `json:"phone"` + "`" + `
	Status    int       ` + "`" + `json:"status"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

// UserDTO user data transfer object
type UserDTO struct {
	Id       int64  ` + "`" + `json:"id"` + "`" + `
	UserName string ` + "`" + `json:"username"` + "`" + `
	Email    string ` + "`" + `json:"email"` + "`" + `
	Phone    string ` + "`" + `json:"phone"` + "`" + `
	Status   int    ` + "`" + `json:"status"` + "`" + `
}

// ToDTO convert to DTO
func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		Id:       u.Id,
		UserName: u.UserName,
		Email:    u.Email,
		Phone:    u.Phone,
		Status:   u.Status,
	}
}
`

const Middleware_Tel = `package middleware

import (
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/middlewares"
)

// NewRecovery creates recovery middleware
func NewRecovery() func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
	return middlewares.NewRecovery()
}

// NewLogger creates logger middleware
func NewLogger() func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
	return middlewares.NewLogger()
}

// NewCORS creates CORS middleware
func NewCORS() func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
	return middlewares.NewCORS()
}

// NewRequestTracker creates request tracker middleware
func NewRequestTracker() func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
	return middlewares.NewRequestTracker()
}

// NewJWT creates JWT auth middleware
func NewJWT(secret string) func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
	return func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
		return func(ctx *context.HttpContext) {
			token := ctx.Header("Authorization")
			if token == "" {
				ctx.JSON(401, context.H{"error": "unauthorized"})
				return
			}
			// JWT validation logic
			next(ctx)
		}
	}
}

// CustomMiddleware custom middleware example
func CustomMiddleware(param string) func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
	return func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
		return func(ctx *context.HttpContext) {
			// Pre-processing
			ctx.SetItem("custom_param", param)

			// Call next handler
			next(ctx)

			// Post-processing
		}
	}
}
`

const ConfigDev_Tel = `yuanboot:
  application:
    name: {{.ModelName}}
    metadata: "develop"
    server:
      type: "fasthttp"
      address: ":8080"
      path: ""
      max_request_size: 2097152
      session:
        name: "SESSION_ID"
        timeout: 3600
      jwt:
        header: "Authorization"
        secret: "your-dev-secret-key-change-in-production"
        prefix: "Bearer"
        expires: 24
        enable: true
        skip_path:
          - "/info"
          - "/api/users/login"
          - "/health"
          - "/metrics"
      cors:
        allow_origins:
          - "*"
        allow_methods:
          - "POST"
          - "GET"
          - "PUT"
          - "DELETE"
          - "PATCH"
        allow_credentials: true
        allow_headers:
          - "*"
  datasource:
    db:
      name: default
      url: tcp(localhost:3306)/{{.ModelName}}?charset=utf8mb4&parseTime=True&loc=Local
      username: root
      password: root
      debug: true
      pool:
        init_cap: 5
        max_cap: 30
        idletimeout: 300
    pool:
      max_open_conns: 100
      max_idle_conns: 10
      conn_max_lifetime: 3600
    redis:
      default:
        addr: localhost:6379
        password: ""
        db: 0
        pool_size: 10
  servicediscovery:
    nacos:
      enabled: true
      server_addr: localhost:8848
      namespace: public
      group: DEFAULT_GROUP
      weight: 100
`

const ConfigProd_Tel = `yuanboot:
  application:
    name: {{.ModelName}}
    metadata: "production"
    server:
      type: "fasthttp"
      address: ":8080"
      path: ""
      max_request_size: 2097152
      session:
        name: "SESSION_ID"
        timeout: 3600
      jwt:
        header: "Authorization"
        secret: "${JWT_SECRET}"
        prefix: "Bearer"
        expires: 24
        enable: true
        skip_path:
          - "/health"
          - "/metrics"
      cors:
        allow_origins:
          - "https://your-domain.com"
        allow_methods:
          - "POST"
          - "GET"
          - "PUT"
          - "DELETE"
          - "PATCH"
        allow_credentials: true
  datasource:
    db:
      name: default
      url: tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local
      username: ${DB_USERNAME}
      password: ${DB_PASSWORD}
      debug: false
      pool:
        init_cap: 10
        max_cap: 100
        idletimeout: 600
    pool:
      max_open_conns: 200
      max_idle_conns: 20
      conn_max_lifetime: 7200
    redis:
      default:
        addr: ${REDIS_HOST}:${REDIS_PORT}
        password: ${REDIS_PASSWORD}
        db: 0
        pool_size: 20
  servicediscovery:
    nacos:
      enabled: true
      server_addr: ${NACOS_HOST}:${NACOS_PORT}
      namespace: ${NACOS_NAMESPACE}
      group: ${NACOS_GROUP}
      weight: 100
  mq:
    kafka:
      brokers:
        - ${KAFKA_BROKERS}
      topic: {{.ModelName}}
      group_id: {{.ModelName}}-group
`

const Mod_Tel = `
module {{.ModelName}}

go 1.18

require (
	github.com/liangboceo/dependencyinjection v1.0.0
	github.com/liangboceo/yuanboot {{.Version}}
)

require (
	github.com/valyala/fasthttp v1.51.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`

const Makefile_Tel = `
.PHONY: build run test clean docker

APP_NAME := {{.ModelName}}
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

LDFLAGS := -ldflags "-s -w"
BUILD_DIR := ./bin
DEV_CONFIG := config_dev.yml
PROD_CONFIG := config_prod.yml

build-linux:
	@echo "Building $(APP_NAME) for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/main.go

build-windows:
	@echo "Building $(APP_NAME) for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/main.go

build-darwin:
	@echo "Building $(APP_NAME) for macOS..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/main.go

build:
	@echo "Building $(APP_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go

run-dev:
	@echo "Running $(APP_NAME) in development mode..."
	$(GOCMD) run ./cmd/main.go --config=$(DEV_CONFIG)

run-prod:
	@echo "Running $(APP_NAME) in production mode..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go && $(BUILD_DIR)/$(APP_NAME) --config=$(PROD_CONFIG)

test:
	@echo "Running tests..."
	$(GOTEST) -v -cover ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

lint:
	@echo "Linting code..."
	golangci-lint run ./...

docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(APP_NAME):latest

help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-darwin  - Build for macOS"
	@echo "  run-dev       - Run in development mode"
	@echo "  run-prod      - Run in production mode"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build files"
	@echo "  deps          - Download dependencies"
	@echo "  verify        - Verify dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
`

const Readme_Tel = `
# {{.ModelName}}

Microservice application based on Yuanboot framework

## Features

- RESTful API design
- MVC architecture
- JWT authentication
- Service discovery (Nacos)
- Complete middleware support
- Database and cache integration

## Project Structure

` + "```" + `
{{.ModelName}}/
├── cmd/                    # Application entry
│   └── main.go
├── internal/               # Internal packages
│   ├── controller/         # Controller layer
│   ├── service/            # Service layer
│   ├── repository/         # Data access layer
│   ├── model/              # Data model
│   └── middleware/         # Middleware
├── config/                 # Configuration files
│   ├── config_dev.yml      # Development environment
│   └── config_prod.yml     # Production environment
├── go.mod
├── go.sum
└── Makefile
` + "```" + `

## Quick Start

### Requirements

- Go 1.18+
- MySQL 5.7+
- Redis 6.0+

### Install dependencies

` + "```bash" + `
make deps
` + "```" + `

### Run development version

` + "```bash" + `
make run-dev
` + "```" + `

### Build

` + "```bash" + `
# Local build
make build

# Cross compile
make build-linux
make build-windows
make build-darwin
` + "```" + `

### Test

` + "```bash" + `
make test
make test-coverage
` + "```" + `

## API Endpoints

### User Management

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/users/register | User registration |
| POST | /api/users/login | User login |
| GET | /api/users/:id | Get user info |
| PUT | /api/users/:id | Update user info |
| DELETE | /api/users/:id | Delete user |
| GET | /api/users | Get user list |

### Operations

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /metrics | Prometheus metrics |
| GET | /pprof | Performance analysis |

## Configuration

For detailed configuration, please refer to [Yuanboot Documentation](https://yuanboot.star2cloud.com)

## License

MIT License
`
