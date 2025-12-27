/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/CarlosEduardoAD/broadcast-server/internal/server"
	"github.com/CarlosEduardoAD/broadcast-server/internal/utils"
	"github.com/spf13/cobra"
)

var (
	port = 8080
	srv  = server.NewServer(port)
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server [subcommand]",
	Short: "Commands to manage the websocket server",
	Long:  `The server command allows the user to start and stop it.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "start":
			if os.Getenv("IS_DAEMON") != "1" {
				pid, err := utils.Fork()
				if err != nil {
					fmt.Printf("Unable to fork process")
					os.Exit(1)
				}
				fmt.Printf("Start child - pid = %d\n", pid)
				os.Exit(0) // parent must exit after spawning child
			}

			srv.Connect()
		case "stop":
			fmt.Println("Stopping server...")
			srv.Disconnect()
		default:
			fmt.Println("Unknown subcommand. Use 'start' or 'stop'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
