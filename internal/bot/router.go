package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler func(b *Bot, m *tgbotapi.Message)

type Router struct {
	commands map[string]Handler
}

func NewRouter() *Router {
	return &Router{commands: make(map[string]Handler)}
}

func (r *Router) RegisterCommand(key string, h Handler) {
	r.commands[key] = h
}

func (r *Router) Resolve(text string) (Handler, bool) {
	h, ok := r.commands[text]
	return h, ok
}
