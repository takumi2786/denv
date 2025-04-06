package processors

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/m-mizutani/goerr"
)

// DeleteOptions is parameters used in DeleteProcessor
type DeleteOptions struct {
	ImageMapPath string
	Identity     string
}

func (o *DeleteOptions) String() string {
	return fmt.Sprintf(
		"DeleteOptions: ImageMapPath: %s, Identity: %s", o.ImageMapPath, o.Identity,
	)
}

// DeleteProcessor is Processor to delete container
type DeleteProcessor struct {
}

var _ Processor = (*DeleteProcessor)(nil)

func NewDeleteProcessor() *DeleteProcessor {
	return &DeleteProcessor{}
}

// Run deletes selected container
func (dp *DeleteProcessor) Run(
	options any,
	stdin *os.File,
	stdout *os.File,
	stderr *os.File,
) error {
	if options == nil {
		return goerr.New("InternalError: options is nil")
	}

	// convert any to DeleteOptions
	dOptions, ok := options.(*DeleteOptions)
	if !ok {
		return goerr.New("failed to parse options")
	}

	fmt.Println("Deleting Container...", dOptions)

	// create command
	exCmd := exec.Command(
		"docker", "rm", "-f", dOptions.Identity,
	)

	// Bind input/output to parent process terminal
	exCmd.Stdin = stdin
	exCmd.Stdout = stdout
	exCmd.Stderr = stderr

	// Ececute command
	if err := exCmd.Run(); err != nil {
		fmt.Println("Failed to delete container:", err)
		return err
	}

	return nil
}
