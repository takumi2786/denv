package processors

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/takumi2786/denv/pkg/denv"
)

// RunOptions is parameters for RunProcessor
type RunOptions struct {
	ImageMapPath string
	Identity     string
}

func (o *RunOptions) String() string {
	return fmt.Sprintf(
		"RunOptions: ImageMapPath: %s, Identity: %s", o.ImageMapPath, o.Identity,
	)
}

// RunProcessor is Processor to run container
type RunProcessor struct {
}

var _ Processor = (*RunProcessor)(nil)

func NewRunProcessor() *RunProcessor {
	return &RunProcessor{}
}

// Run deletes selected container
func (rp *RunProcessor) Run(
	options any,
	stdin *os.File,
	stdout *os.File,
	stderr *os.File,
) error {
	if options == nil {
		return goerr.New("InternalError: options is nil")
	}

	// convert any to RunOptions
	rOptions, ok := options.(*RunOptions)
	if !ok {
		return goerr.New("failed to parse options")
	}

	fmt.Println("Starting Container...", rOptions)

	reader := denv.NewImageMapReader()
	err := reader.Read(rOptions.ImageMapPath)
	if err != nil {
		return goerr.Wrap(err, "Faled to Parse image map")
	}

	entry, err := reader.Loadded(rOptions.Identity)
	if err != nil {
		return goerr.Wrap(err, "Faled to Load image map")
	}

	/*
		Start container
	*/
	commandArgs := []string{
		"run",
		"-itd",
		"-v", ".:/Workspace",
		"--name", rOptions.Identity,
		"-w", "/Workspace",
	}
	commandArgs = append(commandArgs, strings.Split(entry.Option, " ")...)
	commandArgs = append(commandArgs, entry.ImageURI)
	if entry.Entrypoint != "" {
		commandArgs = append(commandArgs, entry.Entrypoint)
	}
	if entry.Cmd != "" {
		commandArgs = append(commandArgs, entry.Cmd)
	}
	// create command
	exCmdStart := exec.Command(
		"docker",
		commandArgs...,
	)

	// Bind input/output to parent process terminal
	exCmdStart.Stdin = stdin
	exCmdStart.Stdout = stdout
	exCmdStart.Stderr = stderr

	// Execute command
	if err := exCmdStart.Run(); err != nil {
		fmt.Println("Failed to Start Container:", err)
		return err
	}

	/*
		Attach container
	*/

	exCmd := exec.Command("docker", "exec", "-it", rOptions.Identity, entry.Shell)

	// Bind input/output to parent process terminal
	exCmd.Stdin = stdin
	exCmd.Stdout = stdout
	exCmd.Stderr = stderr

	// Execute command
	if err := exCmd.Run(); err != nil {
		fmt.Println("Failed to attach container:", err)
	}

	return nil
}
