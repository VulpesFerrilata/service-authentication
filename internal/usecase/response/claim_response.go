package response

import "github.com/VulpesFerrilata/auth/internal/domain/model"

func NewClaimResponse(claim *model.Claim) *ClaimResponse {
	claimResponse := new(ClaimResponse)
	claimResponse.UserID = claim.GetUserId().String()
	return claimResponse
}

type ClaimResponse struct {
	UserID string `json:"user_id"`
}
