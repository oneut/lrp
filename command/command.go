package command

import (
	"bufio"
	"github.com/mattn/go-shellwords"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/log"
	"os/exec"
	"strings"
	"syscall"
)

func NewCommand(name string, commandConfig config.Command) *Command {
	log.Info(name, "Start command")
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
	c.Cmd.Start()
	oneByte := make([]byte, 100)
	go func() {
		for {
			_, err := stdout.Read(oneByte)
			if err != nil {
				break
			}
			r := bufio.NewReader(stdout)
			line, _, _ := r.ReadLine()
			s := string(line)
			log.Info(c.Name, s)

			if len(c.CommandConfig.WatchStdout) == 0 {
				continue
			}

			for _, value := range c.CommandConfig.WatchStdout {
				if strings.Contains(s, value) {
					c.Callback("stdout notify")
				}
			}
		}
		c.Cmd.Wait()
	}()
}

func (c *Command) Restart() {
	if !(c.CommandConfig.NeedsRestart) {
		c.Kill()
		c.Start()
	}
}

func (c *Command) Kill() {
	if c.Cmd.Process == nil {
		return
	}

	c.Cmd.Process.Signal(syscall.SIGTERM)
}
