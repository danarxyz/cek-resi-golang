package bootstrap

import (
	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/cache"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/crypt"
	"github.com/goravel/framework/database"
	"github.com/goravel/framework/event"
	"github.com/goravel/framework/filesystem"
	"github.com/goravel/framework/foundation"
	"github.com/goravel/framework/grpc"
	"github.com/goravel/framework/hash"
	"github.com/goravel/framework/http"
	"github.com/goravel/framework/log"
	"github.com/goravel/framework/mail"
	"github.com/goravel/framework/queue"
	"github.com/goravel/framework/route"
	"github.com/goravel/framework/schedule"
	"github.com/goravel/framework/session"
	"github.com/goravel/framework/testing"
	"github.com/goravel/framework/translation"
	"github.com/goravel/framework/validation"
	"github.com/goravel/framework/view"

	"goravel/app/providers"
	"goravel/config"
	"goravel/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithConfig(config.Boot).
		WithProviders(func() []contractsfoundation.ServiceProvider {
			return []contractsfoundation.ServiceProvider{
				&http.ServiceProvider{},
				&log.ServiceProvider{},
				&cache.ServiceProvider{},
				&session.ServiceProvider{},
				&validation.ServiceProvider{},
				&view.ServiceProvider{},
				&route.ServiceProvider{},
				&database.ServiceProvider{},
				&auth.ServiceProvider{},
				&crypt.ServiceProvider{},
				&queue.ServiceProvider{},
				&event.ServiceProvider{},
				&grpc.ServiceProvider{},
				&hash.ServiceProvider{},
				&translation.ServiceProvider{},
				&mail.ServiceProvider{},
				&schedule.ServiceProvider{},
				&filesystem.ServiceProvider{},
				&testing.ServiceProvider{},
				&providers.AppServiceProvider{},
				&providers.AuthServiceProvider{},
				&providers.RouteServiceProvider{},
				&providers.GrpcServiceProvider{},
				&providers.ConsoleServiceProvider{},
				&providers.QueueServiceProvider{},
				&providers.EventServiceProvider{},
				&providers.ValidationServiceProvider{},
				&providers.DatabaseServiceProvider{},
			}
		}).
		WithRouting(func() {
			routes.Api()
			routes.Grpc(nil)
		}).
		Create()
}
