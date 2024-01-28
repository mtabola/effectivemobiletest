package transportmodels

import dbmodels "effectivemobiletest/internal/models/db-models"

type Request struct {
	FName string `json:"name"`
	SName string `json:"surname"`
	PName string `json:"patronymic,omitempty"`
}

type Response struct {
	Status  int                `json:"status"`
	Users   []dbmodels.UserOut `json:"users,omitempty"`
	Error   string             `json:"error,omitempty"`
	Message string             `json:"message,omitempty"`
}
