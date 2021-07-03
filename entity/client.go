package entity
import(
	"github.com/gorilla/websocket"
)
type Client struct {
	Id int
	Websocket *websocket.Conn
}
