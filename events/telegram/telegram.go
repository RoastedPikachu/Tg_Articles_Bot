package telegram

import "read-adviser-bot/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
	// storage
}

func New(client *telegram.Client, storage)
