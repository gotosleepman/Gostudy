package response

import "lending-copy/api/models"

type Search struct {
	Rows  []models.Pool `json:"rows"`
	Count int64         `json:"count"`
}
