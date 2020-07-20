package discgov

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type state int

type guildInfo struct {
	usersLoc sync.Map
}

var guildInfos sync.Map
var stLock = sync.RWMutex{}

func init() {
}

//GetUsers Returns a list of users in a channel
func GetUsers(guildID, channelID string) []string {
	gInfoTmp, ok := guildInfos.Load(guildID)
	if !ok {
		return nil
	}
	gInfo := gInfoTmp.(*guildInfo)
	result := make([]string, 0)
	gInfo.usersLoc.Range(func(user, channel interface{}) bool {
		if channel.(string) == channelID {
			result = append(result, user.(string))
		}

		return true
	})
	return result
}

//UserVoiceTrackerHandler this is the handler which must be added as a handler to discrodgo it will update the
// Staes whenever a user moves between channels
func UserVoiceTrackerHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	gInfoTmp, _ := guildInfos.LoadOrStore(v.GuildID, &guildInfo{})
	gInfo := gInfoTmp.(*guildInfo)
	//Left Channel
	if v.ChannelID == "" {
		_, ok := gInfo.usersLoc.Load(v.UserID)
		if !ok {
			return
		}
		gInfo.usersLoc.Delete(v.UserID)
	} else {
		lastUserCh, ok := gInfo.usersLoc.Load(v.UserID)
		//First time user has been seen or user is moving between channels
		if !ok || lastUserCh.(string) != v.ChannelID {
			gInfo.usersLoc.Store(v.UserID, v.ChannelID)
		}
	}
}
