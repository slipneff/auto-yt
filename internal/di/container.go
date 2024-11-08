package di

import (
	"context"
	"net/http"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/mymmrac/telego"

	"github.com/slipneff/auto-yt/internal/bot"
	"github.com/slipneff/auto-yt/internal/service"
	"github.com/slipneff/auto-yt/internal/storage/sql"
	"github.com/slipneff/auto-yt/pkg/clients/youtube"
	"gorm.io/gorm"

	"github.com/slipneff/auto-yt/internal/utils/config"
)

type Container struct {
	cfg *config.Config
	ctx context.Context

	telebot            *telego.Bot
	bot                *bot.Bot
	httpServer         *http.Server
	service            *service.Service
	ytClient           *youtube.Client
	storage            *sql.Storage
	db                 *gorm.DB
	transactionManager trm.Manager
}

func New(ctx context.Context, cfg *config.Config) *Container {
	return &Container{cfg: cfg, ctx: ctx}
}

func (c *Container) Bot() *telego.Bot {
	if c.telebot == nil {
		bot, err := telego.NewBot(c.cfg.BotToken, telego.WithDefaultDebugLogger())
		if err != nil {
			panic(err)
		}
		c.telebot = bot
	}

	return c.telebot
}

func (c *Container) NewBot() *bot.Bot {
	if c.bot == nil {
		c.bot = bot.New(c.Bot(), c.GetService())
	}

	return c.bot
}
func (c *Container) GetYoutubeClient() *youtube.Client {
	return get(&c.ytClient, func() *youtube.Client {
		return youtube.New(c.cfg)
	})
}

func (c *Container) GetPostgresDB() *sql.Storage {
	return get(&c.storage, func() *sql.Storage {
		return sql.New(c.GetDB(), trmgorm.DefaultCtxGetter)
	})
}

func (c *Container) GetDB() *gorm.DB {
	return get(&c.db, func() *gorm.DB {
		return sql.MustNewSQLite(c.cfg)
	})
}

func (c *Container) GetTransactionManager() trm.Manager {
	return get(&c.transactionManager, func() trm.Manager {
		return manager.Must(trmgorm.NewDefaultFactory(c.db))
	})
}

func (c *Container) GetService() *service.Service {
	return get(&c.service, func() *service.Service {
		return service.NewService(c.GetPostgresDB(), c.GetYoutubeClient())
	})
}

func get[T comparable](obj *T, builder func() T) T {
	if *obj != *new(T) {
		return *obj
	}

	*obj = builder()
	return *obj
}
