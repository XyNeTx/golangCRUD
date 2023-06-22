package responses

import (
	"uplevel-api/models"
)

type MyOrgResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []models.MyOrg `json:"details"`
}
type AllOrgsResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []models.AllOrgs `json:"details"`
}
type CreateOrgResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"details"`
}
type OrgSummaryResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    []models.OrgSummary `json:"details"`
}
