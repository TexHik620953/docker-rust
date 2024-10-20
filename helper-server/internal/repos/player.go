package repos

import "helper-server/internal/models"

func (h *Repository) AttachPlayer(ply *models.Player) error {
	r := h.db.Where(&models.Player{SteamID: ply.SteamID}).
		Assign(&models.Player{
			PlayerName: ply.PlayerName,
		}).
		FirstOrCreate(ply)
	return r.Error
}

func (h *Repository) UpdateBalance(ply *models.Player) error {
	return h.db.Model(ply).Update("balance", ply.Balance).Error
}
