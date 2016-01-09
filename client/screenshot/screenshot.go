package screenshot

import (
	"fmt"
	"os/exec"
)

// Do executes the screenshot, saving the file to the file path fp.
func Do(fp string, window bool) error {
	_, err := exec.LookPath("import")
	if err != nil {
		return fmt.Errorf("Command 'import' is not installed.")
	}
	var cmd *exec.Cmd
	if window {
		cmd = exec.Command("import", fp)
	} else {
		cmd = exec.Command("import", "-window", "root", fp)
	}
	cmd.Stdout = nil
	cmd.Stderr = nil

	err = cmd.Run()
	return err
}
