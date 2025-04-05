package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "denv",
	Long: "即席のDockerコンテナを立ち上げ、アタッチする",
	// メイン処理
	RunE: func(cmd *cobra.Command, args []string) error {
		identity, err := cmd.Flags().GetString("identity")
		if err != nil {
			return err
		}

		fmt.Println(identity)
		return nil
	},
}

func init() {
	{
		// initialize cli flags
		rootCmd.Flags().StringP("identity", "i", "ubuntu", "Docker Image identity defined in image_map.json")
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
