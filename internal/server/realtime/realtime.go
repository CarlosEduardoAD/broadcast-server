package realtime

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/CarlosEduardoAD/broadcast-server/internal/pub"
	"github.com/CarlosEduardoAD/broadcast-server/internal/request/message"
	"github.com/CarlosEduardoAD/broadcast-server/internal/response/error_response"
	"github.com/CarlosEduardoAD/broadcast-server/internal/sub"
	"github.com/CarlosEduardoAD/broadcast-server/internal/utils"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: utils.Authorize,
	}

	publisher = pub.NewPublisher([]sub.Subscriber{})
)

func WebsocketRoute(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	host, _, _ := net.SplitHostPort(c.RemoteAddr().String())

	sub := sub.Subscriber{
		Name: "RANDOM",
		Ip:   host,
		Conn: c,
	}

	close := c.CloseHandler()
	c.SetCloseHandler(func(code int, text string) error {
		publisher.Remove(sub)

		err := close(code, text)

		return err
	})

	publisher.Subscribe(sub)

	defer c.Close()
	for {
		var unmarshelled message.Message

		mt, message, errWs := c.ReadMessage()

		errJson := json.Unmarshal(message, &unmarshelled)

		if errJson != nil {
			result, err := json.Marshal(error_response.NewErrorResponse(errJson))

			if err != nil {
				panic(err)
			}

			c.WriteMessage(1, result)
			break
		}

		if errWs != nil {
			result, err := json.Marshal(error_response.NewErrorResponse(errWs))

			if err != nil {
				panic(err)
			}

			c.WriteMessage(1, result)

			if strings.Contains(errWs.Error(), "abnormal closure") {
				c.Close()

				publisher.Remove(sub)
				break
			} else {
				break
			}
		}

		publisher.Publish(message, mt)
	}
}
