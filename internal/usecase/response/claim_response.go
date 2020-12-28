package response

import "github.com/VulpesFerrilata/auth/internal/domain/model"

func NewClaimResponse(claim *model.Claim) *ClaimResponse {
	claimResponse := new(ClaimResponse)
	claimResponse.UserID = int(claim.GetUserId())
	return claimResponse
}

type ClaimResponse struct {
	UserID int `json:"userId"`
}
