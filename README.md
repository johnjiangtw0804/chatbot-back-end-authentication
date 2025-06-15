
# This is the secured backend authentication service

This is a secure authentication backend service designed for a chatbot frontend application.
It supports user registration, login, and profile management.

## How to run?

Setup .env / config.env
Edit config.env and set your local DB settings:
```
DB_USER=root
DB_PASSWORD=root
DB_NAME=chat-bot
DB_HOST=localhost
DB_PORT=5432
APP_ENV=development
```

🐘 Database (PostgreSQL via Docker)
Start PostgreSQL container:
```
make start-pg
```

Run DB migration
```
make migrate up
```

Start the server
```
make run
```

Server will start at:
http://localhost:8080

