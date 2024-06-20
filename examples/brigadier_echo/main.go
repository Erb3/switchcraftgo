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

	root := switchcraftgo.NewBrigadier(cb, "Echo")
	root.Register(root.Literal("echo").String("content").Executes(func(ev *switchcraftgo.BrigadierInvocation) {
		ev.ReplyMarkdown(ev.ReadString("content"))
	}))

	cb.Connect()
	cb.Listen()
}
