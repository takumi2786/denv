/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/spf13/cobra"
	"github.com/takumi2786/denv/pkg/denv"
	"github.com/takumi2786/denv/pkg/denv/processors"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete selected container",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := denv.NewLogger()

		options, err := parseDeleteCmd(cmd)
		if err != nil {
			goerr.New("Failed to Parse Command")
		}

		processor := processors.NewDeleteProcessor(logger)
		processor.Run(options, os.Stdin, os.Stdout, os.Stderr)
		if err != nil {
			logger.Error("Error Occured in Run.", slog.Any("err", err))
		}
		return nil
	},
}

// init initialize comand line args
func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("identity", "i", "ubuntu", "Docker Image identity defined in image_map.json")
}

func parseDeleteCmd(cmd *cobra.Command) (*processors.DeleteOptions, error) {
	identity, err := cmd.Flags().GetString("identity")
	if err != nil {
		return nil, err
	}

	filepath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	return &processors.DeleteOptions{
		Identity:     identity,
		ImageMapPath: filepath,
	}, err
}
