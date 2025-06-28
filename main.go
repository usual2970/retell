package main

import (
	"log"
	"os"
	"strings"

	"github.com/usual2970/retell/internal/routes"
	"github.com/usual2970/retell/internal/util/app"

	"github.com/pocketbase/pocketbase/core"

	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/usual2970/retell/migrations"

	_ "github.com/pocketbase/pocketbase/migrations"
)

func main() {
	app := app.Get()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnRecordCreateRequest("essay").BindFunc(func(e *core.RecordRequestEvent) error {
		return routes.OnessayCreate(e)
	})

	app.OnRecordUpdateRequest("essay").BindFunc(func(e *core.RecordRequestEvent) error {
		return routes.OnessayUpdate(e)
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		routes.Route(e.Router)
		if err := routes.Register(); err != nil {
			return err
		}
		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
