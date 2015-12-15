package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	// picarto "github.com/liru/picarto/stream"
)

type tomlConfig struct {
	Channels map[string]Channel
}

type Channel struct {
	Artists []string
}

func main() {
	var config tomlConfig
	if _, err := toml.DecodeFile("artists.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	for channelName, channelInfo := range config.Channels {
		fmt.Printf("Chan: %s; Artists: %s\n", channelName, channelInfo)
	}
}
