// +build !windows

package command

import (
	"syscall"
)

func (c *Command) sysProcAttr() {
	c.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func (c *Command) kill() {
	// kill process with child process.
	pgid, err := syscall.Getpgid(c.cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGKILL)
	}
}
