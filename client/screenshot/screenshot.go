package screenshot

import (
	"fmt"
	"os/exec"
)

// Program is a simple struct defining a screenshot program
type Program struct {
	Command        string
	WindowArgs     []string
	FullScreenArgs []string
}

// GetCommand returns an *exec.Cmd
func (p *Program) GetCommand(fp string, window bool) *exec.Cmd {
	var args []string
	if window {
		args = append(args, p.WindowArgs...)
		args = append(args, fp)
	} else {
		args = append(args, p.FullScreenArgs...)
		args = append(args, fp)
	}
	return exec.Command(p.Command, args...)
}

// Do executes a screenshot with the given program
func (p *Program) Do(fp string, window bool) error {
	cmd := p.GetCommand(fp, window)
	cmd.Stderr = nil
	cmd.Stdout = nil
	return cmd.Run()
}

var progs = []Program{
	{
		Command:        "gnome-screenshot",
		WindowArgs:     []string{"-a", "-f"},
		FullScreenArgs: []string{"-f"},
	},
	{
		Command:        "import",
		WindowArgs:     []string{},
		FullScreenArgs: []string{"-window", "root"},
	},
	{
		Command:        "scrot",
		WindowArgs:     []string{"-s"},
		FullScreenArgs: []string{},
	},
}

// Do executes the screenshot, saving the file to the file path fp.
func Do(fp string, window bool) error {
	var err error

	for _, p := range progs {
		if _, err = exec.LookPath(p.Command); err == nil {
			return p.Do(fp, window)
		}
	}

	return fmt.Errorf("no screenshot program found")
}
