package auth

import (
	"app/pkg/exception"
	"app/pkg/quiz/config"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents JWT claims for user authentication
type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles user authentication operations
type AuthService interface {
	// Register creates a new user account
	Register(userDTO *entity.UserDTO) (*entity.User, error)

	// Login authenticates a user and returns a JWT token
	Login(username, password string) (string, error)

	// ValidateToken validates a JWT token and returns the claims
	ValidateToken(tokenString string) (*Claims, error)

	// HashPassword creates a bcrypt hash of the password
	HashPassword(password string) (string, error)

	// CheckPassword compares a password with a hash
	CheckPassword(password, hash string) error
}

// authService implements the AuthService interface
type authService struct {
	userRepo repository.UserRepository
	config   *config.AuthConfig
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, config *config.AuthConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

// Register creates a new user account
func (s *authService) Register(userDTO *entity.UserDTO) (*entity.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(userDTO.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, exception.InvalidPayload(map[string]string{
			"username": "Username already exists",
		})
	}

	// Hash the password
	hashedPassword, err := s.HashPassword(userDTO.Password)
	if err != nil {
		return nil, exception.InternalError("Failed to hash password")
	}

	// Create user entity
	user := &entity.User{
		Name:      userDTO.Name,
		Username:  userDTO.Username,
		Password:  hashedPassword,
		Role:      userDTO.Role,
		Status:    userDTO.Status,
		BirthDate: userDTO.BirthDate,
	}

	// Save user to repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(username, password string) (string, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", exception.InvalidPayload(map[string]string{
			"username": "Username not found",
		})
	}

	// Check password
	err = s.CheckPassword(password, user.Password)
	if err != nil {
		return "", exception.InvalidPayload(map[string]string{
			"password": "Incorrect password",
		})
	}

	// Check if user is active
	if !user.Status {
		return "", exception.Http(403, "Account is inactive")
	}

	// Generate token
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", exception.InternalError("Failed to generate token")
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, exception.Http(401, "Invalid token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, exception.Http(401, "Invalid token claims")
}

// HashPassword creates a bcrypt hash of the password
func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.PasswordHashCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func (s *authService) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
