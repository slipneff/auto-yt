package ffmpeg

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func UniqueVideo(videoPath string) {
	filename := filepath.Base(videoPath)
	rand.Seed(time.Now().UnixNano())

	projectDir := "./"

	colorTemperature := 6500 + rand.Float64()*1000 - 500 // Диапазон: 6000-7000
	vibrance := rand.Float64()*0.1 - 0.05                // Диапазон: -0.1 до 0.1
	brightness := rand.Float64()*0.1 - 0.05              // Диапазон: -0.1 до 0.1
	gamma := 1.0 + rand.Float64()*0.1 - 0.05             // Диапазон: 0.9 до 1.1
	contrast := 1.0 + rand.Float64()*0.1 - 0.05          // Диапазон: 0.9 до 1.1
	speed := 1.0 + rand.Float64()*0.1 - 0.05             // Диапазон: 0.9 до 1.1
	blur := rand.Intn(2)                                 // Диапазон: до 3
	// flipping := rand.Intn(1)

	tempDir := filepath.Join(projectDir, "temp")
	readyVideosDir := filepath.Join(projectDir, "ready_videos")
	os.MkdirAll(tempDir, os.ModePerm)
	os.MkdirAll(readyVideosDir, os.ModePerm)
	defer removeTempFiles(tempDir, filename)
	// 1. Изменяем температуру
	tempColorTemp := filepath.Join(tempDir, fmt.Sprintf("color_temp_%s", filename))
	err := ffmpeg_go.Input(videoPath).
		Filter("colortemperature", ffmpeg_go.Args{fmt.Sprintf("%f", colorTemperature)}).
		Output(tempColorTemp, ffmpeg_go.KwArgs{
			"codec:v": "libx264",
			"preset":  "veryslow",
			"crf":     "10",
			"codec:a": "copy",
			"map":     "0:a",
			// "vf":      "scale=iw:ih", // Сохраняем исходное разрешение
		}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}

	// 2. Редактируем насыщенность
	tempVibrance := filepath.Join(tempDir, fmt.Sprintf("vibrance_%s", filename))
	err = ffmpeg_go.Input(tempColorTemp).
		Filter("vibrance", ffmpeg_go.Args{fmt.Sprintf("%f", vibrance)}).
		Output(tempVibrance, ffmpeg_go.KwArgs{"codec:v": "libx264", "preset": "veryslow", "crf": "10", "codec:a": "copy", "map": "0:a"}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}

	// 3. Изменяем яркость
	tempBrightness := filepath.Join(tempDir, fmt.Sprintf("brightness_%s", filename))
	err = ffmpeg_go.Input(tempVibrance).
		Filter("eq", ffmpeg_go.Args{fmt.Sprintf("brightness=%f", brightness)}).
		Output(tempBrightness, ffmpeg_go.KwArgs{"codec:v": "libx264", "preset": "veryslow", "crf": "10", "codec:a": "copy", "map": "0:a"}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}

	// 4. Изменяем гамму
	tempGamma := filepath.Join(tempDir, fmt.Sprintf("gamma_%s", filename))
	err = ffmpeg_go.Input(tempBrightness).
		Filter("eq", ffmpeg_go.Args{fmt.Sprintf("gamma=%f", gamma)}).
		Output(tempGamma, ffmpeg_go.KwArgs{"codec:v": "libx264", "preset": "veryslow", "crf": "10", "codec:a": "copy", "map": "0:a"}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}

	// 5. Изменяем контрастность
	tempContrast := filepath.Join(tempDir, fmt.Sprintf("contrast_%s", filename))
	err = ffmpeg_go.Input(tempGamma).
		Filter("eq", ffmpeg_go.Args{fmt.Sprintf("contrast=%f", contrast)}).
		Output(tempContrast, ffmpeg_go.KwArgs{"codec:v": "libx264", "preset": "veryslow", "crf": "10", "codec:a": "copy", "map": "0:a"}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}

	// 6. Изменяем скорость
	tempSpeed := filepath.Join(tempDir, fmt.Sprintf("speed_%s", filename))
	err = ffmpeg_go.Input(tempContrast).
		Filter("setpts", ffmpeg_go.Args{fmt.Sprintf("%f*PTS", speed)}).
		Output(tempSpeed, ffmpeg_go.KwArgs{"codec:v": "libx264", "preset": "veryslow", "crf": "10", "codec:a": "copy", "map": "0:a"}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}

	// 7. Применяем блюр
	finalOutput := filepath.Join(readyVideosDir, fmt.Sprintf("final_%s", filename))
	err = ffmpeg_go.Input(tempSpeed).
		Filter("boxblur", ffmpeg_go.Args{fmt.Sprintf("%d", blur)}).
		Output(finalOutput, ffmpeg_go.KwArgs{"codec:v": "libx264", "preset": "veryslow", "crf": "10", "codec:a": "copy", "map": "0:a"}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}
	// 8. Отзеркаливание
	// var tempFlip = tempBlur
	// if flipping != 0 {
	// 	tempFlip = filepath.Join(readyVideosDir, fmt.Sprintf("flip_%s", filename))
	// 	err = ffmpeg_go.Input(tempBlur).
	// 		Filter("hflip", ffmpeg_go.Args{}).
	// 		Output(tempFlip, ffmpeg_go.KwArgs{"codec:v": "libx264", "preset": "veryslow", "crf": "10", "codec:a": "copy","map":"0:a"}).
	// 		OverWriteOutput().
	// 		Run()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// // 10. Удаление метаданных
	// finalOutput := filepath.Join(readyVideosDir, fmt.Sprintf("final_%s", filename))
	// err = ffmpeg_go.Input(tempFlip).
	// 	Output(finalOutput, ffmpeg_go.KwArgs{"map_metadata": "-1"}).
	// 	OverWriteOutput().
	// 	Run()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Видео успешно обработано и сохранено в:", finalOutput)

}
func removeTempFiles(tempDir, filename string) {
	// Удаление всех файлов во временной директории, название которых заканчивается на filename
	files, err := os.ReadDir(tempDir)
	if err != nil {
		fmt.Println("Ошибка при чтении временной директории:", err)
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), filename) {
			filePath := filepath.Join(tempDir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println("Ошибка при удалении файла:", err)
			}
		}
	}
}
