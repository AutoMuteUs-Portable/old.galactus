package handler

import (
	"encoding/json"
	redis_utils "github.com/automuteus/galactus/internal/redis"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func GuildDeleteHandler(logger *zap.Logger, client *redis.Client) func(s *discordgo.Session, m *discordgo.GuildDelete) {
	return func(s *discordgo.Session, m *discordgo.GuildDelete) {
		byt, err := json.Marshal(m)
		if err != nil {
			logger.Error("error marshalling json for GuildDelete message",
				zap.Error(err))
		}
		err = redis_utils.PushDiscordMessage(client, redis_utils.GuildDelete, byt)
		if err != nil {
			logger.Error("error pushing to Redis for GuildDelete message",
				zap.Error(err))
		} else {
			logger.Info("received GuildDelete message",
				zap.String("ID", m.ID),
			)
		}
	}
}
