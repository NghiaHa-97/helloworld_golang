package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// the topic and broker address are initialized as constants

//run zookeeper:zookeeper-server-start.bat config\zookeeper.properties
//run 3 broker:kafka-server-start.bat config\server.properties
//run 3 broker:kafka-server-start.bat config\server1.properties
//run 3 broker:kafka-server-start.bat config\server2.properties

const (
	topic          = "replication_topic"
	broker2Address = "127.0.0.1:9092"
	broker1Address = "127.0.0.1:9093"
	broker3Address = "127.0.0.1:9094"
)

func main() {
	ctx := context.Background()

	// conn, _ := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, 0)

	go produce(ctx)
	// time.Sleep(10 * time.Second)
	consume(ctx)
}

func produce(ctx context.Context) {
	logger := log.New(os.Stdout, "kafka Writer: ", 0)
	i := 0

	// send to broker
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker1Address, broker2Address, broker3Address},
		RequiredAcks: 0,
		Topic:        topic,
		Logger:       logger,
		//set defaunt partitions in file server.properties num.partitions=3
	})

	for {

		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),

			Value: []byte(" message" + strconv.Itoa(i)),
		})
		if err != nil {
			// panic("could not write" + err.Error())
			fmt.Println("could not write" + err.Error())
		}

		fmt.Println("writes:", i)
		i++
		//sleep 1s
		time.Sleep(1 * time.Second)

	}
}

func consume(ctx context.Context) {
	logger := log.New(os.Stdout, "kafka Reader: ", 0)
	//read topic
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker1Address, broker2Address, broker3Address},
		Topic:    topic,
		GroupID:  "my-group",
		MinBytes: 5,
		// the kafka library requires you to set the MaxBytes
		// in case the MinBytes are set
		MaxBytes: 1e6,
		Logger:   logger,
	})
	// r.SetOffset(13)
	for {

		msg, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Println("could not read message " + err.Error())
		}

		fmt.Println("received: ", string(msg.Value))
	}
}
