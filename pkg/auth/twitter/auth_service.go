package twitter

import (
	"app/pkg/exception"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// Claims represents the JWT claims for Twitter authentication
type Claims struct {
	UserID       string `json:"user_id"`       // Twitter user ID
	Username     string `json:"username"`      // Twitter username
	DisplayName  string `json:"display_name"`  // Twitter display name
	ProfileImage string `json:"profile_image"` // Twitter profile image
	Role         string `json:"role"`          // User role from provider
	DomainUser   any    `json:"domain_user"`   // Application user data
	jwt.RegisteredClaims
}

// TwitterUser represents Twitter user data
type TwitterUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	ProfileURL  string `json:"profile_image_url"`
	Description string `json:"description"`
}

// UserProvider defines the interface for application user storage
type UserProvider interface {
	// GetOrCreateUser creates or updates a user from Twitter data
	GetOrCreateUser(twitterUser *TwitterUser) (any, error)
	// GetUserRole returns the user's role for JWT claims
	GetUserRole(twitterUser *TwitterUser) string
}

type AuthService struct {
	config       *oauth2.Config
	jwtSecret    string
	userProvider UserProvider
}

// NewAuthService creates a new Twitter authentication service
func NewAuthService(clientID, clientSecret, redirectURL, jwtSecret string, userProvider UserProvider) *AuthService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://x.com/i/oauth2/authorize",
			TokenURL: "https://api.x.com/2/oauth2/token",
		},
		Scopes: []string{
			"tweet.read",
			"users.read",
			"offline.access",
		},
	}

	return &AuthService{
		config:       config,
		jwtSecret:    jwtSecret,
		userProvider: userProvider,
	}
}

// GetAuthURL returns the Twitter OAuth2 authorization URL
func (s *AuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

// ExchangeCode exchanges the OAuth2 code for tokens
func (s *AuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, exception.InternalError("Failed to exchange code: " + err.Error())
	}
	return token, nil
}

// GetUserInfo fetches the user information from Twitter API
func (s *AuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*TwitterUser, error) {
	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://api.x.com/2/users/me?user.fields=profile_image_url,description")
	if err != nil {
		return nil, exception.InternalError("Failed to fetch user info: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, exception.InternalError(fmt.Sprintf("Twitter API error: %s", string(body)))
	}

	var response struct {
		Data TwitterUser `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, exception.InternalError("Failed to decode user info: " + err.Error())
	}

	return &response.Data, nil
}

// GenerateToken generates a JWT token for the authenticated user
func (s *AuthService) GenerateToken(user *TwitterUser) (string, error) {
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

	claims := Claims{
		UserID:       user.ID,
		Username:     user.Username,
		DisplayName:  user.Name,
		ProfileImage: user.ProfileURL,
		Role:         role,
		DomainUser:   domainUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "twitter-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", exception.InternalError("Failed to generate token: " + err.Error())
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, exception.Http(401, "Invalid token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, exception.Http(401, "Invalid token claims")
}
