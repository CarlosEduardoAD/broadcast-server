package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/CarlosEduardoAD/broadcast-server/internal/helpers"
	"github.com/CarlosEduardoAD/broadcast-server/pkg/client"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// clientTtyCmd runs a simple TTY-style client that can send and receive messages
var clientTtyCmd = &cobra.Command{
	Use:   "tty",
	Short: "Run a simple TTY client to send and receive messages",
	Run: func(cmd *cobra.Command, args []string) {
		cl := client.NewClient(&helpers.Dialer{})

		if err := cl.Connect("ws://localhost:8080/connect", nil); err != nil {
			log.Fatal("erro ao conectar:", err)
		}
		defer cl.Close()

		done := make(chan struct{})

		go func() {
			for {

				var msgDeserialized map[string]interface{}

				_, msg, err := cl.Conn.ReadMessage()
				if err != nil {
					log.Println("erro ao ler mensagem:", err)
					close(done)
					return
				}
				fmt.Println()

				err = json.Unmarshal(msg, &msgDeserialized)

				if err != nil {
					log.Println("erro ao deserializar mensagem:", err)
				}

				fmt.Println("<-", msgDeserialized["message"])
				fmt.Print("> ")
			}
		}()

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigc
			cl.Close()
			close(done)
		}()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		for {
			select {
			case <-done:
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Println("erro ao ler stdin:", err)
					return
				}
				line = strings.TrimSpace(line)
				if line == "" {
					fmt.Print("> ")
					continue
				}

				if err := cl.SendMessage(websocket.TextMessage, []byte("{\"message\" : \""+line+"\"}")); err != nil {
					log.Println("erro ao enviar mensagem:", err)
				}
				fmt.Print("> ")
			}
		}
	},
}

func init() {
	clientCmd.AddCommand(clientTtyCmd)
}
