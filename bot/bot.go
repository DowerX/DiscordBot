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

	"os/exec"

	"../config"
	"../errorcheck"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var bot *discordgo.Session
var botID string
var botUser discordgo.User

var voiceConnection *discordgo.VoiceConnection
var musicStop chan bool
var clear string = `
‎`
var conf config.Config

// Start _
func Start(c config.Config) {
	var err error
	conf = c
	musicStop = make(chan bool)
	bot, err = discordgo.New("Bot " + conf.Token)
	errorcheck.Check(err)
	botUser, _ := bot.User("@me")
	botID = botUser.ID

	bot.UpdateStatus(0, conf.Status)

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
	// Look for bots
	if m.Author.Bot {
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
	case conf.Prefix + "ping":
		cmdPing(s, m)
		return

	case conf.Prefix + "join":
		cmdJoin(s, m)
		return
	case conf.Prefix + "dc":
		cmdDisconnect()
		return
	case conf.Prefix + "clear":
		cmdClear(s, m)
		return
	case conf.Prefix + "stop":
		cmdStop()
	}

	//Commands with arguments
	if strings.HasPrefix(m.Content, conf.Prefix+"play") {
		cmdPlay(s, m)
		return
	}

	if strings.HasPrefix(m.Content, conf.Prefix+"random") {
		cmdRandom(s, m)
		return
	}
}

func cmdPing(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "pong")
	errorcheck.Check(err)
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

func cmdPlay(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.Split(m.Content, " ")
	err := os.Remove("./temp.mp3")
	errorcheck.Check(err)
	cmd := exec.Command("youtube-dl", "-o", "./temp.mp3", parts[1], "-x", "--audio-format", "mp3")
	_ = cmd.Run()
	errorcheck.Check(err)
	cmdStop()
	cmdJoin(s, m)
	time.Sleep(time.Second)
	dgvoice.PlayAudioFile(voiceConnection, "./temp.mp3", musicStop)
}

func cmdStop() {
	cmd := exec.Command("pkill", "ffmpeg")
	err := cmd.Run()
	errorcheck.Check(err)
}

func cmdJoin(s *discordgo.Session, m *discordgo.MessageCreate) {
	vs, err := findUserVoiceState(bot, m.Author.ID)
	errorcheck.Check(err)
	voiceConnection, _ = bot.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, false)
}

func cmdDisconnect() {
	voiceConnection.Disconnect()
}

func cmdClear(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, strings.Repeat(clear, 200))
	errorcheck.Check(err)
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
