package command

import (
	"bufio"
	"github.com/mattn/go-shellwords"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/logger"
	"os/exec"
	"strings"
	"syscall"
)

func NewCommand(name string, commandConfig config.Command) *Command {
	return &Command{
		CommandConfig: commandConfig,
		Name:          name,
	}
}

type Command struct {
	Cmd           *exec.Cmd
	CommandConfig config.Command
	Name          string
	Callback      func(string)
}

func (c *Command) Run(fn func(string)) {
	if c.isUndefined() {
		return
	}

	logger.InfoEvent(c.Name, "command", "start")
	c.Callback = fn
	c.Start()
}

func (c *Command) Start() {
	if c.isUndefined() {
		return
	}

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
			logger.InfoEvent(c.Name, "command", line)
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
	if c.isUndefined() {
		return false
	}

	if c.Cmd.Process == nil {
		return false
	}

	c.Cmd.Process.Signal(syscall.SIGTERM)

	return true
}

func (c *Command) Stop() {
	if c.isUndefined() {
		return
	}

	if c.Kill() {
		logger.InfoEvent(c.Name, "command", "stop")
	}
}

func (c *Command) isUndefined() bool {
	if c.CommandConfig.Execute == "" {
		return true
	}

	return false
}
