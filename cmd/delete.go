/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/spf13/cobra"
	"github.com/takumi2786/denv/pkg/denv/processors"
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

		err = processors.Delete(options, os.Stdin, os.Stdout, os.Stderr)
		if err != nil {
			fmt.Println("Error Occured in run.", err)
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
