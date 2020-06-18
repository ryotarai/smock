package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/ryotarai/smock/pkg/client"
	"github.com/slack-go/slack"
)

var slashCommandRegexp = regexp.MustCompile(`(?s)\A(/[^\s]+)(.*)\z`)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type CommandHandler func(command, text string) error

type CLI struct {
	Client *client.Client
}

func New() *CLI {
	return &CLI{}
}

func (c *CLI) Start() {
	p := prompt.New(c.execute, c.completer,
		prompt.OptionPrefix(">>> "))
	p.Run()
}

func (c *CLI) execute(input string) {
	if err := c._execute(input); err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}
}

func (c *CLI) completer(in prompt.Document) []prompt.Suggest {
	return nil
}

func (c *CLI) _execute(input string) error {
	if input == "" {
		return nil
	}

	if match := slashCommandRegexp.FindStringSubmatch(input); match != nil {
		command := match[1]
		text := match[2]

		text = strings.TrimLeft(text, " ")

		msg, err := c.Client.SendCommand(command, text)
		if err != nil {
			return err
		}

		c.OnMessage(msg)
	} else if input == "exit" {
		// TODO: care defer functions
		os.Exit(0)
	} else {
		if err := c.Client.SendMessage(input); err != nil {
			return err
		}
	}

	return nil
}

func (c *CLI) OnMessage(msg *slack.Msg) {
	if msg.Text != "" {
		fmt.Printf("\033[2K\033[1000D%s %s\n", Yellow("<<<"), msg.Text)
	}
}
