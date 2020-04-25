package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"../errorcheck"
	"github.com/bwmarrin/discordgo"
	//youtube "github.com/knadh/go-get-youtube"
)

var bot *discordgo.Session
var botID string
var botUser discordgo.User

// BotPerfix _
var BotPerfix string = "/"
var voiceConnection *discordgo.VoiceConnection

// MusicPath _
var MusicPath string = "/mnt/d/Music/"

var clear string = `
‎`

// Start _
func Start(token string) {
	var err error
	bot, err = discordgo.New("Bot " + token)
	errorcheck.Check(err)
	botUser, _ := bot.User("@me")
	botID = botUser.ID
	fmt.Println(botID)
	bot.AddHandler(messageHandler)
	err = bot.Open()
	errorcheck.Check(err)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID {
		return
	}

	// Look for mentions
	for i := 0; i < len(m.Mentions); i++ {
		u := *m.Mentions[i]
		if botID == u.ID {
			_, err := s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Viszont Kívánom!")
			errorcheck.Check(err)
			return
		}
	}

	// Single word commands
	switch m.Content {
	case BotPerfix + "ping":
		_, err := s.ChannelMessageSend(m.ChannelID, "pong")
		errorcheck.Check(err)
		return

	case BotPerfix + "join":
		vs, err := findUserVoiceState(bot, m.Author.ID)
		errorcheck.Check(err)
		voiceConnection, _ = bot.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, false)
		return
	case BotPerfix + "disconnect":
		voiceConnection.Disconnect()
		return
	case BotPerfix + "clear":
		_, err := s.ChannelMessageSend(m.ChannelID, strings.Repeat(clear, 150))
		errorcheck.Check(err)
		return
	}

	// Commands with arguments
	// 1st: play
	// if strings.HasPrefix(m.Content, BotPerfix+"play") {
	// 	youtube, err := youtube.Get(parts[1])
	// 	options := &youtube.Options{
	// 		Rename: true,
	// 		Resume: true,
	// 		Mp3:    true,
	// 	}
	// 	video.Download(0, MusicPath+"music.mp3", options)
	// 	dgvoice.PlayAudioFile(voiceConnection, MusicPath+"music.mp3", make(chan bool))
	// }

	if strings.HasPrefix(m.Content, BotPerfix+"random") {
		cmdRandom(s, m)
		return
	}
}

func cmdRandom(s *discordgo.Session, m *discordgo.MessageCreate) {
	rand.Seed(time.Now().UnixNano())
	parts := strings.Split(m.Content, " ")
	min, err := strconv.Atoi(parts[1])
	errorcheck.Check(err)
	max, err := strconv.Atoi(parts[2])
	errorcheck.Check(err)
	i := rand.Intn(max-min+1) + min
	_, err = s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" "+strconv.FormatInt(int64(i), 10))
}

func findUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}
