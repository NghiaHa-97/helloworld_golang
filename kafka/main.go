package main

import (
	"context"
	"fmt"
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

	i := 0

	// init broker and topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker1Address, broker2Address, broker3Address},
		RequiredAcks: -1,
		Topic:        topic,
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
		time.Sleep(3 * time.Second)

	}
}

func consume(ctx context.Context) {

	//read topic
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker1Address, broker2Address, broker3Address},
		Topic:       topic,
		GroupID:     "my-group",
		StartOffset: kafka.FirstOffset,
	})
	r.SetOffset(13)
	for {

		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		fmt.Println("received: ", string(msg.Value))
	}
}
