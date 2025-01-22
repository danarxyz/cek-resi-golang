package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250122103527AddPackageNameToResiTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250122103527AddPackageNameToResiTable) Signature() string {
	return "20250122103527_add_package_name_to_resi_table"
}

// Up Run the migrations.
func (r *M20250122103527AddPackageNameToResiTable) Up() error {
	return facades.Schema().Table("resi", func(table schema.Blueprint) {
		table.String("package_name", 255)
	})
}

// Down Reverse the migrations.
func (r *M20250122103527AddPackageNameToResiTable) Down() error {
	return nil
}
