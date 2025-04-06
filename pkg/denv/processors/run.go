package processors

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/takumi2786/denv/pkg/denv"
)

type RunOptions struct {
	ImageMapPath string
	Identity     string
}

func (o *RunOptions) String() string {
	return fmt.Sprintf(
		"RunOptions: ImageMapPath: %s, Identity: %s", o.ImageMapPath, o.Identity,
	)
}

func Run(options *RunOptions, stdin *os.File, stdout *os.File, stderr *os.File) error {
	if options == nil {
		return goerr.New("InternalError: options is nil")
	}
	fmt.Println("Starting Container...", options)
	reader := denv.NewImageMapReader()
	err := reader.Read(options.ImageMapPath)
	if err != nil {
		return goerr.Wrap(err, "Faled to Parse image map")
	}

	entry, err := reader.Loadded(options.Identity)
	if err != nil {
		return goerr.Wrap(err, "Faled to Load image map")
	}

	// args
	commandArgs := []string{
		"run",
		"-itd",
		"-v", ".:/Workspace",
		"--name", options.Identity,
		"-w", "/Workspace",
	}
	// option
	commandArgs = append(commandArgs, strings.Split(entry.Option, " ")...)
	// image uri
	commandArgs = append(commandArgs, entry.ImageURI)
	// entrypoint
	if entry.Entrypoint != "" {
		commandArgs = append(commandArgs, entry.Entrypoint)
	}
	// cmd
	if entry.Cmd != "" {
		commandArgs = append(commandArgs, entry.Cmd)
	}
	// create command
	exCmdStart := exec.Command(
		"docker",
		commandArgs...,
	)

	// 入出力を親プロセスのターミナルにバインド
	exCmdStart.Stdin = stdin
	exCmdStart.Stdout = stdout
	exCmdStart.Stderr = stderr

	// 実行
	if err := exCmdStart.Run(); err != nil {
		fmt.Println("コンテナの起動に失敗:", err)
		return err
	}

	exCmd := exec.Command("docker", "exec", "-it", options.Identity, entry.Shell)

	// 入出力を親プロセスのターミナルにバインド
	exCmd.Stdin = stdin
	exCmd.Stdout = stdout
	exCmd.Stderr = stderr

	// 実行
	if err := exCmd.Run(); err != nil {
		fmt.Println("実行エラー:", err)
	}

	return nil
}
