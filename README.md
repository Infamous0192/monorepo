# Quiz API

A RESTful API for managing quizzes, questions, answers, and submissions.

## Features

- Quiz management
- Question management with options and answers
- Submission handling with bulk submission support
- RESTful API with proper error handling
- Database migrations
- Authentication and authorization
  - Local username/password authentication
  - Role-based access control
  - JWT-based authentication
  - Support for multiple auth providers (extensible)

## Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- Git

## Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/quiz-api.git
cd quiz-api
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables
```bash
# Create a .env file from the example
cp .env.example .env
# Edit the .env file with your database credentials and JWT secret
```

## Running the application

1. Start the application
```bash
go run ./cmd/quiz
```

The server will start on port 8080 (or the port specified in your .env file).

## Authentication

The API supports authentication using JWT tokens. To authenticate:

1. Register a new user:
```bash
curl -X POST http://localhost:8080/auth/local/register \
  -H "Content-Type: application/json" \
  -d '{"username": "user1", "password": "password123", "email": "user@example.com", "fullName": "User One"}'
```

2. Login to get a token:
```bash
curl -X POST http://localhost:8080/auth/local/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user1", "password": "password123"}'
```

3. Use the token for authenticated requests:
```bash
curl -X GET http://localhost:8080/api/v1/quizzes \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## API Endpoints

### Authentication
- `POST /auth/local/register` - Register a new user
- `POST /auth/local/login` - Login with username and password

### Quizzes
- `GET /api/v1/quizzes` - Get all quizzes (public)
- `GET /api/v1/quizzes/:id` - Get quiz by ID (public)
- `POST /api/v1/quizzes` - Create a new quiz (admin only)
- `PUT /api/v1/quizzes/:id` - Update a quiz (admin only)
- `DELETE /api/v1/quizzes/:id` - Delete a quiz (admin only)

### Questions
- `GET /api/v1/questions` - Get all questions (public)
- `GET /api/v1/questions/:id` - Get question by ID (public)
- `POST /api/v1/questions` - Create a new question (admin only)
- `PUT /api/v1/questions/:id` - Update a question (admin only)
- `DELETE /api/v1/questions/:id` - Delete a question (admin only)

### Answers
- `GET /api/v1/answers` - Get all answers (admin only)
- `GET /api/v1/answers/:id` - Get answer by ID (admin only)
- `POST /api/v1/answers` - Create a new answer (admin only)
- `PUT /api/v1/answers/:id` - Update an answer (admin only)
- `DELETE /api/v1/answers/:id` - Delete an answer (admin only)

### Submissions
- `GET /api/v1/submissions` - Get all submissions (admin only)
- `GET /api/v1/submissions/:id` - Get submission by ID (authenticated user)
- `POST /api/v1/submissions` - Create a new submission (authenticated user)
- `POST /api/v1/submissions/bulk` - Create multiple submissions (authenticated user)
- `PUT /api/v1/submissions/:id` - Update a submission (admin only)
- `DELETE /api/v1/submissions/:id` - Delete a submission (admin only)

## Project Structure

```
.
├── cmd
│   └── quiz            # Application entry point
├── pkg
│   ├── auth            # Authentication packages
│   │   ├── domain      # Domain entities and repositories for auth
│   │   ├── local       # Username/password authentication
│   │   ├── telegram    # Telegram authentication
│   │   └── repository  # Repository implementations
│   └── quiz
│       ├── domain      # Domain entities and repositories
│       ├── handlers    # HTTP handlers
│       ├── repository  # Repository implementations
│       └── services    # Business logic
```

## Database Schema

The application uses the following entities:
- User (for authentication)
- Quiz
- Question
- Answer
- Option (for multiple choice questions)
- Submission

## License

MIT 