package command

import (
	"bufio"
	"os/exec"
	"strings"
	"syscall"

	"github.com/mattn/go-shellwords"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/logger"
)

func NewCommand(name string, commandName string, commandConfig config.Command) Commander {
	if !(commandConfig.IsValid()) {
		return &NilCommand{}
	}

	return &Command{
		CommandConfig: commandConfig,
		CommandName:   commandName,
		Name:          name,
	}
}

type Commander interface {
	Run(func(string))
	Start()
	Restart()
	Kill() bool
	Stop()
}

type Command struct {
	Cmd           *exec.Cmd
	CommandConfig config.Command
	CommandName   string
	Name          string
	Callback      func(string)
}

func (c *Command) Run(fn func(string)) {
	logger.InfoCommand(c.Name, c.CommandName, "start")
	c.Callback = fn
	c.Start()
}

func (c *Command) Start() {
	args, err := shellwords.Parse(c.CommandConfig.Execute)
	if err != nil {
		panic(err)
	}

	switch len(args) {
	case 0:
		panic("command.execute is required")
	case 1:
		c.Cmd = exec.Command(args[0])
	default:
		c.Cmd = exec.Command(args[0], args[1:]...)
	}

	stdout, _ := c.Cmd.StdoutPipe()
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			logger.InfoCommandStdout(c.Name, c.CommandName, line)
			c.watchStdout(line)
		}
	}()

	defer c.Kill()
	c.Cmd.Start()
	c.Cmd.Wait()
}

func (c *Command) watchStdout(line string) {
	if len(c.CommandConfig.WatchStdout) == 0 {
		return
	}

	for _, value := range c.CommandConfig.WatchStdout {
		if strings.Contains(line, value) {
			c.Callback("stdout notify")
			logger.InfoCommand(c.Name, c.CommandName, "watch_stdout is fired:"+value)
			break
		}
	}
}

func (c *Command) Restart() {
	if c.CommandConfig.NeedsRestart {
		c.Kill()
		c.Start()
	}
}

func (c *Command) Kill() bool {
	if c.Cmd.Process == nil {
		return false
	}

	c.Cmd.Process.Signal(syscall.SIGTERM)

	return true
}

func (c *Command) Stop() {
	if c.Kill() {
		logger.InfoCommand(c.Name, c.CommandName, "stop")
	}
}
