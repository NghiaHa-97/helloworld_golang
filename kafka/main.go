package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// the topic and broker address are initialized as constants
const (
	topic          = "message-log"
	broker1Address = "127.0.0.1:9092"
	broker2Address = "127.0.0.1:9093"
	broker3Address = "127.0.0.1:9094"
)

func main() {
	ctx := context.Background()
	go produce(ctx)
	consume(ctx)
}

func produce(ctx context.Context) {

	i := 0

	// init broker and topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},

		Topic: topic,
		//set defaunt partitions in file server.properties num.partitions=3
	})

	for {

		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),

			Value: []byte(" message" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write" + err.Error())
		}

		fmt.Println("writes:", i)
		i++
		//sleep 1s
		time.Sleep(time.Second)
	}
}

func consume(ctx context.Context) {

	//read topic
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {

		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		fmt.Println("received: ", string(msg.Value))
	}
}
