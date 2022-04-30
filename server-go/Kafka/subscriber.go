package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	brokerList        = kingpin.Flag("brokerList", "List of brokers to connect").Default("127.0.0.1:9092").Strings()
	topic             = kingpin.Flag("topic", "Topic name").Default("important").String()
	partition         = kingpin.Flag("partition", "Partition number").Default("0").String()
	offsetType        = kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

type Data struct {
	GameID   int32  `bson:"game_id"`
	Players  int32  `bson:"players"`
	GameName string `bson:"game_name"`
	Winner   int32  `bson:"winner"`
	Queue    string `bson:"queue"`
}

func main() {
	kingpin.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				*messageCountStart++
				go insertar_mongo(msg.Value)
				log.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", *messageCountStart, "messages")
}


func insertar_mongo(contenido []byte) {
	// MongoDB
	var djson Data
	err := json.Unmarshal(contenido, &djson)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongoadmin:Arqui2022_2022@35.223.112.104:27017/?compressors=disabled&gssapiServiceName=mongodb"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("proyectof2").Collection("JUEGO")

	if err != nil {
		log.Println(err)
	}
	currentTime := time.Now()
	t := currentTime.Add(-time.Hour * 6)
	doc := bson.D{{"game_id", djson.GameID}, {"players", djson.Players}, {"game_name", djson.GameName}, {"winner", djson.Winner}, {"queue", djson.Queue}, {"fecha_hora", t.Format("2006-01-02 15:04:05")}}

	result, insertErr := collection.InsertOne(ctx, doc)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	fmt.Println(result)

	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var posts []Data
	if err = cur.All(ctx, &posts); err != nil {
		panic(err)
	}
	fmt.Println(posts)

}
//docker rmi wurstmeister/zookeeper wurstmeister/kafka server-go_subscriber server-go_producer
//docker rmi edson2021/subscribergo edson2021/server_grpc_201701029 edson2021/grpclientnode
