
set quiet := true
set dotenv-load

MAIN_PACKAGE_PATH := "."
BINARY_NAME := "aggregator"

# =========================================================================== #
# HELPERS
# =========================================================================== #

# Display help
[private]
help:
  just --list --unsorted

# Confirm or cancel recipe
_confirm:
  echo "Are you sure? [y/N] " && read ans && [ ${ans:-N} = y ]

# =========================================================================== #
# DEVELOPMENT
# =========================================================================== #

# Run dev server
run: 
  go run .

# Start docker containers
up:
  docker compose up -d

# Log into Postgres using psql 
psql:
  psql $DSN

# Preform migration on DB
migrate-up version="": _confirm
  echo 'Running up migrations...'
  migrate -path ./sql/migrations -database $DSN up {{version}}

# Undo migrations on DB
migrate-down version="": 
  echo 'Running down migrations...'
  migrate -path ./sql/migrations -database $DSN down {{version}}

# Create new migration files
migrate-new name: 
  echo 'Creating migration files for {{name}}...'
  migrate create -seq -ext=.sql -dir=./sql/migrations {{name}}

migrate-force version: 
  migrate -path ./sql/migrations -database $DSN force {{version}} 

sqlc-generate:
  sqlc generate

# =========================================================================== #
# QUALITY CONTROL
# =========================================================================== #

# Verify and Vet all Go files in project
audit:
  echo 'Checking module dependencies'
  go mod tidy -diff
  go mod verify
  echo 'Vetting code...'
  go vet ./...
  echo 'Running tests...'
  go test -race -vet=off ./...

# Run formatter and tidy over all Go files in project
tidy:
  echo 'Formatting .go files...'
  go fmt ./...
  echo 'Tidying module dependencies...'
  go mod tidy -v


# =========================================================================== #
# BUILD
# =========================================================================== #

# Build for current OS/Arch
build:
  go build -ldflags='-s' -o=./tmp/{{BINARY_NAME}} {{MAIN_PACKAGE_PATH}}

# Build for current OS/Arch and then run the binary
preview: build
  ./tmp/{{BINARY_NAME}}

