# Telegram Service Migration

This directory contains database migration scripts for the Telegram service.

## Overview

The migration script creates initial data for development and testing purposes, including:
- Sample bot clients
- Sample users

## Usage

To run the migration:

```bash
go run cmd/telegram/migration/migrate.go --config=cmd/telegram/config/config.yml
```

By default, it will look for the config file at `cmd/telegram/config/config.yml`. You can specify a different path using the `--config` flag.

## Environment Variables

If no config file is found, the script will attempt to read configuration from environment variables:

- `MONGODB_HOST` - MongoDB host (default: localhost)
- `MONGODB_PORT` - MongoDB port (default: 27017)
- `MONGODB_DATABASE` - MongoDB database name (default: database)
- `MONGODB_USERNAME` - MongoDB username (optional)
- `MONGODB_PASSWORD` - MongoDB password (optional)

## Sample Data

The migration creates the following sample data:

### Clients (Bots)
1. Support Bot
   - Username: test_bot_1
   - Type: support
   - Status: active
   - Allowed Updates: message, callback_query

2. Notification Bot
   - Username: test_bot_2
   - Type: notification
   - Status: active
   - Allowed Updates: message, callback_query, channel_post

### Users
- User 1: John Doe (TelegramID: 123456789)
- User 2: Jane Smith (TelegramID: 987654321) 