package main

import (
	"context"
	"errors"
	"fmt"
	"kry127.ru/weddingbot/config"
	"os"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func configureBot(config *config.Config) (*telego.Bot, error) {
	botToken, hasBotToken := os.LookupEnv("BOT_TOKEN")
	if !hasBotToken {
		return nil, fmt.Errorf("Specify Telegram bot token with 'BOT_TOKEN' environment variable")
	}

	logger := telego.WithDefaultLogger(config.Debug, true)
	bot, err := telego.NewBot(botToken, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot make new bot: %w", err)
	}

	// Call method getMe (https://core.telegram.org/bots/api#getme)
	botUser, err := bot.GetMe()
	if err != nil {
		return nil, fmt.Errorf("getme error: %w", err)
	}

	// Print Bot information
	fmt.Printf("Bot user: %+v\n", botUser)
	return bot, nil
}

func processUpdate(bot *telego.Bot, update telego.Update) error {
	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("✍️ Подписаться на обновления"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("📍 Где и когда свадьба?"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("🍽 Выбрать блюдо"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("🎵 Предложить песню DJ"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("🎉 Поздравить молодожёнов"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("📝 Расписание мероприятия"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("🤔 Что происходит сейчас"),
		),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Select something")

	// Creating message
	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Свадьба Марии и Виктора! Выберите пункт меню:",
	).WithReplyMarkup(keyboard).WithProtectContent() // Multiple `with` method

	_, err := bot.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("send message error: %w", err)
	}
	return nil
}

func startProcessingUpdates(ctx context.Context, bot *telego.Bot, workerCount int) error {
	updates, err := bot.UpdatesViaLongPolling(new(telego.GetUpdatesParams).WithTimeout(5))
	if err != nil {
		return fmt.Errorf("receive updates error: %w", err)
	}
	defer bot.StopLongPolling()

	var wg sync.WaitGroup
	wg.Add(workerCount)
	ctxx, cancel := context.WithCancelCause(ctx)
	for workerID := 0; workerID < workerCount; workerID++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case update := <-updates:
					err := processUpdate(bot, update)
					if err != nil {
						cancel(fmt.Errorf("worker %d stopped with error: %w", workerID, err))
					}
				case <-ctxx.Done():
					return
				}
			}
		}(workerID)
	}
	wg.Wait()
	if context.Cause(ctxx) != nil {
		return fmt.Errorf("worker stopped working with error: %w", context.Cause(ctxx))
	}
	return ctxx.Err()
}

func launchBot(ctx context.Context, cfg *config.Config) error {
	bot, err := configureBot(cfg)
	if err != nil {
		return config.NewConfigError(fmt.Errorf("configure error: %w", err))
	}

	err = startProcessingUpdates(ctx, bot, 1)
	if err != nil {
		return fmt.Errorf("processing updates error: %w", err)
	}

	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("cannot load config: %+v\n", err)
		os.Exit(1)
	}

	for {
		err = launchBot(context.Background(), cfg)
		var cfgErr *config.ConfigError
		if errors.As(err, &cfgErr) {
			fmt.Println(err)
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("An error occured: %+v\n", err)
		}

		timeoutSeconds := cfg.RestartTimeoutSeconds
		if timeoutSeconds < 0 {
			break
		}
		fmt.Printf("Restarting bot processing after %d seconds", timeoutSeconds)
		time.Sleep(time.Duration(timeoutSeconds) * time.Second)
	}
	fmt.Printf("Finish bot processing")
}
