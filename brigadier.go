package switchcraftgo

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// An instance of Brigadier.
// Must be created with the [NewBrigadier] function.
type Brigadier struct {
	conn *Chatbox
	cmds []*BrigadierCommand
	name string
}

// A registered Brigadier Command.
// Must be made with [Brigadier.Literal].
type BrigadierCommand struct {
	name         string
	sub_commands []*BrigadierCommand
	executes     func(*BrigadierInvocation)
	arguments    map[string]BrigadierArgumentDefinition
	_arg_index   uint16
}

// Definition for an argument.
// Must be made with any of the argument helpers, such as [BrigadierCommand.String] or [BrigadierCommand.Number].
type BrigadierArgumentDefinition struct {
	name       string
	value_type string
	index      uint16
}

// The invocation of a Brigadier command.
type BrigadierInvocation struct {
	args         []string
	string_args  map[string]string
	number_args  map[string]int
	boolean_args map[string]bool
	User         *ChatboxIngameUser
	parent       *BrigadierCommand
	brigadier    *Brigadier
}

// Creates a new instance of Brigadier.
// Requires an instance of Chatbox passed, along with the name of the command.
// Returns a reference to a [Brigadier] struct.
func NewBrigadier(sc *Chatbox, name string) *Brigadier {
	b := &Brigadier{
		conn: sc,
		cmds: []*BrigadierCommand{},
		name: name,
	}

	sc.OnCommand = func(packet ChatboxCommandPacket) {
		var cmd *BrigadierCommand

		for _, value := range b.cmds {
			if value.name == packet.Command {
				cmd = value
			}
		}

		if cmd == nil {
			return
		}

		b.parse(cmd, packet)
	}

	return b
}

// Internal function to parse Command packets from Chatbox to Brigadier.
// Will recurse itself with subcommands.
func (b *Brigadier) parse(cmd *BrigadierCommand, packet ChatboxCommandPacket) {
	var target *BrigadierCommand

	if len(packet.Args) == 0 {
		target = cmd
	} else {
		for _, sub := range cmd.sub_commands {
			if sub.name == packet.Args[0] {
				b.parse(sub, ChatboxCommandPacket{
					Event:     packet.Event,
					User:      packet.User,
					Command:   packet.Args[0],
					Args:      packet.Args[1:],
					OwnerOnly: packet.OwnerOnly,
				})

				return
			}
		}

		if len(cmd.arguments) != 0 {
			target = cmd
		}
	}

	if target == nil || target.executes == nil {
		b.tellError(packet.User.Uuid, fmt.Sprintf("No subcommand or argument found. Check out &7\\%s help &cfor more information.", cmd.name))
		return
	}

	var string_args map[string]string = make(map[string]string)
	var number_args map[string]int = make(map[string]int)
	var bool_args map[string]bool = make(map[string]bool)

	for arg_name, arg := range target.arguments {
		if len(packet.Args) <= int(arg.index) {
			b.tellError(packet.User.Uuid, fmt.Sprintf("Missing argument \"%s\"", arg_name))
			return
		}
		str := packet.Args[arg.index]

		switch arg.value_type {
		case "string":
			if strings.HasPrefix(str, "'") || strings.HasPrefix(str, "\"") || strings.HasPrefix(str, "«") {
				closingArgument := -1

				for idx, future_arg := range packet.Args[arg.index:] {
					if strings.HasSuffix(future_arg, "'") || strings.HasSuffix(future_arg, "\"") || strings.HasSuffix(future_arg, "»") {
						closingArgument = idx
					}
				}

				if closingArgument != -1 {
					str = strings.Join(packet.Args[arg.index:closingArgument+1], " ")
					str = str[1 : len(str)-1]
				}
			}

			string_args[arg_name] = str
		case "number":
			num, err := strconv.Atoi(str)

			if err != nil {
				b.tellError(packet.User.Uuid, fmt.Sprintf("Unable to convert \"%s\" to number", str))
				return
			}

			number_args[arg_name] = num
		case "boolean":
			boolean, err := strconv.ParseBool(strings.ToLower(str))

			if err != nil {
				b.tellError(packet.User.Uuid, fmt.Sprintf("Unable to convert \"%s\" to boolean", str))
				return
			}

			bool_args[arg_name] = boolean
		}
	}

	target.executes(&BrigadierInvocation{
		parent:       cmd,
		User:         &packet.User,
		brigadier:    b,
		args:         packet.Args,
		string_args:  string_args,
		number_args:  number_args,
		boolean_args: bool_args,
	})
}

// Registers command(s) into this Brigadier instance
func (b *Brigadier) Register(commands ...*BrigadierCommand) {
	for _, command := range commands {
		err := command.verify()
		if err != nil {
			log.Panicf("error while parsing command %s: %s", command.name, err.Error())
		}

		command.Then(b.Literal("help").Executes(func(bi *BrigadierInvocation) {
			var content []string
			content = append(content, fmt.Sprintf("**\\%s Help Page**", command.name))
			content = append(content, command.getHelp("", 0)...)

			bi.ReplyMarkdown(strings.Join(content, "\n"))
		}))
	}

	b.cmds = append(b.cmds, commands...)
}

// Creates a new command with the prefix supplied.
// Returns a [BrigadierCommand] which you will use to make the command.
func (*Brigadier) Literal(command string) *BrigadierCommand {
	return &BrigadierCommand{
		name:         command,
		sub_commands: []*BrigadierCommand{},
		arguments:    map[string]BrigadierArgumentDefinition{},
		_arg_index:   0,
	}
}

