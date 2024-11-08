package main

import (
	"fmt"

	"github.com/slipneff/auto-yt/internal/ffmpeg"
	"github.com/slipneff/auto-yt/internal/utils/config"
	"github.com/slipneff/auto-yt/internal/utils/flags"
	"github.com/slipneff/gogger"
	"github.com/slipneff/gogger/log"
)

func main() {
	flags := flags.MustParseFlags()
	cfg := config.MustLoadConfig(flags.EnvMode, flags.ConfigPath)
	gogger.ConfigureZeroLogger()

	// container := di.New(context.Background(), cfg)
	log.Info(fmt.Sprintf("Server starting at %s:%d", cfg.Host, cfg.Port))

	// err := container.GetYoutubeClient().SearchVideos("новый год")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// client := container.GetYoutubeClient()
	// err := client.UploadVideo(&youtube.Video{
	// 	Title:       "test video222",
	// 	Description: "random",
	// 	FileName:    "files/1.mp4",
	// 	Category:    "10",
	// 	Keywords:    "viperr",
	// 	Privacy:     "public",
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// ffmpeg.UniqueVideo("")

	ffmpeg.UniqueVideo("files/2.mp4")

	// bot := container.NewBot()

	// updates, err := bot.UpdatesViaLongPolling(nil)
	// if err != nil {
	// 	panic(err)
	// }
	// defer bot.StopLongPolling()

	// bot.HandleUpdates(updates)
}
