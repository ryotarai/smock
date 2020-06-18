package command

import (
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ryotarai/smock/pkg/cli"
	"github.com/ryotarai/smock/pkg/client"
	"github.com/ryotarai/smock/pkg/server"
	urfavecli "github.com/urfave/cli/v2"
)

const (
	flagListen      = "listen"
	flagServerLog   = "server-log"
	flagEventURL    = "event-url"
	flagExternalURL = "external-url"
	flagUserID      = "user-id"
	flagUserName    = "user-name"
)

type Command struct {
}

func New() *Command {
	return &Command{}
}

func (c *Command) App() *urfavecli.App {
	return &urfavecli.App{
		Commands: []*urfavecli.Command{
			{
				Name:   "start",
				Action: c.actionStart,
				Flags: []urfavecli.Flag{
					&urfavecli.StringFlag{
						Name:  flagListen,
						Usage: "server listen on",
						Value: ":8000",
					},
					&urfavecli.StringFlag{
						Name:  flagServerLog,
						Usage: "server log file path",
					},
					&urfavecli.StringFlag{
						Name:     flagEventURL,
						Usage:    "event subscription URL",
						Required: true,
					},
					&urfavecli.StringFlag{
						Name:     flagExternalURL,
						Usage:    "external URL",
						Required: true,
					},
					&urfavecli.StringFlag{
						Name:  flagUserID,
						Usage: "user ID",
						Value: "USERID",
					},
					&urfavecli.StringFlag{
						Name:  flagUserName,
						Usage: "user name",
						Value: "USERNAME",
					},
				},
			},
		},
	}
}

func (c *Command) actionStart(ctx *urfavecli.Context) error {
	// Setup client
	client := client.New()
	client.EventURL = ctx.String(flagEventURL)
	client.ExternalURL = ctx.String(flagExternalURL)
	client.UserName = ctx.String(flagUserName)
	client.UserID = ctx.String(flagUserID)

	// Setup interface
	cli := cli.New()
	cli.Client = client

	// Setup server
	if serverLog := ctx.String(flagServerLog); serverLog != "" {
		logFile, err := os.OpenFile(serverLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer logFile.Close()
		gin.DefaultWriter = logFile
	} else {
		gin.DefaultWriter = ioutil.Discard
	}

	server := server.New()
	server.CLI = cli
	go func() {
		if err := server.Run(":3002"); err != nil {
			panic(err)
		}
	}()

	// Start loop
	cli.Start()

	return nil
}
