package vehicle

import (
	"context"
	"simple-api-go/internal/app/entity"
)

type Repository interface {
	GetVehicle(ctx context.Context, motorBrand string, motorType string, motorTransmission string) (returnData []entity.Vehicle, err error)
}
