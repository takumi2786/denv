/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/m-mizutani/goerr"
	"github.com/spf13/cobra"
	"github.com/takumi2786/denv/pkg/denv"
	"golang.org/x/term"
)

type Options struct {
	ImageMapPath string
	Identity     string
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

	// アタッチ
	attachResp, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	defer attachResp.Close()

	// raw mode に入る
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	defer term.Restore(int(syscall.Stdin), oldState) // 終了時に戻す

	// 入出力を繋げる
	go func() {
		_, _ = io.Copy(attachResp.Conn, os.Stdin) // キーボード → コンテナ
	}()
	go func() {
		_, _ = io.Copy(os.Stdout, attachResp.Reader) // コンテナ出力 → 画面
		// _, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, attachResp.Reader) // tty trueだとこれはダメ
		// ttyが有効の場合、stdout/stderr は1本のストリームに multiplex されてる。stdcopy.StdCopyはこれを分離して扱うためのメソッド。
	}()

	fmt.Printf("コンテナ %s にアタッチしました（終了するには Ctrl+D)\n", resp.ID)

	// 終了まで待機
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
		fmt.Println("コンテナ終了")
	}
	return nil
}
