```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/uladzimirSTR/terraform-provider-trino/internal/client"
	"github.com/uladzimirSTR/terraform-provider-trino/internal/sqlrender"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	r, err := sqlrender.NewRenderer[sqlrender.CreateSchemaData]()

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	sql, err := r.Render("create_schema.sql.tmpl", sqlrender.CreateSchemaData{
		IfNotExists: true,
		TableSchema: sqlrender.TableSchema{
			Catalog:  "datalake",
			Name:     "_debug_tf",
			Location: "s3a://prod-datalake-gypsy",
		},
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf(sql)

	client, err := client.NewClient(client.Config{
		Host:       "91.98.21.193",
		Port:       8443,
		User:       "",
		Password:   "",
		Catalog:    "datalake",
		HTTPScheme: "https",

		PathToPEM:   "./",
		FileNamePEM: "trino.pem",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	err = client.Exec(ctx, sql)
	if err != nil {
		log.Fatal(err)
	}
}

DROP
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/uladzimirSTR/terraform-provider-trino/internal/client"
	"github.com/uladzimirSTR/terraform-provider-trino/internal/sqlrender"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	r, err := sqlrender.NewRenderer[sqlrender.DropSchemaData]()

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	sql, err := r.Render("drop_schema.sql.tmpl", sqlrender.DropSchemaData{
		IfExists: true,
		TableSchema: sqlrender.TableSchema{
			Catalog:  "datalake",
			Name:     "_debug_tf",
			Location: "s3a://prod-datalake-gypsy",
		},
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf(sql)

	client, err := client.NewClient(client.Config{
		Host:       "91.98.21.193",
		Port:       8443,
		User:       "",
		Password:   "",
		Catalog:    "datalake",
		HTTPScheme: "https",

		PathToPEM:   "./",
		FileNamePEM: "trino.pem",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	err = client.Exec(ctx, sql)
	if err != nil {
		log.Fatal(err)
	}
}

```
