// /*
// Copyright © 2025 NAME HERE <EMAIL ADDRESS>
// */
package cmd

import (
	"fmt"
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/spf13/cobra"
	"github.com/takumi2786/denv/pkg/denv"
	"github.com/takumi2786/denv/pkg/denv/processors"
)

func parseRunCmd(cmd *cobra.Command) (*processors.RunOptions, error) {
	identity, err := cmd.Flags().GetString("identity")
	if err != nil {
		return nil, err
	}

	filepath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	return &processors.RunOptions{
		Identity:     identity,
		ImageMapPath: filepath,
	}, err
}

// runCmd strt container and attach
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start Instant Container",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := denv.NewLogger()

		options, err := parseRunCmd(cmd)
		if err != nil {
			goerr.New("Failed to Parse Command")
		}

		processor := processors.NewRunProcessor(logger)
		err = processor.Run(options, os.Stdin, os.Stdout, os.Stderr)
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
