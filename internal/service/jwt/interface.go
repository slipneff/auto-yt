package jwt

import "github.com/slipneff/auto-yt/internal/utils/config"

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}
