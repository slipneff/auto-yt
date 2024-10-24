package main

import (
	"context"
	"fmt"

	"github.com/slipneff/auto-yt/internal/di"
	"github.com/slipneff/auto-yt/internal/utils/config"
	"github.com/slipneff/auto-yt/internal/utils/flags"
	"github.com/slipneff/auto-yt/pkg/clients/youtube"
	"github.com/slipneff/gogger"
	"github.com/slipneff/gogger/log"
)

func main() {
	flags := flags.MustParseFlags()
	cfg := config.MustLoadConfig(flags.EnvMode, flags.ConfigPath)
	gogger.ConfigureZeroLogger()

	container := di.New(context.Background(), cfg)
	log.Info(fmt.Sprintf("Server starting at %s:%d", cfg.Host, cfg.Port))

	// err := container.GetYoutubeClient().SearchVideos("новый год")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	client := container.GetYoutubeClient()
	client.UploadTokens()
	err := client.UploadVideo(&youtube.Video{
		Title:       "test video",
		Description: "random",
		FileName:    "files/1.mp4",
		Category:    "42",
		Keywords:    "viperr",
		Privacy:     "unlisted",
	})
	if err != nil {
		fmt.Println(err)
	}
}
