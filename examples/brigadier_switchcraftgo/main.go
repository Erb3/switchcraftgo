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

	root := switchcraftgo.NewBrigadier(cb, "&dSwitchCraftGo")
	root.Register(root.Literal("switchcraftgo").Then(root.Literal("api").Executes(func(bi *switchcraftgo.BrigadierInvocation) {
		bi.ReplyMarkdown("The API part of SwitchCraftGo is a wrapper for the [SwitchCraft API](https://docs.sc3.io/faq/api.html).")
	})).Then(root.Literal("chatbox").Executes(func(bi *switchcraftgo.BrigadierInvocation) {
		bi.ReplyMarkdown("The chatbox part of SwitchCraftGo is a very simple Chatbox wrapper.")
	}).Then(root.Literal("brigadier").Executes(func(bi *switchcraftgo.BrigadierInvocation) {
		bi.ReplyMarkdown("Brigadier is a part of the chatbox part.")
	}))))

	cb.Connect()
	cb.Listen()
}
