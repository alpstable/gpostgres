# GPostgres 

Gidari Postgres is the PostgreSQL storage implementation for Gidari.

## Usage

```go
package main

import (
	"context"
	"database/sql"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/config"
	"github.com/alpstable/gpostgres"
)

func main() {
	ctx := context.TODO()

	database, _ := sql.Open("postgres", "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable")
	mdbStorage, _ := gpostgres.New(ctx, database)

	err := gidari.Transport(ctx, &config.Config{
		Storage: []Storage{mdbStorage},
	})
}
```
