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

	root := switchcraftgo.NewBrigadier(cb, "Owner Only Test")
	root.Register(root.Literal("owneronly").Executes(func(bi *switchcraftgo.BrigadierInvocation) {
		status := "not"
		if bi.OwnerOnly {
			status = "actually"
		}

		bi.ReplyMarkdown(fmt.Sprintf("You did **%s** call this command owner-only", status))
	}))

	cb.Connect()
	cb.Listen()
}
