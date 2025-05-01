# Telegram Authentication Example

This example demonstrates how to implement Telegram Web App authentication with JWT tokens using our `pkg/auth/telegram` package.

## Features
- Telegram Web App authentication
- JWT token generation and validation
- Role-based access control
- Protected routes
- CORS configuration

## Backend Setup

1. Set environment variables:
```bash
export JWT_SECRET="your-secret-key"
```

2. Run the example:
```bash
go run main.go
```

The server starts on `http://localhost:3000` with endpoints:
- `POST /auth/telegram/login` - Login with Telegram
- `GET /api/protected` - Requires authentication
- `GET /api/admin` - Requires admin role

## Frontend Integration

1. Add Telegram Web App script to your HTML:
```html
<!DOCTYPE html>
<html>
<head>
    <title>Telegram Auth Example</title>
    <script src="https://telegram.org/js/telegram-web-app.js"></script>
</head>
<body>
    <div id="app">
        <button onclick="login()">Login with Telegram</button>
        <button onclick="fetchProtected()">Access Protected Route</button>
    </div>
    <pre id="output"></pre>
</body>
</html>
```

2. Add JavaScript code:
```javascript
// Initialize Telegram Web App
const initData = window.Telegram.WebApp.initData;

// Login with Telegram
async function login() {
    try {
        const response = await fetch('http://localhost:3000/auth/telegram/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ initData }),
        });

        const data = await response.json();
        if (!data.data?.token) {
            throw new Error('Login failed');
        }

        // Store token
        localStorage.setItem('token', data.data.token);
        document.getElementById('output').textContent = 'Logged in as: ' + 
            JSON.stringify(data.data.user, null, 2);
    } catch (error) {
        document.getElementById('output').textContent = 'Login failed: ' + error.message;
    }
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
```

## Testing with cURL

```bash
# Login (replace with actual initData)
curl -X POST http://localhost:3000/auth/telegram/login \
  -H "Content-Type: application/json" \
  -d '{"initData":"your-init-data"}'

# Access protected route
curl http://localhost:3000/api/protected \
  -H "Authorization: Bearer your-jwt-token"

# Access admin route
curl http://localhost:3000/api/admin \
  -H "Authorization: Bearer your-jwt-token"
```

## Security Notes

1. Always use HTTPS in production
2. Use a strong JWT secret
3. Validate all Telegram Web App data
4. Implement proper role management
5. Add rate limiting
6. Configure CORS for your domains only

## Telegram Web App Setup

1. Create a bot using [@BotFather](https://t.me/botfather)
2. Enable Web App in bot settings
3. Set your domain in the bot's menu button
4. Configure allowed domains for your Web App

For more details, see [Telegram Web Apps documentation](https://core.telegram.org/bots/webapps) 