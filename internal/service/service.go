package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/slipneff/auto-yt/internal/models"
	"github.com/slipneff/auto-yt/internal/utils/parser"
	yt "google.golang.org/api/youtube/v3"
)

// func (s *Service) PublishVideo(c context.Context, users []uuid.UUID) {
// 	s.ytClient.UploadVideo(&youtube.Video{
// 		Title:       "",
// 		Description: "",
// 		FileName:    "",
// 		Category:    "",
// 		Keywords:    "",
// 		Privacy:     "",
// 	})
// 	return
// }

func (s *Service) CreateUser(c context.Context, user *models.Account) (*models.Account, error) {
	account, err := s.storage.CreateAccount(c, user)
	if err != nil {
		return nil, err
	}
	return account, nil
}
func (s *Service) GetUser(c context.Context, user *models.Account) (*models.Account, error) {
	account, err := s.storage.GetAccount(c, user.ID)
	if err != nil {
		return nil, err
	}
	return account, nil
}
func (s *Service) GetUsers(c context.Context) ([]*models.Account, error) {
	account, err := s.storage.GetAccounts(c)
	if err != nil {
		return nil, err
	}
	return account, nil
}
func (s *Service) BatchCreateUsers(c context.Context, user *parser.Accounts) error {
	accs := make([]*models.Account, 0, len(user.Accounts))
	for _, acc := range user.Accounts {
		accs = append(accs, &models.Account{
			Email:    acc.Email,
			Password: acc.Password,
			Recovery: acc.Recovery,
		})
	}
	err := s.storage.BatchCreateAccounts(c, accs)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetAuthUser(c context.Context, id uuid.UUID) error {
	client := s.ytClient.GetClient(yt.YoutubeUploadScope, id.String())
	if client == nil {
		return nil
	}
	err := s.storage.ConfirmAuth(c, id)
	if err != nil {
		return err
	}
	return nil
}
