package response

import "github.com/VulpesFerrilata/auth/internal/usecase/dto"

func NewTokenResponse(tokenDTO *dto.TokenDTO) *TokenResponse {
	tokenResponse := new(TokenResponse)
	tokenResponse.AccessToken = tokenDTO.AccessToken
	tokenResponse.RefreshToken = tokenDTO.RefreshToken
	return tokenResponse
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
