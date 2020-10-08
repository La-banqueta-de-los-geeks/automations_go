package main
import(
	"fmt"
	"github.com/go-redis/redis/v8"
	"context"
	"time"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

var ctx = context.Background()
var client = redis.NewClient(&redis.Options{
	Addr: "redis_go_rails:6379",
	Password:"vurokrazia",
	DB: 0,
})

type Message struct {
	Id int
	Title string
	Body string
}

type DevicePort struct {
	Id int
	Status string
	Port string
}

type Client struct {
	Id int
	websocket *websocket.Conn
}

var Clients = make(map[int]Client)

func ConnectNewClient(channel_request chan DevicePort)  {
	pubsub := client.PSubscribe(ctx, "device_1*")
	defer pubsub.Close()
	for {
		message, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("Dont posible read the message")
		}
		
		request := DevicePort{}
		if err := json.Unmarshal([]byte(message.Payload), &request); err != nil {
			fmt.Println("Imposible read json")
		}
		fmt.Println(request.Id)
		fmt.Println(request.Status)
		fmt.Println(request.Port)
		channel_request <- request
		// fmt.Println(message)
		// fmt.Println(message.Payload)
	}
}

func SendMessageRedis()  {
	err := client.Set(ctx, "key", time.Now(), 0).Err()
	if err != nil {
			panic(err)
	}
}

func main(){
	fmt.Println("Hello World")
	SendMessageRedis()
	channel_request := make(chan DevicePort)
	go ConnectNewClient(channel_request)
	go ValidateChannel(channel_request)
	mux := mux.NewRouter()
	mux.HandleFunc("/Subscribe/", Subscribe).Methods("GET")
	http.Handle("/", mux)
	fmt.Println("Server mountend 8000")
	http.ListenAndServe(":8000",nil)
}

func Subscribe(w http.ResponseWriter, r *http.Request){
	
	ws, err := websocket.Upgrade(w,r,nil,1024,1024)

	if err != nil {
		return
	}

	fmt.Println("New Web Socket")

	count := len(Clients)
	new_client := Client{ count, ws}
	Clients[count] = new_client
	fmt.Printf("hello %v\n", count )
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(Clients, new_client.Id)
			fmt.Printf("see you later %v\n", count )
		}
	}
}

func ValidateChannel(request chan DevicePort ) {
	for {
		select {
		case r := <- request:
			SendMessage(r)
		}
	}
}

func SendMessage(request DevicePort)  {
	for _, client := range Clients {
		if err := client.websocket.WriteJSON(request); err != nil {
			return 
		}
	}
}
