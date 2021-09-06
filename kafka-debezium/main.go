package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"io/ioutil"
	"net/http"
)

func CreateKafkaConsumer() *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094",
		"group.id":          "myGroup11",
		"auto.offset.reset": "earliest",
		"go.events.channel.enable": true,
	})

	if err != nil {
		panic(err)
	}
	return consumer
}

func SubscribeTopic(consumer *kafka.Consumer)  {
	consumer.Subscribe("postgres.debezium.product",nil)
	fmt.Println("Subscribed to product topic")
}

func ReadTopicMessages(consumer *kafka.Consumer) string {

	var message string
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var raw map[string]interface{}
			if err1 := json.Unmarshal(msg.Value, &raw); err1 != nil {
				panic(err1)
			}
			jsonFormat,_ := json.MarshalIndent(raw,"","\t")
			fmt.Printf("Message on [ Topic: %s , Partition: %d, Offset: %s ] ", *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
			fmt.Println(string(jsonFormat))
			//message = message + string(msg.Value)
		} else {

			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	CloseConsumer(consumer)

	return message
}

func CloseConsumer(consumer *kafka.Consumer){
	consumer.Close()
}

func RegisterConnector() *http.Response {
	plan, e := ioutil.ReadFile("./debezium_connect.json")
	if e==nil{
		fmt.Println("read ok",plan)
	}
	response, err := http.Post("http://localhost:8083/connectors/","application/json",bytes.NewBuffer(plan))

	if err != nil{
		panic(err)
	}else {
		fmt.Println("success")
	}
	return response
}

func CheckConnector() {
	response, err := http.Get("http://localhost:8083/connectors/product_connector")

	fmt.Println(response.StatusCode)
	defer response.Body.Close()

	if err != nil{
		panic(err)
	}
	if response.StatusCode != 200 {
		fmt.Println("register connect" )
		RegisterConnector()
	}

	//body, _ := ioutil.ReadAll(response.Body)

	//fmt.Println(string(body))
}

func main()  {
	consumer := CreateKafkaConsumer()
	CheckConnector()

	SubscribeTopic(consumer)
	ReadTopicMessages(consumer)

	fmt.Scan()
}
