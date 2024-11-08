package bot

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/slipneff/auto-yt/internal/utils/parser"
)

func (b *Bot) HandleUpdates(updates <-chan telego.Update) {
	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				b.Chat[update.Message.Chat.ID] = SceneMain
				b.handleStart(update.Message.Chat.ID)
				continue
			case "/parse":
				b.parseAccs(update.Message.Chat.ID)
				continue
			case "/get_users":
				b.handleGetUsers(update.Message.Chat.ID)
				continue

			default:
				if b.Chat[update.Message.Chat.ID] == SceneMain {
					b.handleStart(update.Message.Chat.ID)
				}
			}
		}
	}
}
func (b *Bot) handleStart(chatID int64) {
	b.SendMessage(&telego.SendMessageParams{
		ChatID: telego.ChatID{
			ID: chatID,
		},
		Text: "Hello, world!",
	})
}
func (b *Bot) handleGetUsers(chatID int64) {
	accs, err := b.service.GetUsers(context.Background())
	if err != nil {
		b.SendMessage(&telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatID,
			},
			Text: err.Error(),
		})
	}
	msg := "Users:\n\n"
	for _, acc := range accs {
		msg += acc.String() + "\n"
	}
	b.SendMessage(&telego.SendMessageParams{
		ChatID: telego.ChatID{
			ID: chatID,
		},
		Text: msg,
	})
}
func (b *Bot) parseAccs(chatID int64) {
	accs, err := parser.ReadFile("environments/accs/data.txt")
	if err != nil {
		b.SendMessage(&telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatID,
			},
			Text: err.Error(),
		})
	}
	err = b.service.BatchCreateUsers(context.Background(), accs)
	if err != nil {
		b.SendMessage(&telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatID,
			},
			Text: err.Error(),
		})
	}
	b.SendMessage(&telego.SendMessageParams{
		ChatID: telego.ChatID{
			ID: chatID,
		},
		Text: "Good!",
	})
}
