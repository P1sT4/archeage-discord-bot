package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/bwmarrin/discordgo"
	archeage "github.com/geeksbaek/archeage-kr-web-go"
)

// Variables used for command line parameters
var (
	Token string = "MzAyNTc0NTA0MzQ4MDkwMzY4.C-F8cQ.5ZUMS_5TRGupESp7vNkYViOAN8M"
	BotID string
)

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandlerOnce(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	aa := archeage.New(&http.Client{})
	var oldNotices archeage.Notices
	ticker := time.Tick(time.Second * 10)
	for _ = range ticker {
		fmt.Println(".")
		newNotices, err := aa.FetchNotice()
		if err != nil || len(newNotices) == 0 {
			log.Println(err)
			continue
		}
		diffNotices := oldNotices.Diff(newNotices)
		if len(diffNotices) > 0 && len(oldNotices) > 0 {
			for _, notice := range diffNotices {
				msg := fmt.Sprintf("[%s] %s %s", notice.Category, notice.Title, notice.URL)
				s.ChannelMessageSend(m.ChannelID, msg)
				log.Println(msg)
			}
		}
		oldNotices = newNotices
	}
}
