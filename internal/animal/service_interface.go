package animal

import "context"

type ServiceInterface interface {
	CreateAnimal(ctx context.Context, input animal) (CreateAnimalOutput, error)
}
