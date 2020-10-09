package viewmodel

import (
	"github.com/VulpesFerrilata/auth/internal/usecase/dto"
	"github.com/VulpesFerrilata/grpc/protoc/auth"
)

func NewClaimResponse(claimResponsePb *auth.ClaimResponse) *ClaimResponse {
	return &ClaimResponse{
		claimResponsePb: claimResponsePb,
	}
}

type ClaimResponse struct {
	claimResponsePb *auth.ClaimResponse
}

func (cr *ClaimResponse) FromClaimDTO(claimDTO *dto.ClaimDTO) {
	cr.claimResponsePb.UserID = int64(claimDTO.UserID)
}
