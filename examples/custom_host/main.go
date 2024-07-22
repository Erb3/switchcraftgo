package main

import (
	"log"
	"net/url"
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
		Base: url.URL{
			Scheme: "ws",
			Host:   "localhost:8080",
			Path:   "/v2/",
		},
	})

	root := switchcraftgo.NewBrigadier(cb, "Custom host eh?")
	root.Register(root.Literal("proxied").Executes(func(bi *switchcraftgo.BrigadierInvocation) {
		bi.ReplyMarkdown("**Hi**")
	}))

	cb.Connect()
	cb.Listen()
}
