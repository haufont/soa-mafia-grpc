package cli

import (
	"context"
	"fmt"
	"mafia-grpc/internal/client"
	"mafia-grpc/internal/player"
	"mafia-grpc/internal/state"
	"mafia-grpc/internal/util"
	"os"
	"strings"

	"github.com/fatih/color"
)

func GetUsername() string {
	var username string
	for {
		fmt.Println("Enter your name: ")
		fmt.Scanf("%s", &username)
		err := player.ValidatePlayerName(username)
		if err == nil {
			break
		}
		fmt.Println(err)
	}
	return username
}

type Command struct {
	Name        string
	Args        []string
	Description string
	Handler     func() (string, bool)
}

type CLI struct {
	ctx      context.Context
	waiter   util.Waiter
	client   *client.Client
	commands []Command
}

func NewCLI() CLI {
	ctx, waiter := util.NewStopCtx()
	go func() {
		waiter.Wait()
	}()
	return CLI{
		ctx:    ctx,
		waiter: waiter,
	}
}

func (c CLI) GetUsername() string {
	var username string
	for {
		select {
		case <-c.ctx.Done():
			c.waiter.Cancel()
			os.Exit(0)
		default:
			fmt.Print("Enter your name: ")
			fmt.Scanf("%s", &username)
			err := player.ValidatePlayerName(username)
			if err == nil {
				return username
			}
			color.Red("%s\n", err)
		}
	}
}

func (c *CLI) StartClient(options client.ClientOptions) error {
	client, err := client.NewClient(c.GetUsername(), c.ctx, options)
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

func (c *CLI) HandleStateCommand() (rsp string, cont bool) {
	s, err := c.client.GetState()
	if err != nil {
		return fmt.Sprint(err), true
	}
	return state.StateToString(s, util.BytesToUUID(c.client.Player.Uuid)), true
}

func (c *CLI) HandleLeaveCommand() (rsp string, cont bool) {
	return "", false
}

func (c *CLI) HandleKillCommand() (rsp string, cont bool) {
	var playerName string
	fmt.Scanf("%s", &playerName)
	err := c.client.Kill(playerName)
	if err != nil {
		return color.RedString("%s", err), true
	}
	return "", true
}

func (c *CLI) HandleSkipCommand() (rsp string, cont bool) {
	err := c.client.Kill("")
	if err != nil {
		return color.RedString("%s", err), true
	}
	return "", true
}

func (c *CLI) HandleCheckCommand() (rsp string, cont bool) {
	var playerName string
	fmt.Scanf("%s", &playerName)
	isMafia, err := c.client.Check(playerName)
	if err != nil {
		return color.RedString("%s", err), true
	}
	if isMafia {
		return color.RedString("%s is a mafia", playerName), true
	} else {
		return color.GreenString("%s isn't a mafia", playerName), true
	}
}

func (c *CLI) HandleSkipCheckCommand() (rsp string, cont bool) {
	_, err := c.client.Check("")
	if err != nil {
		return color.RedString("%s", err), true
	}
	return "", true
}

func (c *CLI) HandlePublishCommand() (rsp string, cont bool) {
	var playerName string
	fmt.Scanf("%s", &playerName)
	err := c.client.Publish(playerName)
	if err != nil {
		return color.RedString("%s", err), true
	}
	return "", true
}

func (c *CLI) HandleHelpCommand() (rsp string, cont bool) {
	var descriptions []string
	descriptions = append(descriptions, "Commands: ")
	for _, command := range c.commands {
		args := ""
		if len(command.Args) != 0 {
			args = fmt.Sprintf("[%s] ", strings.Join(command.Args, ", "))
		}
		descriptions = append(descriptions, fmt.Sprintf("%s %s- %s", command.Name, args, command.Description))
	}
	return color.GreenString(strings.Join(descriptions, "\n")), true
}

func (c *CLI) HandleCommand() (cont bool) {
	var scommand string
	fmt.Scanf("%s", &scommand)

	rsp := color.RedString(fmt.Sprintf("Unknown command \"%s\". Use \"%s\" for details", scommand, HelpCommand))
	cont = true
	for _, command := range c.commands {
		if scommand == command.Name {
			rsp, cont = command.Handler()
			break
		}
	}

	if rsp != "" {
		fmt.Printf("%s\n", rsp)
	}

	return cont
}

const (
	StateCommand     string = "state"
	LeaveCommand     string = "leave"
	KillCommand      string = "kill"
	SkipCommand      string = "skip"
	CheckCommand     string = "check"
	SkipCheckCommand string = "skip-check"
	PublishCommand   string = "publish"
	HelpCommand      string = "help"
)

func (c *CLI) addCommands() {
	c.commands = append(c.commands,
		Command{
			Name:        StateCommand,
			Description: "get the state of the game",
			Handler:     c.HandleStateCommand,
		},
		Command{
			Name:        LeaveCommand,
			Description: "log out of the session",
			Handler:     c.HandleLeaveCommand,
		},
		Command{
			Name:        KillCommand,
			Args:        []string{"player name"},
			Description: "vote for execution/murder (available to everyone during the day or only to the mafia at night)",
			Handler:     c.HandleKillCommand,
		},
		Command{
			Name:        SkipCommand,
			Description: "abstain from voting at this stage (equivalent to an empty vote)",
			Handler:     c.HandleSkipCommand,
		},
		Command{
			Name:        CheckCommand,
			Description: "checks if the player is a mafia (available only to the commissioner at night)",
			Handler:     c.HandleCheckCommand,
		},
		Command{
			Name:        SkipCheckCommand,
			Description: "abstain from check (available only to the commissioner at night)",
			Handler:     c.HandleSkipCheckCommand,
		},
		Command{
			Name:        PublishCommand,
			Description: "publish information about the mafia that the commissioner checked (available only to the commissioner at day)",
			Handler:     c.HandlePublishCommand,
		},
		Command{
			Name:        HelpCommand,
			Description: "information about commands",
			Handler:     c.HandleHelpCommand,
		},
	)
}

func (c *CLI) Serve() {
	c.addCommands()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if !c.HandleCommand() {
				c.waiter.Cancel()
				return
			}
		}
	}
}
