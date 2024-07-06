package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type ConfigError struct {
	error
}

func NewConfigError(err error) *ConfigError {
	return &ConfigError{error: err}
}

func configureBot() (*telego.Bot, error) {
	botToken, hasBotToken := os.LookupEnv("BOT_TOKEN")
	if !hasBotToken {
		return nil, fmt.Errorf("Specify Telegram bot token with 'BOT_TOKEN' environment variable")
	}

	_, isDebug := os.LookupEnv("DEBUG")

	logger := telego.WithDefaultLogger(isDebug, true)
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
	// Creating keyboard
	keyboard := tu.Keyboard(
		tu.KeyboardRow( // Row 1
			// Column 1
			tu.KeyboardButton("Button"),

			// Column 2, `with` method
			tu.KeyboardButton("Poll Regular").WithRequestPoll(tu.PollTypeRegular()),
		),
		tu.KeyboardRow( // Row 2
			// Column 1, `with` method
			tu.KeyboardButton("Contact").WithRequestContact(),

			// Column 2, `with` method
			tu.KeyboardButton("Vote for").WithRequestUsers(&telego.KeyboardButtonRequestUsers{}),
		),
		tu.KeyboardRow( // Row 3
			tu.KeyboardButton("Griatech").WithWebApp(tu.WebAppInfo("https://gria.tech")),
			tu.KeyboardButton("Requestchat").WithRequestChat(&telego.KeyboardButtonRequestChat{
				RequestID: int32(update.UpdateID),
			}),
		),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Select something")
	// Multiple `with` methods can be chained

	// Creating message
	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Hello World",
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

func launchBot(ctx context.Context) error {
	bot, err := configureBot()
	if err != nil {
		return NewConfigError(fmt.Errorf("configure error: %w", err))
	}

	err = startProcessingUpdates(ctx, bot, 2)
	if err != nil {
		return fmt.Errorf("processing updates error: %w", err)
	}

	return nil
}

func main() {
	for {
		err := launchBot(context.Background())
		var cfgErr *ConfigError
		if errors.As(err, &cfgErr) {
			fmt.Println(err)
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("An error occured: %+v\n", err)
		}
		_, isRestarting := os.LookupEnv("RESTARTING")
		if !isRestarting {
			break
		}
		timeoutSeconds := 1
		fmt.Printf("Restarting bot processing after %d seconds", timeoutSeconds)
		time.Sleep(time.Duration(timeoutSeconds) * time.Second)
	}
	fmt.Printf("Finish bot processing")
}
