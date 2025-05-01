#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create extensions
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    -- Create schemas if they don't exist
    CREATE SCHEMA IF NOT EXISTS quiz;
    CREATE SCHEMA IF NOT EXISTS article;

    -- Grant privileges to the default user
    GRANT ALL PRIVILEGES ON SCHEMA quiz TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON SCHEMA article TO $POSTGRES_USER;
EOSQL 