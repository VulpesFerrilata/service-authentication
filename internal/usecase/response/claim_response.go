package response

import "github.com/VulpesFerrilata/auth/internal/domain/datamodel"

func NewClaimResponse(claim *datamodel.Claim) *ClaimResponse {
	claimResponse := new(ClaimResponse)
	claimResponse.UserID = int(claim.UserID)
	return claimResponse
}

type ClaimResponse struct {
	UserID int `json:"userId"`
}
