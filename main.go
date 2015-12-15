package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	irc "github.com/fluffle/goirc/client"
	picarto "github.com/liru/picarto/stream"
)

var (
	nick   = flag.String("nick", "StreamAnnouncer", "The bot's nickname on IRC.")
	server = flag.String("server", "irc.rizon.net", "The IRC server to connect to.")

	artists              = picarto.ArtistMap{}
	channelsToAnnounceOn = make(map[string][]string) // map artist to slice of channels
	lastOnline           = make(map[string]time.Time)
	bot                  *irc.Conn
	config               tomlConfig
)

const (
	announcementMessage = "[ \x033PICARTO\x03 ] %s is streaming! https://picarto.tv/%s"
)

type tomlConfig struct {
	Channels map[string]Channel
}

type Channel struct {
	Artists []string
}

func announce(artist string) {
	for _, channel := range channelsToAnnounceOn[artist] {
		bot.Privmsg("#"+channel, fmt.Sprintf(announcementMessage, artist, artist))
	}
}

func makeIrcBot() {

	cfg := irc.NewConfig(*nick)
	cfg.Server = *server
	cfg.NewNick = func(n string) string { return n + "_" }
	c := irc.Client(cfg)
	c.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			for channelName, _ := range config.Channels {
				conn.Join("#" + channelName)
			}
		})
	// c.HandleFunc("privmsg", changeMonitor)
	// c.HandleFunc("privmsg", listArtists)

	if err := c.Connect(); err != nil {
		// log.Fatal("Connection error: ", err)
	}

	bot = c
}

func main() {
	flag.Parse()

	if _, err := toml.DecodeFile("artists.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	for channelName, channelInfo := range config.Channels {
		fmt.Printf("Chan: %s; Artists: %s\n", channelName, channelInfo.Artists)
		for _, x := range channelInfo.Artists {
			artists.AddArtist(x)
			channelsToAnnounceOn[x] = append(channelsToAnnounceOn[x], channelName)
		}
	}

	fmt.Println(channelsToAnnounceOn)
	done := make(chan struct{})

	announceChan := artists.MakeAnnounceChan(done)

	for {
		notification := <-announceChan
		artist, thisTime := notification.Name, notification.Time
		if time.Since(lastOnline[artist]) > (15 * time.Minute) {
			go announce(artist)
		}
		lastOnline[artist] = thisTime
	}

}
