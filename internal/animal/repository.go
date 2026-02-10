package animal

import (
	"context"

	model "github.com/Carolineleless/petcare-jobs/internal/animal"
	"gorm.io/gorm"
)

var _ RepositoryInterface = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, animal *model.Animal) error {
	return r.db.WithContext(ctx).Create(animal).Error
}
