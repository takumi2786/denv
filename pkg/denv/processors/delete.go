package processors

import (
	"fmt"
	"log/slog"
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
	logger *slog.Logger
}

var _ Processor = (*DeleteProcessor)(nil)

func NewDeleteProcessor(logger *slog.Logger) *DeleteProcessor {
	return &DeleteProcessor{logger: logger}
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

	dp.logger.Info("DeleteProcessor", slog.Any("options", dOptions))

	dp.logger.Info("Deleting Container...")
	// create command
	exCmd := exec.Command(
		"docker", "rm", "-f", dOptions.Identity,
	)

	// Bind input/output to parent process terminal
	exCmd.Stdout = stdout
	exCmd.Stderr = stderr

	// Ececute command
	if err := exCmd.Run(); err != nil {
		fmt.Println("Failed to delete container:", err)
		return err
	}
	dp.logger.Info("The Container is Deleted.")

	return nil
}
