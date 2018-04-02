package command

import (
	"os/exec"
	"strconv"
	"syscall"
)

func (c *Command) sysProcAttr() {
	c.cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func (c *Command) kill() {
	// kill process with child process.
	cmd := exec.Command("taskkill", "/f", "/t", "/pid", strconv.Itoa(c.cmd.Process.Pid))
	cmd.Start()
	cmd.Wait()
}
