package repository

import (
	"context"
	"errors"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/library/pkg/db"
	server_errors "github.com/VulpesFerrilata/library/pkg/errors"
	"gorm.io/gorm"
)

type ReadOnlyTokenRepository interface {
	GetByUserId(ctx context.Context, userId uint) (*model.Token, error)
	GetByJti(ctx context.Context, jti string) (*model.Token, error)
}

type TokenRepository interface {
	ReadOnlyTokenRepository
	Insert(ctx context.Context, token *model.Token) error
	SaveByUserId(ctx context.Context, token *model.Token) error
}

func NewTokenRepository(dbContext *db.DbContext) TokenRepository {
	return &tokenRepository{
		dbContext: dbContext,
	}
}

type tokenRepository struct {
	dbContext *db.DbContext
}

func (tr tokenRepository) GetByUserId(ctx context.Context, userId uint) (*model.Token, error) {
	token := new(model.Token)
	err := tr.dbContext.GetDB(ctx).First(token, "user_id = ?", userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = server_errors.NewNotFoundError("token")
	}
	return token, err
}

func (tr tokenRepository) GetByJti(ctx context.Context, jti string) (*model.Token, error) {
	token := new(model.Token)
	err := tr.dbContext.GetDB(ctx).First(token, "jti = ?", jti).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = server_errors.NewNotFoundError("token")
	}
	return token, err
}

func (tr tokenRepository) Insert(ctx context.Context, token *model.Token) error {
	return tr.dbContext.GetDB(ctx).Model(token).Create(token).Error
}

func (tr tokenRepository) SaveByUserId(ctx context.Context, token *model.Token) error {
	if tr.dbContext.GetDB(ctx).Model(token).Where("user_id = ?", token.UserID).Updates(token).RowsAffected == 0 {
		return tr.Insert(ctx, token)
	}
	return nil
}
