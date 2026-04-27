```go
import (
	"fmt"
	"log"

	"github.com/uladzimirSTR/terraform-provider-trino/internal/sqlrender"
)

func main() {
	r, err := sqlrender.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	sql, err := r.Render("create_table.sql.tmpl", sqlrender.CreateTableData{
		IfNotExists: true,
		Table: sqlrender.Table{
			TableSchema: sqlrender.TableSchema{
				Catalog:  "datalake",
				Name:     "stage",
				Location: "s3://my-bucket/data",
			},
			TableName: "users",
			Columns: []sqlrender.Column{
				{ColName: "id", ColType: "BIGINT"},
				{ColName: "email", ColType: "VARCHAR"},
				{ColName: "created_at", ColType: "TIMESTAMP"},
			},
			TableProp: sqlrender.TableProperties{
				Format:        "PARQUET",
				PartitionedBy: []string{"created_date"},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sql)
}
```
