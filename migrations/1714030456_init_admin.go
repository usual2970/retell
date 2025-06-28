// migrations/1687801090_initial_superuser.go
package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		record := core.NewRecord(superusers)

		// note: the values can be eventually loaded via os.Getenv(key)
		// or from a special local config file
		record.Set("email", "admin@ikittg.com")
		record.Set("password", "1234567890")

		return app.Save(record)
	}, func(app core.App) error { // optional revert operation
		record, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "admin@ikittg.com")
		if record == nil {
			return nil // probably already deleted
		}

		return app.Delete(record)
	})
}
