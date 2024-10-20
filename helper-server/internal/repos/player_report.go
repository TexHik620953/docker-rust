package repos

import "helper-server/internal/models"

func (h *Repository) CreateReport(rp *models.PlayerReport) error {
	return h.db.Create(rp).Error
}
