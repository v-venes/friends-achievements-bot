package models

import "time"

// TODO: Mover esse cara lá para o pkg
type AddAccountMessage struct {
	SteamID    string    `json:"steam_id"`
	Username   string    `json:"username"`
	GuildID    string    `json:"guild_id"`
	ChannelID  string    `json:"channel_id"`
	ExecutedAt time.Time `json:"executed_at"`
}
