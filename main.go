package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

const welcomeChannelID = "803527411949895700"
const generalChannelID = "803527660365283349"
const armyRoleID = "797968006310920214"
const traitorID = "803776638356815912"

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	// dg.AddHandler(memberAdd)
	dg.AddHandler(memberUpdate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMembers)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func memberUpdate(s *discordgo.Session, u *discordgo.GuildMemberUpdate) {
	if u.Member.User.Bot {
		return
	}
	if u.Member.Nick == "" {
		s.ChannelMessageSend(welcomeChannelID, "NOT A JEPH?!?!?!")
		name := fmt.Sprintf("jeph - %v", StringWithCharset(5))
		s.GuildMemberNickname(u.GuildID, u.User.ID, name)
		s.ChannelMessageSend(welcomeChannelID, fmt.Sprintf("I dub you %v", name))
		s.GuildMemberRoleAdd(u.GuildID, u.User.ID, armyRoleID)
		return
	}
	if !strings.Contains(strings.ToLower(u.Member.Nick), "jeph") {
		name := fmt.Sprintf("traitor jeph - %v", StringWithCharset(5))
		s.ChannelMessageSend(generalChannelID, fmt.Sprintf("Nuh uh uhhh! <@%v> We are many. We are one.  You shall now be: %v", u.User.ID, name))
		s.GuildMemberNickname(u.GuildID, u.User.ID, name)
		s.GuildMemberRoleRemove(u.GuildID, u.User.ID, armyRoleID)
		s.GuildMemberRoleAdd(u.GuildID, u.User.ID, traitorID)
	}
}
