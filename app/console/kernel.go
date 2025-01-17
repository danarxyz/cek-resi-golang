package console

import (
	"goravel/app/http/controllers"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"
)

type Kernel struct {
}

func (kernel Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule().Call(func() {
			controllers.CheckExpedition()
		}).EveryFifteenMinutes(),
	}
}

func (kernel Kernel) Commands() []console.Command {
	return []console.Command{}
}
