package repos

import "helper-server/internal/models"

func (h *Repository) CreateFeedback(fb *models.PlayerFeedback) error {
	return h.db.Create(fb).Error
}
