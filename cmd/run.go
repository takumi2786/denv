// /*
// Copyright © 2025 NAME HERE <EMAIL ADDRESS>
// */
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/m-mizutani/goerr"
	"github.com/spf13/cobra"
	"github.com/takumi2786/denv/pkg/denv"
)

type Options struct {
	ImageMapPath string
	Identity     string
}

func parseCmd(cmd *cobra.Command) (*Options, error) {
	identity, err := cmd.Flags().GetString("identity")
	if err != nil {
		return nil, err
	}

	filepath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	return &Options{
		Identity:     identity,
		ImageMapPath: filepath,
	}, err
}

// greetCmd represents the greet command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start Instant Container",
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := parseCmd(cmd)
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

func run(options *Options) error {
	if options == nil {
		return goerr.New("InternalError: options is nil")
	}
	reader := denv.NewImageMapReader()
	err := reader.Read(options.ImageMapPath)
	if err != nil {
		return goerr.Wrap(err, "Faled to Parse image map")
	}

	entry, err := reader.Loadded(options.Identity)
	if err != nil {
		return goerr.Wrap(err, "Faled to Load image map")
	}

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	entrypoint := []string{}
	if entry.Entrypoint != "" {
		entrypoint = append(entrypoint, entry.Entrypoint)
	}
	cmd := []string{}
	if entry.Cmd != "" {
		cmd = append(cmd, entry.Cmd)
	}
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        entry.ImageURI,
			Entrypoint:   strslice.StrSlice(entrypoint),
			Cmd:          strslice.StrSlice(cmd),
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			AttachStdin:  true,
			OpenStdin:    true,
			WorkingDir:   "/Workspace",
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/home/takumi/denv",
					Target: "/Workspace",
				},
			},
		},
		nil,
		nil,
		options.Identity,
	)
	if err != nil {
		return err
	}

	// コンテナ開始
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	exCmd := exec.Command("docker", "exec", "-it", options.Identity, "zsh")

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
