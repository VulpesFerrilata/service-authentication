package service

type tokenServiceType int

const (
	AccessToken tokenServiceType = iota
	RefreshToken
)

type TokenServiceFactory interface {
	GetTokenService(tokenServiceType tokenServiceType) TokenService
}

func NewTokenServiceFactory(accessTokenService TokenService, refreshTokenService TokenService) TokenServiceFactory {
	return &tokenServiceFactory{
		accessTokenService:  accessTokenService,
		refreshTokenService: refreshTokenService,
	}
}

type tokenServiceFactory struct {
	accessTokenService  TokenService
	refreshTokenService TokenService
}

func (t tokenServiceFactory) GetTokenService(tokenServiceType tokenServiceType) TokenService {
	switch tokenServiceType {
	case AccessToken:
		return t.accessTokenService
	case RefreshToken:
		return t.refreshTokenService
	default:
		return nil
	}
}
