package main

import (
	"fmt"
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

	root := switchcraftgo.NewBrigadier(cb, "Multiplier")
	root.Register(root.Literal("multiply").Number("factor1").Number("factor2").Executes(func(ev *switchcraftgo.BrigadierInvocation) {
		ev.ReplyMarkdown(fmt.Sprintf("Result is `%d`", ev.ReadNumber("factor1")*ev.ReadNumber("factor2")))
	}))

	cb.Connect()
	cb.Listen()
}
