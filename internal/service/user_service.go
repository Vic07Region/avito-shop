package service

import (
	"context"
	"errors" //nolint:gci
	"github.com/Vic07Region/avito-shop/internal/storage"

	"github.com/Vic07Region/avito-shop/internal/utils" //nolint:gci
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt" //nolint:gci
)

func (s *Service) LoginUser(ctx context.Context, userdata UserData) (string, error) {
	userAuthData, err := s.Storage.GetUserAuthData(ctx, userdata.Username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userdata.Password), bcrypt.DefaultCost)
			if err != nil {
				//s.log.Error("LoginUser Generate passowrd hash Error:", zap.Error(err))
				return "", ErrGeneratePasswordHash
			}

			userID, err := s.Storage.NewUser(ctx, userdata.Username, string(hashedPassword))
			if err != nil {
				s.log.Error("LoginUser NewUser Error:", zap.Error(err))
				return "", err
			}

			token, err := utils.GenerateJWT(utils.User{UserID: userID})
			if err != nil {
				s.log.Error("LoginUser NewUser GenerateJWT error:", zap.Error(err))
				return "", ErrGenerateJWT
			}
			return token, nil
		}

		s.log.Error("LoginUser GetUserAuthData error:", zap.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAuthData.PasswordHash), []byte(userdata.Password))
	if err != nil {
		s.log.Error("LoginUser validate error:", zap.Error(err))
		return "", ErrInvelidPassword
	}

	token, err := utils.GenerateJWT(utils.User{UserID: userAuthData.UserID})
	if err != nil {
		s.log.Error("LoginUser GenerateJWT error:", zap.Error(err))
		return "", ErrGenerateJWT
	}

	return token, nil
}
