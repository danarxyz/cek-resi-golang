package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250117093802CreateResiTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250117093802CreateResiTable) Signature() string {
	return "20250117093802_create_resi_table"
}

// Up Run the migrations.
func (r *M20250117093802CreateResiTable) Up() error {
	if !facades.Schema().HasTable("resi") {
		return facades.Schema().Create("resi", func(table schema.Blueprint) {
			table.ID()
			table.Timestamps()
			table.String("tracking_num")
			table.Unique("tracking_num")
			table.Enum("expedition", []any{"spx", "jnt-cargo", "jnt", "tokopedia", "sicepat", "jne"})
			table.Enum("status", []any{"tracking", "done"}).Default("tracking")
			table.Text("details").Nullable()
			table.String("email")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250117093802CreateResiTable) Down() error {
	return facades.Schema().DropIfExists("resi")
}
