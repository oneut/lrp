package command

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/logger"
)

func NewCommand(name string, commandName string, commandConfig config.Command) CommandInterface {
	if !(commandConfig.IsValid()) {
		return &NilCommand{}
	}

	return &Command{
		commandName:  commandName,
		name:         name,
		executes:     commandConfig.Executes,
		needsRestart: commandConfig.NeedsRestart,
		watchStdouts: commandConfig.WatchStdouts,
	}
}

type CommandInterface interface {
	Run(func(string))
	Start()
	NeedsRestart() bool
	Kill() bool
	Stop()
}

type Command struct {
	cmd          *exec.Cmd
	commandName  string
	name         string
	callback     func(string)
	executes     []string
	needsRestart bool
	watchStdouts []string
}

func (c *Command) Run(fn func(string)) {
	logger.InfoCommand(c.name, c.commandName, "start")
	c.callback = fn
	c.Start()
}

func (c *Command) Start() {
	for _, execute := range c.executes {
		args, err := shellwords.Parse(execute)
		if err != nil {
			panic(err)
		}

		switch len(args) {
		case 0:
			panic("command.execute is required")
		case 1:
			c.cmd = exec.Command(args[0])
		default:
			c.cmd = exec.Command(args[0], args[1:]...)
		}

		stdout, _ := c.cmd.StdoutPipe()
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				logger.InfoCommandStdout(c.name, c.commandName, line)
				c.watchStdout(line)
			}
		}()

		stderr, _ := c.cmd.StderrPipe()
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				logger.InfoCommandStdout(c.name, c.commandName, line)
			}
		}()

		defer c.Kill()
		c.sysProcAttr()
		c.cmd.Start()
		logger.InfoCommand(c.name, c.commandName, "execute: "+execute)
		c.cmd.Wait()
	}
}

func (c *Command) watchStdout(line string) {
	if len(c.watchStdouts) == 0 {
		return
	}

	for _, value := range c.watchStdouts {
		if strings.Contains(line, value) {
			c.callback("stdout notify")
			logger.InfoCommand(c.name, c.commandName, "watch_stdouts is fired:"+value)
			break
		}
	}
}

func (c *Command) NeedsRestart() bool {
	return c.needsRestart
}

func (c *Command) Kill() bool {
	if c.cmd.Process == nil {
		return false
	}

	c.kill()
	return true
}

func (c *Command) Stop() {
	if c.Kill() {
		logger.InfoCommand(c.name, c.commandName, "stop")
	}
}
