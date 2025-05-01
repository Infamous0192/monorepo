package telegram

import (
	"app/pkg/exception"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// MaxAuthAge is the maximum age of auth data in seconds
	MaxAuthAge = 86400 // 24 hours
	// JWTExpiration is the expiration time for JWT tokens
	JWTExpiration = 24 * time.Hour
)

// TelegramUser represents authenticated Telegram user data
type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	AuthDate  int64  `json:"auth_date"`
}

// UserProvider defines the interface for custom user storage
type UserProvider interface {
	// GetOrCreateUser creates or updates a user from Telegram data
	GetOrCreateUser(telegramUser *TelegramUser) (any, error)
	// GetUserRole returns the user's role for JWT claims
	GetUserRole(telegramUser *TelegramUser) string
}

// Claims represents JWT claims for Telegram user
type Claims struct {
	UserID     int64  `json:"user_id"`
	FirstName  string `json:"first_name"`
	Username   string `json:"username,omitempty"`
	Role       string `json:"role"`
	DomainUser any    `json:"domain_user,omitempty"`
	jwt.RegisteredClaims
}

// AuthService handles Telegram authentication operations
type AuthService struct {
	jwtSecret    []byte
	botToken     string
	userProvider UserProvider // Add user provider
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(botToken string, jwtSecret string, userProvider UserProvider) *AuthService {
	return &AuthService{
		botToken:     botToken,
		jwtSecret:    []byte(jwtSecret),
		userProvider: userProvider,
	}
}

// ValidateInitData validates Telegram Web App initData string
func (s *AuthService) ValidateInitData(initData string) (*TelegramUser, error) {
	// Parse the initData string
	params, err := url.ParseQuery(initData)
	if err != nil {
		return nil, exception.BadRequest("Invalid initData format")
	}

	// Get and validate hash
	hash := params.Get("hash")
	if hash == "" {
		return nil, exception.BadRequest("Missing hash")
	}
	params.Del("hash")

	// Sort params
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create data check string
	var dataCheckParts []string
	for _, k := range keys {
		dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckParts, "\n")

	// Validate hash
	if !s.validateHash(dataCheckString, hash) {
		return nil, exception.Http(401, "Invalid hash")
	}

	// Parse user data
	var user TelegramUser
	if userData := params.Get("user"); userData != "" {
		if err := json.Unmarshal([]byte(userData), &user); err != nil {
			return nil, exception.BadRequest("Invalid user data")
		}
	}

	// Validate auth date
	if authDate, err := strconv.ParseInt(params.Get("auth_date"), 10, 64); err == nil {
		user.AuthDate = authDate
		if time.Now().Unix()-authDate > MaxAuthAge {
			return nil, exception.BadRequest("Authentication data is too old")
		}
	} else {
		return nil, exception.BadRequest("Invalid auth_date")
	}

	return &user, nil
}

// validateHash validates the Telegram hash using bot token
func (s *AuthService) validateHash(dataCheckString, hash string) bool {
	// Create secret key using bot token
	secretKeyHmac := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyHmac.Write([]byte(s.botToken))
	secretKey := secretKeyHmac.Sum(nil)

	// Calculate hash
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(hash), []byte(expectedHash))
}

// GenerateToken generates a JWT token for the authenticated user
func (s *AuthService) GenerateToken(user *TelegramUser) (string, error) {
	var domainUser any
	var role string = "user" // Default role

	// If user provider is configured, get domain user and role
	if s.userProvider != nil {
		var err error
		domainUser, err = s.userProvider.GetOrCreateUser(user)
		if err != nil {
			return "", exception.InternalError("Failed to get or create user: " + err.Error())
		}
		role = s.userProvider.GetUserRole(user)
	}

	now := time.Now()
	claims := &Claims{
		UserID:     user.ID,
		FirstName:  user.FirstName,
		Username:   user.Username,
		Role:       role,       // Add role from provider
		DomainUser: domainUser, // Add domain user data
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates and parses a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, exception.Http(401, "Invalid token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, exception.Http(401, "Invalid token claims")
}
