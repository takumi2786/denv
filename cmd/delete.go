/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/m-mizutani/goerr"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete selected container",
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := parseDeleteCmd(cmd)
		if err != nil {
			goerr.New("Failed to Parse Command")
		}

		err = delete(options)
		if err != nil {
			fmt.Println("Error Occured in run.", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("identity", "i", "ubuntu", "Docker Image identity defined in image_map.json")
}

type DeleteOptions struct {
	ImageMapPath string
	Identity     string
}

func (o *DeleteOptions) String() string {
	return fmt.Sprintf(
		"DeleteOptions: ImageMapPath: %s, Identity: %s", o.ImageMapPath, o.Identity,
	)
}

func parseDeleteCmd(cmd *cobra.Command) (*DeleteOptions, error) {
	identity, err := cmd.Flags().GetString("identity")
	if err != nil {
		return nil, err
	}

	filepath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	return &DeleteOptions{
		Identity:     identity,
		ImageMapPath: filepath,
	}, err
}

func delete(options *DeleteOptions) error {
	if options == nil {
		return goerr.New("InternalError: options is nil")
	}

	fmt.Println("Deleting Container...", options)

	// create command
	exCmd := exec.Command(
		"docker", "rm", "-f", options.Identity,
	)

	// 入出力を親プロセスのターミナルにバインド
	exCmd.Stdin = os.Stdin
	exCmd.Stdout = os.Stdout
	exCmd.Stderr = os.Stderr

	// 実行
	if err := exCmd.Run(); err != nil {
		fmt.Println("コンテナの削除に失敗:", err)
	}

	return nil
}
