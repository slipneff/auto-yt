package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/slipneff/auto-yt/internal/utils/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	cfg       *config.Config
	ytService *youtube.Service
}

func New(cfg *config.Config) *Client {
	ytService, err := youtube.NewService(context.Background(), option.WithAPIKey(cfg.DeveloperKey))
	if err != nil {
		panic(err)
	}
	return &Client{
		cfg:       cfg,
		ytService: ytService,
	}
}

func (c *Client) SearchVideos(query string) error {
	call := c.ytService.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(25)
	response, err := call.Do()
	if err != nil {
		return err
	}
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	fmt.Println("Videos", videos)
	fmt.Println("Channels", channels)
	fmt.Println("Playlists", playlists)
	return nil
}

type Video struct {
	Title       string
	Description string
	FileName    string
	Category    string
	Keywords    string
	Privacy     string
}

func (c *Client) UploadVideo(video *Video) error {

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       video.Title,
			Description: video.Description,
			CategoryId:  video.Category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: video.Privacy},
	}
	call := c.ytService.Videos.Insert([]string{"snippet", "status"}, upload)
	file, err := os.Open(video.FileName)
	if err != nil {
		log.Fatalf("Error opening %v: %v", video.FileName, err)
		return err
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		return err
	}
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
	return nil
}

func (c *Client) UploadTokens() *youtube.Service {
	ctx := context.Background()
	b, err := ioutil.ReadFile("environments/tokens/client_secret1.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл client_secret.json: %v", err)
	}

	// Создайте конфигурацию OAuth 2.0.
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
	if err != nil {
		log.Fatalf("Не удалось создать конфигурацию OAuth 2.0: %v", err)
	}
	token, err := getTokenFromWeb(config)
	if err != nil {
		log.Fatalf("Не удалось получить токен OAuth 2.0: %v", err)
	}
	client := config.Client(ctx, token)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Не удалось создать клиент YouTube Data API: %v", err)
	}
	c.ytService = service
	return service
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по ссылке для авторизации: %v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("не удалось прочитать код авторизации: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("не удалось обменять код авторизации на токен: %v", err)
	}

	// Сохраните токен для будущих запусков.
	saveToken("token.json", token)
	return token, nil
}
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Сохранение токена в файл: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Не удалось создать файл для сохранения токена: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
