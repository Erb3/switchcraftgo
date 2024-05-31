package main

import (
	"log"
	"os"

	"github.com/Erb3/switchcraftgo"
)

func main() {
	token := os.Getenv("CHATBOX_TOKEN")
	if token == "" {
		log.Fatalf("No CHATBOX_TOKEN environment variable set. Exiting.")
	}

	cb := switchcraftgo.NewChatbox(switchcraftgo.NewChatboxOptions{
		Token: token,
	})

	cb.OnCommand = func(ccp switchcraftgo.ChatboxCommandPacket) {
		if ccp.Command != "hello" {
			return
		}

		cb.Tell(ccp.User.Uuid, "&aHello World!", "Hello World Bot", switchcraftgo.ChatboxFormattingFormat)
	}

	cb.Connect()
	cb.Listen()
}
