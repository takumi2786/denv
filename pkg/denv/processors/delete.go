package processors

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/m-mizutani/goerr"
)

type DeleteOptions struct {
	ImageMapPath string
	Identity     string
}

func (o *DeleteOptions) String() string {
	return fmt.Sprintf(
		"DeleteOptions: ImageMapPath: %s, Identity: %s", o.ImageMapPath, o.Identity,
	)
}

func Delete(options *DeleteOptions, stdin *os.File, stdout *os.File, stderr *os.File) error {
	if options == nil {
		return goerr.New("InternalError: options is nil")
	}

	fmt.Println("Deleting Container...", options)

	// create command
	exCmd := exec.Command(
		"docker", "rm", "-f", options.Identity,
	)

	// 入出力を親プロセスのターミナルにバインド
	exCmd.Stdin = stdin
	exCmd.Stdout = stdout
	exCmd.Stderr = stderr

	// 実行
	if err := exCmd.Run(); err != nil {
		fmt.Println("Failed to delete container:", err)
		return err
	}

	return nil
}
