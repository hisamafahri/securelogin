# Database Setup and Management

## First-Time Setup

1. Prerequisites
  - Docker
  - Go 1.25.0 or higher
  - [Atlas CLI](https://atlasgo.io/getting-started#installation) (`brew install ariga/tap/atlas`)
  - [`make`](https://en.wikipedia.org/wiki/Make_(software)) utility

2. Run PostgreSQL in Docker

```bash
# example
docker run -d \
  --name securelogin-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=securelogin \
  -p 5432:5432 \
  postgres:latest
```

3. Configure environment variables

Create a `.env` file in the project root:

```bash
# adjust the credentials, port, and database name as needed
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/securelogin
```

4. Apply migrations

```bash
make migrate-up
```

## Creating and Applying Database Changes

1. Define or update your model in `infrastructure/pgsql/models/`. Example: `infrastructure/pgsql/models/user.go`:

```go
package models

import "github.com/google/uuid"

type User struct {
    ID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Email string    `gorm:"type:varchar(255);not null;uniqueIndex"`
    Name  string    `gorm:"type:varchar(255);not null"`
}

func (User) TableName() string {
    return "users"
}
```

2. Register the model in `cmd/atlas/main.go`

```go
func main() {
    stmts, err := gormschema.New("postgres").Load(
        // ... other models ...
        &models.User{},  // Add your new model here
    )
}
```

3. Generate migration

```bash
make migrate-generate
# Enter migration name when prompted (e.g., add_users_table)
```

This creates a new migration file in `infrastructure/pgsql/migrations/` with:
- Unique prefix (e.g., `0002_add_users_table.sql`)
- Proper SQL for creating/altering tables
- Automatic checksum in `atlas.sum`

4. Review the generated migration

Check the generated SQL file in `infrastructure/pgsql/migrations/` to ensure it matches your expectations.

5. Apply/rollback the migration

```bash
make migrate-up

# To rollback the last applied migration:
make migrate-down
```

## Available Commands

See `Makefile`

### Migration conflicts

If multiple developers create migrations simultaneously, resolve conflicts by:
1. Keeping both migration files with different timestamps
2. Regenerating the `atlas.sum` file: `atlas migrate hash --env gorm`
