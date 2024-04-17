package context

import (
	"simple-api-go/internal/app/repository"
	vehicleRepo "simple-api-go/internal/app/repository/vehicle"
	"simple-api-go/internal/app/service"
	vehicleSvc "simple-api-go/internal/app/service/vehicle"
)

// WiringRepositories ...
func (a *appContext) WiringRepositories() {
	a.SetRepositories(
		&repository.Repositories{

			Vehicle: vehicleRepo.NewDB(
				a.GetDBMysqlWrite(),
				a.GetDBMysqlRead(),
				a.GetTimeZone(),
				a.GetTimeFormat(),
				a.GetTimePostgresFriendlyFormat(),
				a.GetTimeLocation(),
			),
		},
	)
}

// WiringServices ...
func (a *appContext) WiringServices() {
	a.SetServices(&service.Services{
		Vehicle: vehicleSvc.NewSvc(
			a.GetConfig(),
			a.GetRepositories(),
			a.GetWatermillPubSub(),
			a.GetLogger(),
			a.GetCtx(),
			a.GetEnvironment(),
		),
	})
}

// WiringMachineryTaskMap ...
func (a *appContext) WiringMachineryTaskMap() {
	a.SetMachineryTaskMap(
		map[string]interface{}{},
	)
}
