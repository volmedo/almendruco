package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"

	"github.com/volmedo/almendruco.git/internal/notifier"
	"github.com/volmedo/almendruco.git/internal/raices"
	"github.com/volmedo/almendruco.git/internal/repo"
	"github.com/volmedo/almendruco.git/internal/repo/dynamodbrepo"
)

const appName = "almendruco"

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler() error {
	cfg := config{}
	if err := envconfig.Process(appName, &cfg); err != nil {
		return fmt.Errorf("configuration processing failed: %w", err)
	}

	r, err := dynamodbrepo.NewRepo()
	if err != nil {
		return fmt.Errorf("unable to initialize repository: %w", err)
	}

	rc, err := raices.NewClient(cfg.Raices.BaseURL)
	if err != nil {
		return fmt.Errorf("error creating Raíces client: %w", err)
	}

	n, err := notifier.NewTelegramNotifier(cfg.Telegram.BaseURL, cfg.Telegram.BotToken)
	if err != nil {
		return fmt.Errorf("error creating notifier: %w", err)
	}

	if err := notifyMessages(r, rc, n); err != nil {
		return fmt.Errorf("error notifying messages: %w", err)
	}

	log.Println("Success!")

	return nil
}

func notifyMessages(r repo.Repo, rc raices.Client, n notifier.Notifier) error {
	chats, err := r.GetChats()
	if err != nil {
		return fmt.Errorf("unable to fetch chats from repo: %s", err)
	}

	for _, c := range chats {
		chatID, err := strconv.ParseUint(c.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("bad chatID %s: %w", c.ID, err)
		}

		msgs, err := rc.FetchMessages(c.Credentials, c.LastNotifiedMessage)
		if err != nil {
			return fmt.Errorf("error fetching messages from Raíces: %s", err)
		}

		if len(msgs) != 0 {
			last, err := n.Notify(notifier.ChatID(chatID), msgs)
			if err != nil {
				// Notify notifies messages until it encounters an error, so even in the case of an error
				// happening we can still update last notified message to avoid notifying again messages
				// that have already been notified
				_ = r.UpdateLastNotifiedMessage(strconv.FormatUint(chatID, 10), last)
				return fmt.Errorf("error notifying messages: %s", err)
			}

			if err := r.UpdateLastNotifiedMessage(strconv.FormatUint(chatID, 10), last); err != nil {
				return fmt.Errorf("error updating last notified message: %s", err)
			}
		}
	}

	return nil
}
