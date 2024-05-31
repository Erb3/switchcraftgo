package switchcraftgo

import (
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	ChatboxFormattingMarkdown = iota
	ChatboxFormattingFormat
)

type ChatboxIngameUser struct {
	Type        string              `json:"type"`
	Name        string              `json:"name"`
	Uuid        string              `json:"uuid"`
	DisplayName string              `json:"displayName"`
	Group       string              `json:"group"`
	Pronouns    string              `json:"pronouns"`
	World       string              `json:"world"`
	Afk         bool                `json:"afk"`
	Alt         bool                `json:"alt"`
	Bot         bool                `json:"bot"`
	Supporter   uint8               `json:"supporter"`
	LinkedUser  *ChatboxDiscordUser `json:"linkedUser"`
}

type ChatboxDiscordUser struct {
	Type          string                `json:"type"`
	Id            uint64                `json:"id"`
	Name          string                `json:"name"`
	DisplayName   string                `json:"displayName"`
	Discriminator uint16                `json:"discriminator"`
	Avatar        string                `json:"avatar"`
	Roles         []*ChatboxDiscordRole `json:"roles"`
	LinkedUser    *ChatboxIngameUser    `json:"linkedUser"`
}

type ChatboxDiscordRole struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	Colour uint32 `json:"colour"`
}

type ChatboxTellPacket struct {
	Type string `json:"type"`
	User string `json:"user"`
	Text string `json:"text"`
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type ChatboxCommandPacket struct {
	Event     string            `json:"event"`
	User      ChatboxIngameUser `json:"user"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	OwnerOnly bool              `json:"ownerOnly"`
}

type ChatboxGenericEventPacket struct {
	Event string `json:"event"`
}

type Chatbox struct {
	Conn      *websocket.Conn
	scUrl     url.URL
	OnRaw     func(int, []byte)
	OnCommand func(ChatboxCommandPacket)
}

type NewChatboxOptions struct {
	Token string
	Base  url.URL
}

func GetDefaultBase() url.URL {
	return url.URL{
		Scheme: "wss",
		Host:   "chat.sc3.io",
		Path:   "/v2/",
	}
}

func NewChatbox(opts NewChatboxOptions) *Chatbox {
	scUrl := opts.Base
	if scUrl.Host == "" {
		scUrl = GetDefaultBase()
	}

	if opts.Token == "" {
		opts.Token = "guest"
	}

	scUrl.Path += opts.Token

	sc := &Chatbox{
		scUrl:     scUrl,
		OnRaw:     func(_ int, _ []byte) {},
		OnCommand: func(_ ChatboxCommandPacket) {},
	}

	return sc
}

func (sc *Chatbox) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(sc.scUrl.String(), nil)
	if err != nil {
		return err
	}
	sc.Conn = conn

	return nil
}

func (sc *Chatbox) Listen() {
	for {
		messageType, message, err := sc.Conn.ReadMessage()
		if err != nil {
			break
		}

		sc.OnRaw(messageType, message)

		var parsed ChatboxGenericEventPacket
		json.Unmarshal(message, &parsed)

		switch parsed.Event {
		case "command":
			var command ChatboxCommandPacket
			json.Unmarshal(message, &command)

			sc.OnCommand(command)
		}
	}
}

func (sc *Chatbox) Tell(user, message, name string, mode int) {
	formattingMode := "markdown"
	if mode == ChatboxFormattingFormat {
		formattingMode = "format"
	}
	packet := &ChatboxTellPacket{
		Type: "tell",
		User: user,
		Text: message,
		Name: name,
		Mode: formattingMode,
	}

	sc.Conn.WriteJSON(packet)
}
