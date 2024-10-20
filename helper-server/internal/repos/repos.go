package repos

import (
	"helper-server/internal/config"
	"helper-server/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		return pgErr.Code == "23505"

	}
	return false
}

type Repository struct {
	db *gorm.DB
}

func New(cfg *config.AppConfig) (*Repository, error) {
	h := &Repository{}
	var err error

	h.db, err = gorm.Open(postgres.Open(cfg.DatabaseDSN))
	if err != nil {
		return nil, err
	}
	err = migrate(h.db)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func migrate(db *gorm.DB) error {
	models := []interface{}{
		&models.Player{},
		&models.PlayerReport{},
		&models.PlayerFeedback{},
	}
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}
