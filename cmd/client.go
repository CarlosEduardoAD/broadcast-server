/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client [subcommand]",
	Short: "All the commands a client can use to interact with the chat",
	Long:  `A client can use it's terminal as TTY to send chat messages to random people!`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
