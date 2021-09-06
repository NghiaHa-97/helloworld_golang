# Kafka cluster and Debezium, postgers using golang, docker-compose


## Mô hình kết nối các container:
![](./images/diagram.PNG)

## 1: Tạo container postgres từ ``` debezium/postgres:latest```
### environment:
- `POSTGRES_PASSWORD` : pass loggin vào postgers
- `POSTGRES_USER`     : user name login
```md
postgres:
    image: debezium/postgres:latest
    container_name: postgres
    networks:
      - kafka-network
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
    ports:
      - 5432:5432
```
## 2: Tạo container Zookeeper từ `confluentinc/cp-zookeeper:latest`
### environment:
- `ZOOKEEPER_TICK_TIME`       : tickTime là Được sử dụng để điều chỉnh heart beat và timeout
- `ZOOKEEPER_CLIENT_PORT`     : port mở để các client kết nối tới
```md
zookeeper:
    image: confluentinc/cp-zookeeper:latest
#    image: debezium/zookeeper:latest
    container_name: zookeeper
    networks:
      - kafka-network
    ports:
      - 2181:2181
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
```

## 3: Tạo container kafka từ ``` confluentinc/cp-kafka:latest```
### environment:
[Tham Khảo: ADVERTISED_LISTENER](https://www.confluent.io/blog/kafka-listeners-explained/)
- `KAFKA_ZOOKEEPER_CONNECT`: kết nối tói container zookeeper
- `KAFKA_ADVERTISED_LISTENERS`: danh sách các port cung cấp để thao tác với kafka broker
- `KAFKA_LISTENER_SECURITY_PROTOCOL_MAP`: xác định các cặp khóa / giá trị cho giao thức bảo mật để sử dụng cho mỗi tên người nghe
- `KAFKA_INTER_BROKER_LISTENER_NAME`: Máy chủ / IP được sử dụng phải có thể truy cập được từ máy môi giới cho người khác
- `KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR`: Hệ số nhân rộng cho chủ đề bù trừ mặc định là 3, nếu việc tạo chủ đề tự động mà không đáp ứng dược số replication đc khai báo sẽ không tạo thành công
- `KAFKA_BROKER_ID`: id của broker
- `KAFKA_AUTO_CREATE_TOPICS_ENABLE`: tự động tạo topic (default là true)
- `KAFKA_NUM_PARTITIONS`: số partition của topic đc tạo tự động
- `KAFKA_DEFAULT_REPLICATION_FACTOR`: số nhân bản đc tạo ra (phải nhỏ hơn hoặc bằng số kafka broker : vd: 2 kafka broker không thể có 3 replication của partition đc) 
#### VD tạo 3 kafka broker
```md
kafka-1:
    image: confluentinc/cp-kafka:latest
#    image: debezium/kafka:latest
    container_name: kafka-1
    networks:
      - kafka-network
    depends_on:
      - zookeeper
    ports:
      - 9092:9092
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
#      KAFKA_LOG_CLEANER_DELETE_RETENTION_MS: 5000
      KAFKA_BROKER_ID: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
      KAFKA_NUM_PARTITIONS: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 3
  #      KAFKA_MIN_INSYNC_REPLICAS: 1

  kafka-2:
    image: confluentinc/cp-kafka:latest
    #    image: debezium/kafka:latest
    container_name: kafka-2
    networks:
      - kafka-network
    depends_on:
      - zookeeper
    ports:
      - 9093:9093
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:29093,PLAINTEXT_HOST://localhost:9093
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      #      KAFKA_LOG_CLEANER_DELETE_RETENTION_MS: 5000
      KAFKA_BROKER_ID: 2
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
      KAFKA_NUM_PARTITIONS: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 3
  #      KAFKA_MIN_INSYNC_REPLICAS: 1

  kafka-3:
    image: confluentinc/cp-kafka:latest
    #    image: debezium/kafka:latest
    container_name: kafka-3
    networks:
      - kafka-network
    depends_on:
      - zookeeper
    ports:
      - 9094:9094
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-3:29094,PLAINTEXT_HOST://localhost:9094
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      #      KAFKA_LOG_CLEANER_DELETE_RETENTION_MS: 5000
      KAFKA_BROKER_ID: 3
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: "true"
      KAFKA_NUM_PARTITIONS: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 3
  #      KAFKA_MIN_INSYNC_REPLICAS: 1
```
## 4: Tạo container connector từ ``` debezium/connect:latest``` liên kết kafka và postgres
### Debezium chịu trách nhiệm đọc dữ liệu từ hệ thống dữ liệu nguồn ( postgres là 1 ví dụ ) và đẩy nó vào một chủ đề kafka (được đặt tên tự động theo bảng) ở định dạng phù hợp.
### environment:
- `GROUP_ID`: xác định nhóm Connect cluster 
- `CONFIG_STORAGE_TOPIC`: Tên của chủ đề Kafka nơi lưu trữ cấu hình trình kết nối
- `OFFSET_STORAGE_TOPIC`: Tên của chủ đề Kafka nơi lưu trữ các offset của trình kết nối
- `STATUS_STORAGE_TOPIC`: Tên của chủ đề Kafka nơi trình kết nối và trạng thái tác vụ được lưu trữ
- `BOOTSTRAP_SERVERS`: danh sách các port để liên kết với kafka cluster `host1:port1,host2:port2,...` chỉ cần 1 kafka trong cụm sau đó nó sẽ tự dộng liên kết với tất cả các kafka broker còn lại

```md
connector:
    image: debezium/connect:latest
    container_name: kafka_connect_with_debezium
    networks:
      - kafka-network
    ports:
      - "8083:8083"
    environment:
      GROUP_ID: 1
      CONFIG_STORAGE_TOPIC: my_connect_configs
      OFFSET_STORAGE_TOPIC: my_connect_offsets
      STATUS_STORAGE_TOPIC: my_connect_statuses
      BOOTSTRAP_SERVERS: kafka-1:29092
    depends_on:
      - kafka-1
      - kafka-2
      - kafka-3
      - postgres
```

## 5: đăng ký connect giữa debezium và postgresql (name của topic được tự động đặt có cấu trúc là `<database.server.name>.<schema>.<table> ` )
### file debezium_connect.json 
- `database.*`: cấu hình là các tham số kết nối cho cơ sở dữ liệu postgres
- `database.server.name`:là tên tự gán cho cơ sở dữ liệu của mình
- `table.whitelist`: là trường thông báo cho trình kết nối debezium chỉ đọc các thay đổi dữ liệu từ bảng đó. Tương tự, cũng có thể đưa các bảng hoặc lược đồ vào danh sách trắng hoặc danh sách đen. Theo mặc định, debezium đọc từ tất cả các bảng trong một lược đồ.
- `connector.class`:là trình kết nối được sử dụng để kết nối với cơ sở dữ liệu postgres 
- `name`:tên chỉ định cho trình kết nối
- `key.converter`:  được sử dụng để chuyển đổi các khóa của trình kết nối sang dạng được lưu trữ trong Kafka. Mặc định là `org.apache.kafka.connect.json.JsonConverter`
- `value.converter.schemas.enable`:
- `key.converter.schemas.enable`: Khi các thuộc tính `key.converter.schemas.enable` và `value.converter.schemas.enable` được đặt thành `false` (mặc định), chỉ dữ liệu được truyền cùng mà không có schema. Điều này làm giảm chi phí tải trọng cho các ứng dụng không cần schema.
- `value.converter`: được sử dụng để chuyển đổi các giá trị của trình kết nối sang dạng được lưu trữ trong Kafka. Mặc định là `org.apache.kafka.connect.json.JsonConverter`

- `topic.creation.default.replication.factor`: số bản copy của partition đc tạo ra của topic đó (cấu hình ở đây thì không cần cấu hình ở kafka-broker)
- `topic.creation.default.partitions`: số partition của topic đc tạo ra
- 
```md
{
  "name": "product_connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "tasks.max": 1,
    "database.hostname": "postgres",
    "database.port": "5432",
    "database.user": "postgres",
    "database.password": "postgres",
    "database.dbname" : "postgres",
    "database.server.name": "postgres",

    "key.converter": "org.apache.kafka.connect.storage.StringConverter",
    "key.converter.schemas.enable": "false",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "value.converter.schemas.enable": "false",

    "topic.creation.default.replication.factor": 3,
    "topic.creation.default.partitions": 3
  }
}
```
### đăng ký connect bằng package net/http golang
#### kiểm tra xem đã đăng ký connect chưa
- kiểm tra xem có connnect đến postgres chưa bằng thuộc tính name trong file config `debezium_connect.json` `"name": "product_connector",` nếu `response.StatusCode != 200` thì mới đăng ký 
 ```md
 response, err := http.Get("http://localhost:8083/connectors/product_connector")
 ```
 #### đăng ký connect
- đọc file json `plan, e := ioutil.ReadFile("./debezium_connect.json")`
- send data đến api `http://localhost:8083/connectors/` kiểu body `application/json`
 ```md
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
 ```
 
 ## 6: Cấu hình Consumer và read mesage
 ### Cấu hình consumer 
 - `bootstrap.servers`: các port của kafka broker để lắng nghe sự kiện
 - `group.id`: id group của consumer
 - `auto.offset.reset`:vị trí offset để các group id khác bắt đầu lắng nghe sự kiện phát ra
     - `earliest`: đọc từ đầu đối với các group id khác nhau
     - `latest`: đọc cái mới nhất đc emit ra sau khi subscribe topic đó
 - `go.events.channel.enable`: `true` để nhận thêm sự kiện dùng `consumer.Events()`
```md
func CreateKafkaConsumer() *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094",
		"group.id":          "myGroup1",
		"auto.offset.reset": "earliest",
		"go.events.channel.enable": true,
	})

	if err != nil {
		panic(err)
	}
	return consumer
}
```
### Subscribe topic 
- Topic name được auto create theo cấu trúc `<database.server.name>.<schema>.<table>`
```md
func SubscribeTopic(consumer *kafka.Consumer)  {
	consumer.Subscribe("postgres.debezium.product",nil)
	fmt.Println("Subscribed to product topic")
}
```

### Read topic mesage 
- `consumer.ReadMessage(-1)`: đọc mesage mỗi khi kafka phát ra event nếu param khác `-1` thì sau mỗi khoảng thòi gian không có dữ liệu được emit thì consumer sẽ trả về message `nil`  , ReadMessage chỉ nhận duy nhất được mesage không nhận đc 1 số event khác
- hoặc có thể dùng `consumer.Events()` nếu  `go.events.channel.enable` được bật để nhận thêm các sự kiện khác
#### Đọc mesage và in ra dạng json
 ```md
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

```


# Demo:
## 1: run docker compose: 
```md
docker-compose -f docker-compose.yml up -d
```
## 2: Tạo database
- Truy cập vào container postgres (-it là input terminal để thực hiện input thao tác với container)
```md
docker exec -it postgres bash
```
- sử dụng psql để đăng nhập vào postgres `-h` là host, `-p` là port xuất ra ngoài đc cấu hình ở docker-compose.yml, `-U` user name
```md
psql -h localhost -p 5432 -U postgres
```
-- tạo schema và table và thực hiện insert
```md
CREATE SCHEMA debezium;
```
```md
CREATE TABLE debezium.product (id int);
```
```md
insert into debezium.product values (23);
```

## 3: đăng ký connect và nhận mesage bằng chạy main.go

```md
go run main.go
```
- Result của mesage bước trước insert vào table product
```md
Message on [ Topic: postgres.debezium.product , Partition: 0, Offset: 0 ] {
        "after": {
                "id": 23
        },
        "before": null,
        "op": "r",
        "source": {
                "connector": "postgresql",
                "db": "postgres",
                "lsn": 23700112,
                "name": "postgres",
                "schema": "debezium",
                "sequence": "[null,\"23700168\"]",
                "snapshot": "last",
                "table": "product",
                "ts_ms": 1630916915780,
                "txId": 555,
                "version": "1.6.1.Final",
                "xmin": null
        },
        "transaction": null,
        "ts_ms": 1630916915785
}
```

## 4: xem thông tin chi tiết của topic `postgres.debezium.produc`
- xem danh sách các topic
```md
kafka-topics --zookeeper localhost:2181  --list
```
```md
--kết quả
__consumer_offsets
my_connect_configs
my_connect_offsets
my_connect_statuses
postgres.debezium.product
```
- xem chi tiết của topic `postgres.debezium.product`
```md
kafka-topics --zookeeper localhost:2181  --topic postgres.debezium.product --describe
```
```md
--kết quả
Topic: postgres.debezium.product        TopicId: 3mbRZoQfRg-6kHLdyLGZZA PartitionCount: 3       ReplicationFactor: 3    Configs:
        Topic: postgres.debezium.product        Partition: 0    Leader: 1       Replicas: 1,3,2 Isr: 1,3,2
        Topic: postgres.debezium.product        Partition: 1    Leader: 2       Replicas: 2,1,3 Isr: 2,1,3
        Topic: postgres.debezium.product        Partition: 2    Leader: 3       Replicas: 3,2,1 Isr: 3,2,1
```
- Giải thích và thử xóa 1 broker xem replication hoạt động
    - partition 0 vị trí chính là ở brokerId = 1, và được nhận lên làm 3 xuất hiện ở cả 3 brokerId 1,3,2 (Replicas: 1,3,2), (Isr: 1,3,2) là id của broker được đồng bộ
    - `Message on [ Topic: postgres.debezium.product , Partition: 0, Offset: 0 ]` kết quả của đọc mesage đang ở partition 0 vị trí chính (leader) là broker có Id =1. do đó thử stop broker  đó rồi đọc lại mesage xem có bị mất dự liệu không 
    - `docker ps` để xem tất cả các container đang chạy và lấy ID của container có Borker id = 1
    - `docker stop <id>` để dừng container 
    - kiểm tra lại topic  `kafka-topics --zookeeper localhost:2181  --topic postgres.debezium.product --describe`
```md
--kết quả
Topic: postgres.debezium.product        TopicId: 3mbRZoQfRg-6kHLdyLGZZA PartitionCount: 3       ReplicationFactor: 3    Configs:
        Topic: postgres.debezium.product        Partition: 0    Leader: 3       Replicas: 1,3,2 Isr: 3,2
        Topic: postgres.debezium.product        Partition: 1    Leader: 2       Replicas: 2,1,3 Isr: 2,3
        Topic: postgres.debezium.product        Partition: 2    Leader: 3       Replicas: 3,2,1 Isr: 3,2
```
-
    - theo kết quả trên thì `Isr` broker được đồng bộ đã không còn Id 1. nhưng partition 0 vẫn còn và đc đổi thánh vị trí của Broker ID = 3
    - thử kiểm tra xem còn nhận được mesage không bằng cách thay đổi tên group consumer và chạy lại main.go
    - kết quả vẫn nhận đc
```md
Message on [ Topic: postgres.debezium.product , Partition: 0, Offset: 0 ] {
    "after": {
            "id": 23
    },
    "before": null,
    "op": "r",
    "source": {
            "connector": "postgresql",
            "db": "postgres",
            "lsn": 23700112,
            "name": "postgres",
            "schema": "debezium",
            "sequence": "[null,\"23700168\"]",
            "snapshot": "last",
            "table": "product",
            "ts_ms": 1630916915780,
            "txId": 555,
            "version": "1.6.1.Final",
            "xmin": null
    },
    "transaction": null,
    "ts_ms": 1630916915785
}
```
* [Tham Khảo: Change Data Capture Using Debezium Kafka and Pg](https://www.startdataengineering.com/post/change-data-capture-using-debezium-kafka-and-pg/)
* [Tham Khảo: kafka.apache.org](https://kafka.apache.org/documentation/#configuration)
* [Tham Khảo: ADVERTISED_LISTENER](https://www.confluent.io/blog/kafka-listeners-explained/)
