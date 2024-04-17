package vehicle

import (
	"context"
	"simple-api-go/internal/app/config"
	"simple-api-go/internal/app/entity"
	"simple-api-go/internal/app/repository"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog"
)

type svc struct {
	Config             *config.Data
	Repo               *repository.Repositories
	WatermillPublisher message.Publisher
	Logger             zerolog.Logger
	AppCtx             context.Context
	Environment        string
}

var _ Service = (*svc)(nil)

func NewSvc(
	config *config.Data,
	repo *repository.Repositories,
	watermillPublisher message.Publisher,
	logger zerolog.Logger,
	appCtx context.Context,
	env string,
) (returnData Service) {
	return &svc{
		Config:             config,
		Repo:               repo,
		WatermillPublisher: watermillPublisher,
		Logger:             logger,
		AppCtx:             appCtx,
		Environment:        env,
	}
}

func (s *svc) GetVehicle(ctx context.Context, motorBrand string, motorType string, motorTransmission string) (returnData []entity.Vehicle, err error) {
	return s.Repo.Vehicle.GetVehicle(
		ctx,
		motorBrand,
		motorType,
		motorTransmission,
	)
}
