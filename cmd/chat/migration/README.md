# MongoDB Migrations for Chat Service

This directory contains MongoDB migration code for the chat service. The migrations are implemented in Go, leveraging the existing repository implementations.

## Migration Files

- `migrate.go`: Main migration program that uses repository implementations to set up collections and indexes

## How to Run Migrations

From the project root, run:

```bash
# For local development
go run cmd/chat/migrations/migrate.go

# For specific environment, set environment variables:
MONGODB_HOST=host \
MONGODB_PORT=27017 \
MONGODB_DATABASE=wonderverse_chat \
MONGODB_USERNAME=user \
MONGODB_PASSWORD=pass \
go run cmd/chat/migrations/migrate.go
```

## What Gets Created

1. `users` collection:
   - Indexes: userId (unique), username (unique)
   - Schema validation via Go struct tags

2. `clients` collection:
   - Indexes: clientKey (unique), name
   - Schema validation via Go struct tags

3. `chatrooms` collection:
   - Indexes: participants+timestamp, type+timestamp
   - Schema validation via Go struct tags

4. `chats` collection:
   - Indexes: chatroom+timestamp, sender+timestamp, receiver+timestamp
   - Schema validation via Go struct tags

## Benefits of Go-based Migration

1. Reuses existing repository code
2. Type safety through Go structs
3. Consistent with codebase
4. Handles configuration the same way as the main application
5. Can be extended to handle versioned migrations if needed

## Rollback

To rollback migrations, you can use the MongoDB shell:

```javascript
db.users.drop()
db.clients.drop()
db.chatrooms.drop()
db.chats.drop()
```

Note: Be extremely careful with rollbacks in production. Always backup data first. 