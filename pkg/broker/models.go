package broker

import "time"

type FeedbackMessageTypeEnum int

const (
	SuccessMessage FeedbackMessageTypeEnum = iota
	ErrorMessage
)

type AddAccountMessage struct {
	SteamID    string    `json:"steam_id"`
	Username   string    `json:"username"`
	GuildID    string    `json:"guild_id"`
	ChannelID  string    `json:"channel_id"`
	ExecutedAt time.Time `json:"executed_at"`
}

type SendFeedbackMessage struct {
	Content    string                  `json:"content"`
	Type       FeedbackMessageTypeEnum `json:"type"`
	Username   string                  `json:"username"`
	GuildID    string                  `json:"guild_id"`
	ChannelID  string                  `json:"channel_id"`
	ExecutedAt time.Time               `json:"executed_at"`
}
