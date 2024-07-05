package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
)

func launchBot() error {
	botToken, hasBotToken := os.LookupEnv("BOT_TOKEN")
	if !hasBotToken {
		return fmt.Errorf("Specify Telegram bot token with 'BOT_TOKEN' environment variable")
	}

	_, isDebug := os.LookupEnv("DEBUG")

	logger := telego.WithDefaultLogger(isDebug, true)
	bot, err := telego.NewBot(botToken, logger)
	if err != nil {
		return fmt.Errorf("cannot make new bot: %w", err)
	}

	// Call method getMe (https://core.telegram.org/bots/api#getme)
	botUser, err := bot.GetMe()
	if err != nil {
		return fmt.Errorf("getme error: %w", err)
	}

	// Print Bot information
	fmt.Printf("Bot user: %+v\n", botUser)
	return nil
}

func main() {
	err := launchBot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
