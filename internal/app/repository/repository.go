package repository

import "simple-api-go/internal/app/repository/vehicle"

// Repositories all repo object injected here
type Repositories struct {
	Vehicle vehicle.Repository
}