func (cmd *BrigadierCommand) Then(commands ...*BrigadierCommand) *BrigadierCommand {
	cmd.sub_commands = append(cmd.sub_commands, commands...)
	return cmd
}

// Sets the function that will be ran once the command is triggered.
func (cmd *BrigadierCommand) Executes(fn func(*BrigadierInvocation)) *BrigadierCommand {
	cmd.executes = fn
	return cmd
}

// Internal helper function to push an argument definition.
func (cmd *BrigadierCommand) push_arg_def(arg_name, value_type string) {
	cmd.arguments[arg_name] = BrigadierArgumentDefinition{
		name:       arg_name,
		value_type: value_type,
		index:      cmd._arg_index,
	}

	cmd._arg_index++
}

// Verifies that the command you are trying to register, is valid.
func (cmd *BrigadierCommand) verify() error {
	// TODO: Verify that all literals have executes
	// TODO: Verify that the command name is alright

	return nil
}

// Defines a string argument with the supplied name.
func (cmd *BrigadierCommand) String(arg_name string) *BrigadierCommand {
	cmd.push_arg_def(arg_name, "string")
	return cmd
}

// Defines a number argument with the supplied name.
func (cmd *BrigadierCommand) Number(arg_name string) *BrigadierCommand {
	cmd.push_arg_def(arg_name, "number")
	return cmd
}

// Defines a boolean argument with the supplied name.
func (cmd *BrigadierCommand) Boolean(arg_name string) *BrigadierCommand {
	cmd.push_arg_def(arg_name, "boolean")
	return cmd
}

func (cmd *BrigadierCommand) getHelp(parentName string, depth int) []string {
	var lines []string

	if parentName == "" {
		lines = append(lines, fmt.Sprintf("`\\%s %s`", cmd.name, cmd.getArgsHelp()))
	} else {
		lines = append(lines, fmt.Sprintf("↪ %s `%s %s %s`", strings.Repeat(" ↪", depth-1), parentName, cmd.name, cmd.getArgsHelp()))
	}

	for _, sub := range cmd.sub_commands {
		if parentName == "" {
			lines = append(lines, strings.Join(sub.getHelp(fmt.Sprintf("\\%s", cmd.name), depth+1), "\n"))
		} else {
			lines = append(lines, strings.Join(sub.getHelp(fmt.Sprintf("%s %s", parentName, cmd.name), depth+1), "\n"))
		}
	}

	return lines
}

func (cmd *BrigadierCommand) getArgsHelp() string {
	out := ""

	for _, arg := range cmd.arguments {
		out += fmt.Sprintf("[%s: %s] ", arg.name, arg.value_type)
	}

	return out
}

// Replies to the user with the supplied message, in format mode.
// Use [BrigadierInvocation.ReplyMarkdown] to reply with markdown mode.
func (ev *BrigadierInvocation) Reply(message string) {
	ev.brigadier.conn.Tell(ev.User.Uuid, message, ev.brigadier.name, ChatboxFormattingFormat)
}

// Replies to the user with the supplied message, in markdown mode.
// Use [BrigadierInvocation.Reply] to reply with format mode.
func (ev *BrigadierInvocation) ReplyMarkdown(message string) {
	ev.brigadier.conn.Tell(ev.User.Uuid, message, ev.brigadier.name, ChatboxFormattingMarkdown)
}

// Replies to the uesr with an error.
// The message is prepended with "Error: ".
// Your message has to be in formatting mode, and is by default red.
func (ev *BrigadierInvocation) Error(message string) {
	ev.brigadier.tellError(ev.User.Uuid, fmt.Sprintf("&c&lError: &c%s", message))
}

// Internal version of [Error], that also requires the user uuid to send to.
func (b *Brigadier) tellError(user, message string) {
	b.conn.Tell(user, fmt.Sprintf("&c&lError: &c%s", message), b.name, ChatboxFormattingFormat)
}

// Internal function to validate that you are allowed to read the value you want.
// Panics if you cannot read it.
func (ev *BrigadierInvocation) validateRead(arg_name, arg_type string) {
	val, ok := ev.parent.arguments[arg_name]

	if !ok {
		panic(fmt.Sprintf("attempting to read nonexistant argument \"%s\" as %s", arg_name, arg_type))
	}

	if val.value_type != arg_type {
		panic(fmt.Sprintf("attempting to read argument \"%s\" as %s, but type is defined as %s", arg_name, arg_type, val.value_type))
	}
}

// Function to read a string defined with the [BrigadierCommand.String function.
// Can panic if you are attempting to read a non-existing argument.
func (ev *BrigadierInvocation) ReadString(arg_name string) string {
	ev.validateRead(arg_name, "string")
	return ev.string_args[arg_name]
}

// Function to read a number defined with the [BrigadierCommand.Number] function.
// Can panic if you are attempting to read a non-existing argument.
func (ev *BrigadierInvocation) ReadNumber(arg_name string) int {
	ev.validateRead(arg_name, "number")
	return ev.number_args[arg_name]
}

// Function to read a boolean defined with the [BrigadierCommand.Boolean] function.
// Can panic if you are attempting to read a non-existing argument.
func (ev *BrigadierInvocation) ReadBoolean(arg_name string) bool {
	ev.validateRead(arg_name, "boolean")
	return ev.boolean_args[arg_name]
}
