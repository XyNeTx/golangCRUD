package responses

import "uplevel-api/models"

type UserResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    *models.User `json:"details"`
}

type AccessTokenResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    *models.AccessToken `json:"details"`
}
