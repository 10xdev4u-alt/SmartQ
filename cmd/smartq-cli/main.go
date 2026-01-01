package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "smartq-cli",
	Short: "A CLI for administering the SmartQ system",
	Long:  `A command-line interface for managing queues, users, and other settings in the SmartQ system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no command is provided
		fmt.Println("Welcome to the SmartQ CLI! Use 'smartq-cli help' for a list of commands.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
