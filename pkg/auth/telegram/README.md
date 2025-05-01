# Telegram Authentication Package

This package provides JWT-based authentication using Telegram Web App, with support for application user storage integration.

## Features
- Telegram Web App authentication with hash validation
- JWT token generation and validation
- Application user storage integration
- Role-based access control
- Secure cookie handling

## Installation

```bash
go get github.com/golang-jwt/jwt/v5
```

## Basic Usage

```go
// Initialize the service
authService := telegram.NewAuthService(
    os.Getenv("TELEGRAM_BOT_TOKEN"),
    os.Getenv("JWT_SECRET"),
    nil, // No user provider
)

// Create handler and middleware
authHandler := telegram.NewAuthHandler(authService)
authMiddleware := telegram.NewAuthMiddleware(authService)

// Register routes
authHandler.RegisterRoutes(app)
```

## Domain User Integration

### 1. Implement UserProvider Interface

```go
// Define your domain user model
type AppUser struct {
    ID         string `json:"id"`
    Email      string `json:"email"`
    Role       string `json:"role"`
    TelegramID int64  `json:"telegram_id"`
}

// Define your storage interface
type UserStorage interface {
    FindUserByTelegramID(telegramID int64) (*AppUser, error)
    CreateUser(user *AppUser) error
    UpdateUser(user *AppUser) error
}

// Implement UserProvider
type AppUserProvider struct {
    storage UserStorage
}

func (p *AppUserProvider) GetOrCreateUser(telegramUser *telegram.TelegramUser) (interface{}, error) {
    // Check if user exists
    user, err := p.storage.FindUserByTelegramID(telegramUser.ID)
    if err == nil {
        return user, nil
    }

    // Create new user
    user = &AppUser{
        ID:         uuid.New().String(),
        TelegramID: telegramUser.ID,
        Role:       "user",
    }
    
    if err := p.storage.CreateUser(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (p *AppUserProvider) GetUserRole(telegramUser *telegram.TelegramUser) string {
    user, err := p.storage.FindUserByTelegramID(telegramUser.ID)
    if err != nil {
        return "user" // Default role
    }
    return user.Role
}
```

### 2. Initialize with User Provider

```go
// Initialize your storage
storage := NewUserStorage()

// Create user provider
userProvider := &AppUserProvider{storage: storage}

// Initialize auth service with provider
authService := telegram.NewAuthService(
    os.Getenv("TELEGRAM_BOT_TOKEN"),
    os.Getenv("JWT_SECRET"),
    userProvider,
)
```

### 3. Access Application User Data

```go
func protectedRoute(c *fiber.Ctx) error {
    // Get claims from context
    claims := c.Locals("user").(*telegram.Claims)
    
    // Access Telegram data
    telegramID := claims.UserID
    username := claims.Username
    
    // Access application user data
    appUser := claims.DomainUser.(*AppUser)
    userID := appUser.ID
    role := claims.Role
    
    return c.JSON(fiber.Map{
        "telegram_id": telegramID,
        "user_id": userID,
        "role": role,
    })
}
```

## JWT Claims Structure

```go
type Claims struct {
    UserID     int64       `json:"user_id"`      // Telegram user ID
    FirstName  string      `json:"first_name"`   // Telegram first name
    Username   string      `json:"username"`      // Telegram username
    Role       string      `json:"role"`         // User role from provider
    DomainUser interface{} `json:"domain_user"`  // Application user data
    jwt.RegisteredClaims
}
```

## Security Considerations

1. **Bot Token Security**
   - Store bot token securely using environment variables
   - Never commit bot token to version control

2. **JWT Security**
   - Use a strong, random JWT secret
   - Set appropriate token expiration
   - Validate tokens on every request

3. **User Data**
   - Validate all user data before storage
   - Implement proper role management
   - Handle user updates securely

4. **Error Handling**
   - Implement proper error logging
   - Return appropriate HTTP status codes
   - Don't expose internal errors to clients

## Example Implementation

See `pkg/auth/examples/domain/` for a complete example of application user integration. 