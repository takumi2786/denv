// /*
// Copyright © 2025 NAME HERE <EMAIL ADDRESS>
// */
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/spf13/cobra"
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

func parseRunCmd(cmd *cobra.Command) (*RunOptions, error) {
	identity, err := cmd.Flags().GetString("identity")
	if err != nil {
		return nil, err
	}

	filepath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	return &RunOptions{
		Identity:     identity,
		ImageMapPath: filepath,
	}, err
}

// greetCmd represents the greet command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start Instant Container",
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := parseRunCmd(cmd)
		if err != nil {
			goerr.New("Failed to Parse Command")
		}

		err = run(options)
		if err != nil {
			fmt.Println("Error Occured in run.", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("identity", "i", "ubuntu", "Docker Image identity defined in image_map.json")
}

func run(options *RunOptions) error {
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
	exCmdStart.Stdin = os.Stdin
	exCmdStart.Stdout = os.Stdout
	exCmdStart.Stderr = os.Stderr

	// 実行
	if err := exCmdStart.Run(); err != nil {
		fmt.Println("コンテナの起動に失敗:", err)
		return err
	}

	exCmd := exec.Command("docker", "exec", "-it", options.Identity, entry.Shell)

	// 入出力を親プロセスのターミナルにバインド
	exCmd.Stdin = os.Stdin
	exCmd.Stdout = os.Stdout
	exCmd.Stderr = os.Stderr

	// 実行
	if err := exCmd.Run(); err != nil {
		fmt.Println("実行エラー:", err)
	}

	return nil
}
