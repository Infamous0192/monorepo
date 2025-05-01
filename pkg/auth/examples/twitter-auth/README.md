# Twitter OAuth2 Authentication Example

This example demonstrates how to implement Twitter OAuth2 authentication with JWT tokens using our `pkg/auth/twitter` package.

## Features
- Twitter OAuth2 authentication flow
- JWT token generation and validation
- Role-based access control
- Protected routes
- CORS configuration

## Prerequisites

1. Create a Twitter Developer Account and App:
   - Go to [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard)
   - Create a new app or use an existing one
   - Enable OAuth2 and configure callback URL
   - Get your Client ID and Client Secret

## Backend Setup

1. Set environment variables:
```bash
export TWITTER_CLIENT_ID="your-client-id"
export TWITTER_CLIENT_SECRET="your-client-secret"
export TWITTER_REDIRECT_URL="http://localhost:3000/auth/twitter/callback"
export JWT_SECRET="your-jwt-secret"
```

2. Run the example:
```bash
go run main.go
```

The server starts on `http://localhost:3000` with endpoints:
- `GET /auth/twitter/login` - Start OAuth2 flow
- `GET /auth/twitter/callback` - OAuth2 callback
- `GET /api/protected` - Requires authentication
- `GET /api/admin` - Requires admin role

## Frontend Integration

1. Add Twitter OAuth2 button to your HTML:
```html
<!DOCTYPE html>
<html>
<head>
    <title>Twitter Auth Example</title>
</head>
<body>
    <div id="app">
        <button onclick="login()">Login with Twitter</button>
        <button onclick="fetchProtected()">Access Protected Route</button>
    </div>
    <pre id="output"></pre>

    <script>
        // Start OAuth2 flow
        function login() {
            window.location.href = 'http://localhost:3000/auth/twitter/login';
        }

        // Access protected route
        async function fetchProtected() {
            const token = localStorage.getItem('token');
            if (!token) {
                document.getElementById('output').textContent = 'Not authenticated';
                return;
            }

            try {
                const response = await fetch('http://localhost:3000/api/protected', {
                    headers: {
                        'Authorization': 'Bearer ' + token,
                    },
                });

                const data = await response.json();
                document.getElementById('output').textContent = 
                    'Protected data: ' + JSON.stringify(data, null, 2);
            } catch (error) {
                document.getElementById('output').textContent = 'Error: ' + error.message;
            }
        }

        // Handle OAuth callback
        if (window.location.pathname === '/callback') {
            const urlParams = new URLSearchParams(window.location.search);
            const token = urlParams.get('token');
            if (token) {
                localStorage.setItem('token', token);
                window.location.href = '/'; // Redirect to home
            }
        }
    </script>
</body>
</html>
```

## Testing with cURL

```bash
# Start OAuth2 flow (opens browser)
curl -i http://localhost:3000/auth/twitter/login

# Access protected route (after getting token)
curl http://localhost:3000/api/protected \
  -H "Authorization: Bearer your-jwt-token"

# Access admin route
curl http://localhost:3000/api/admin \
  -H "Authorization: Bearer your-jwt-token"
```

## Security Notes

1. Always use HTTPS in production
2. Use secure environment variables
3. Never commit secrets to version control
4. Use proper session management
5. Implement proper role management
6. Add rate limiting
7. Configure CORS for your domains only

## Twitter OAuth2 Setup

1. Configure OAuth2 settings in Twitter Developer Portal:
   - Set callback URL to `http://localhost:3000/auth/twitter/callback`
   - Enable OAuth2 with PKCE
   - Add required scopes: `tweet.read`, `users.read`, `offline.access`

2. Security best practices:
   - Use state parameter to prevent CSRF
   - Validate all tokens and claims
   - Store tokens securely
   - Implement proper error handling

For more details, see:
- [Twitter OAuth2 Documentation](https://developer.twitter.com/en/docs/authentication/oauth-2-0)
- [Twitter API Reference](https://developer.twitter.com/en/docs/api-reference-index) 