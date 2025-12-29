package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/CarlosEduardoAD/broadcast-server/internal/helpers"
	"github.com/CarlosEduardoAD/broadcast-server/pkg/client"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var clientTtyCmd = &cobra.Command{
	Use:   "join",
	Short: "Run a simple TTY client to send and receive messages",
	Run: func(cmd *cobra.Command, args []string) {

		basicAuth := cmd.Flags().Lookup("basicAuth")

		if basicAuth == nil {
			log.Println("basicAuth flag is required")
			os.Exit(1)
		}

		basicAuthStr := basicAuth.Value.String()

		cl := client.NewClient(&helpers.Dialer{})
		headers := http.Header{"Origin": {"http://localhost:8080"}, "Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(basicAuthStr))}}
		if err := cl.Connect("ws://localhost:8080/connect", headers); err != nil {
			log.Fatal("erro ao conectar:", err)
		}
		defer cl.Close()

		done := make(chan struct{})
		var closeOnce sync.Once
		closeDone := func() {
			closeOnce.Do(func() {
				close(done)
				os.Exit(0)
			})
		}

		go func() {
			for {

				var msgDeserialized map[string]interface{}

				_, msg, err := cl.Conn.ReadMessage()
				if err != nil {

					isNet := err == io.EOF

					if !isNet {
						log.Println("erro ao ler mensagem:", err)
						closeDone()
						return
					}

					closeDone()
					return
				}

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
			closeDone()
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
