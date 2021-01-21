package handler

import (
	"encoding/json"
	redis_utils "github.com/automuteus/galactus/internal/redis"
	"github.com/automuteus/galactus/pkg/discord_message"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func VoiceStateUpdateHandler(logger *zap.Logger, client *redis.Client) func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	return func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		if m == nil {
			return
		}
		// ignore the bot
		if m.UserID == s.State.User.ID {
			return
		}

		// if no active games, completely ignore message reactions
		if !redis_utils.AnyActiveGamesInGuild(client, m.GuildID) {
			return
		}

		// a game is happening in this guild; in the background, make sure it's pruned of inactive games
		go redis_utils.PurgeOldGuildGames(client, m.GuildID)

		byt, err := json.Marshal(m)
		if err != nil {
			logger.Error("error marshalling json for VoiceStateUpdate message",
				zap.Error(err))
		}
		err = redis_utils.PushDiscordMessage(client, discord_message.VoiceStateUpdate, byt)
		if err != nil {
			logger.Error("error pushing to Redis for VoiceStateUpdate message",
				zap.Error(err))
		} else {
			logger.Info("pushed discord message to Redis",
				zap.String("type", discord_message.DiscordMessageTypeStrings[discord_message.VoiceStateUpdate]),
				zap.String("guild_id", m.GuildID),
				zap.String("channel_id", m.ChannelID),
				zap.String("user_id", m.UserID),
				zap.String("id", m.SessionID),
			)
		}
	}
}
