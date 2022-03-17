package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/vurokrazia/automations_go/configs"
	"github.com/vurokrazia/automations_go/entity"
	//"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	//"github.com/stianeikeland/go-rpio"
)

var ctx = context.Background()
var client *redis.Client

var Clients = make(map[int]entity.Client)

func ConnectNewClient(port string, channel_request chan<- entity.DevicePort) {
	fmt.Print(channel_request)
	pubsub := client.PSubscribe(ctx, "device_port_"+port)
	defer pubsub.Close()
	for {
		message, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("Dont posible read the message")
		}

		request := entity.DevicePort{}
		if err := json.Unmarshal([]byte(message.Payload), &request); err != nil {
			fmt.Println("Imposible read json")
		}
		InitLeds(request)
		channel_request <- request
	}
}

func SendMessageRedis() {
	err := client.Set(ctx, "key", time.Now(), 0).Err()
	if err != nil {
		panic(err)
	}
}

func GetRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     configs.GetRedisconfiguration().Addr,
		Password: configs.GetRedisconfiguration().Password,
		DB:       0,
	})
}

func main() {
	client = GetRedis()
	SendMessageRedis()
	channel_request := make(chan entity.DevicePort)
	go ConnectNewClient("1", channel_request)
	go ValidateChannel(channel_request)
	mux := mux.NewRouter()
	http.Handle("/", mux)
	fmt.Println("Server mountend 8000")
	http.ListenAndServe(":8000", nil)
}

func ValidateChannel(request chan entity.DevicePort) {
	for {
		select {
		case r := <-request:
			SendMessage(r)
		}
	}
}

func SendMessage(request entity.DevicePort) {
	for _, client := range Clients {
		if err := client.Websocket.WriteJSON(request); err != nil {
			return
		}
	}
}

func InitLeds(dp entity.DevicePort) {
	fmt.Println("opening gpio")
	fmt.Println(dp.Status)
	// err := rpio.Open()
	// if err != nil {
	// 	panic(fmt.Sprint("unable to open gpio", err.Error()))
	// }

	// defer rpio.Close()
	// i, _ := strconv.ParseInt(dp.Port, 10, 32)
	// pin := rpio.Pin(i)
	// pin.Output()
	// fmt.Println(dp.Status == "1")
	// if dp.Status == "1" {
	// 	pin.High()
	// } else {
	// 	pin.Low()
	// }
}
