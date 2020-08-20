package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	DB_USER     = "dscbot"
	DB_PASSWORD = "dscbot0215"
	DB_NAME     = "dscbot"
)

//DBconnect for connect to postgresql
func DBconnect(s *discordgo.Session, m *discordgo.MessageCreate, state int) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	//vaildate postgrewql db
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//create connection with Postgresql db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//create first input for channel_basic
	if state == 1 {
		message := m.Content
		message = strings.Replace(message, "!item ", "", 1)
		INFO := strings.Split(message, " ")
		//first check
		var first bool
		err := db.QueryRow("select activate from channel_basic where channelid=$1", m.ChannelID).Scan(&first)
		if err != nil {
			panic(err)
		}
		//not the first time
		if first == true {
			//sql for update
			updatesql := `
			UPDATE channel_basic
			SET channelinfo = $1, trellourl = $2
			WHERE channelinfo = $3
			;`
			//update info
			_, err = db.Exec(updatesql, INFO[0], INFO[1], m.ChannelID)
			if err != nil {
				panic(err)
			}
			s.ChannelMessageSend(m.ChannelID, "채널정보갱신")
			return
		}
		sqlStatement := `
		INSERT INTO channel_basic (channelid,channelinfo,trellourl,activate)
		VALUES ($1, $2, $3, true)`
		channelid := m.ChannelID
		_, err = db.Exec(sqlStatement, channelid, INFO[0], INFO[1])
		if err != nil {
			panic(err)
		}
		s.ChannelMessageSend(m.ChannelID, "추가완료")
	}
	//view channelinfo
	if state == 2 {
		//info is channelinfo
		var info string
		err := db.QueryRow("select channelinfo from channel_basic where channelid=$1", m.ChannelID).Scan(&info)
		if err != nil {
			panic(err)
		}
		//url is trellourl
		var url string
		err = db.QueryRow("select trellourl from channel_basic where channelid=$1", m.ChannelID).Scan(&url)
		if err != nil {
			panic(err)
		}
		//show info + url in discord
		s.ChannelMessageSend(m.ChannelID, "chanenlinfo: "+info+"\n trellourl: "+url)
	}
}

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message has "item" add info to DB
	if strings.Contains(m.Content, "!item") {
		DBconnect(s, m, 1)
	}
	if strings.Contains(m.Content, "!채널정보") {
		DBconnect(s, m, 2)
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
