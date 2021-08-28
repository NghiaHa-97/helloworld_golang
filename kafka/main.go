package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

// the topic and broker address are initialized as constants

//run zookeeper:zookeeper-server-start.bat config\zookeeper.properties
//run 3 broker:kafka-server-start.bat config\server.properties
//run 3 broker:kafka-server-start.bat config\server1.properties
//run 3 broker:kafka-server-start.bat config\server2.properties

const (
	//topic          = "replication_topic"
	broker1Address = "127.0.0.1:9092"
	broker2Address = "127.0.0.1:9093"
	broker3Address = "127.0.0.1:9094"
)

func main() {
	//ctx := context.Background()

	// conn, _ := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, 0)

	//go produce(ctx)
	// time.Sleep(10 * time.Second)
	//consume(ctx)

	log.Println("LLL")
	go produce()
	consume()
}

func produce() {
	//logger := log.New(os.Stdout, "kafka Writer: ", 0)
	//i := 0
	//
	////send to broker
	//w := kafka.NewWriter(kafka.WriterConfig{
	//	Brokers:      []string{broker1Address, broker2Address, broker3Address},
	//	RequiredAcks: 0,
	//	Topic:        topic,
	//	Logger:       logger,
	//	//set defaunt partitions in file server.properties num.partitions=3
	//})
	//
	//for {
	//
	//	err := w.WriteMessages(ctx, kafka.Message{
	//		Key: []byte(strconv.Itoa(i)),
	//
	//		Value: []byte(" message" + strconv.Itoa(i)),
	//	})
	//	if err != nil {
	//		// panic("could not write" + err.Error())
	//		fmt.Println("could not write" + err.Error())
	//	}
	//
	//	fmt.Println("writes:", i)
	//	i++
	//	//sleep 1s
	//	time.Sleep(1 * time.Second)
	//
	//}

	fmt.Println(fmt.Sprintf("%s,%s,%s",broker1Address,broker2Address,broker3Address))
	p, _ := kafka.NewProducer(&kafka.ConfigMap{
		"metadata.broker.list": fmt.Sprintf("%s,%s,%s", broker1Address, broker2Address, broker3Address),
		"acks":                 -1,
	})
	defer p.Close()

	// Delivery report cuar message to topic (acks)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	for _, word := range []string{"1", "2", "3", "4", "5", "6", "7"} {
		error := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
		if error!=nil{
			log.Println("ERROR send message: ", error)

		}
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func consume() {
	//fmt.Println(fmt.Sprintf("%s,%s,%s",broker1Address,broker2Address,broker3Address))
	//logger := log.New(os.Stdout, "kafka Reader: ", 0)
	////read topic
	//r := kafka.NewReader(kafka.ReaderConfig{
	//	// Brokers:  []string{broker1Address, broker2Address, broker3Address},
	//	Brokers:  []string{broker2Address},
	//	Topic:    topic,
	//	GroupID:  "my-group",
	//	MinBytes: 5,
	//	// the kafka library requires you to set the MaxBytes
	//	// in case the MinBytes are set
	//	MaxBytes: 1e6,
	//	Logger:   logger,
	//})
	//// r.SetOffset(13)
	//for {
	//
	//	msg, err := r.ReadMessage(ctx)
	//	if err != nil {
	//		fmt.Println("could not read message " + err.Error())
	//	}
	//
	//	fmt.Println("received: ", string(msg.Value))
	//}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        broker2Address,
		"group.id":                 "myGroup",
		"auto.offset.reset":        "earliest",
		"go.events.channel.enable": true,
		"enable.auto.commit":       false,
	})

	if err != nil {
		panic(err)
	}

	//c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)
	c.SubscribeTopics([]string{"myTopic"}, nil)

	for {
		//TODO: ReadMessage
		/*

		msg, err := c.ReadMessage(5*time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		*/


		ev := c.Events()

		fmt.Println("Message on ", <-ev)

	}

	c.Close()
}

//TODO:
// -------------
// dùng ReadMessage chỉ thăm dò đc 1 event duy nhất, Evwnt đc thăm dò nhiều  event cụ thể là for maxEvents = 1000
// if channel == nil {
//		maxEvents = 1
//	}
// -------------
//


