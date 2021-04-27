package response

import "github.com/VulpesFerrilata/auth/internal/domain/model"

func NewClaimResponse(claim *model.Claim) *ClaimResponse {
	claimResponse := new(ClaimResponse)
	claimResponse.ID = claim.GetId().String()
	return claimResponse
}

type ClaimResponse struct {
	ID string `json:"id"`
}
