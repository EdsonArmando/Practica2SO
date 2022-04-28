package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"

	pb "github.com/EdsonArmando/demo-gRCP/proto"
	"github.com/Shopify/sarama"
	"google.golang.org/grpc"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Struct Juego
type Juego struct {
	GameID   int32  `bson:"game_id"`
	Players  int32  `bson:"players"`
	GameName string `bson:"game_name"`
	Winner   int32  `bson:"winner"`
	Queue    string `bson:"queue"`
}

var (
	port = flag.Int("port", 50051, "The server port")
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("127.0.0.1:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("important").String()
	maxRetry   = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedIniciarJuegoServer
}

//Retorna el struct
func GetStructGanador(GameID int32, Players int32, GameName string, Winner int32) Juego {

	return Juego{GameID, Players, GameName, Winner, "Kafka"}
}

// Metodo donde se ejecuta el juego
func (s *server) EjecutarJuego(ctx context.Context, in *pb.JuegoRequest) (*pb.Reply, error) {
	intGameId := in.GetGameId()
	players := in.GetPlayers()
	juego := ""
	var winner int32
	//Ejecutar el juego
	if intGameId == 1 {
		winner = Juego1(players)
		juego = "Adivinar Numero"
	} else if intGameId == 2 {
		winner = Juego2(players)
		juego = "Dados Mayor"
	} else if intGameId == 3 {
		winner = Juego3(players)
		juego = "Dados Iguales"
	} else if intGameId == 4 {
		winner = Juego4(players)
		juego = "Numero Menor"
	} else if intGameId == 5 {
		winner = Juego5(players)
		juego = "Cerillo Quemado"
	}

	//Struct Ganador
	ganadorStruct := GetStructGanador(intGameId, players, juego, winner)
	mongoStructAsBytes, _ := json.Marshal(ganadorStruct)
	buffer := bytes.NewBuffer(mongoStructAsBytes)

	//Enviar datos a Kafka
	kingpin.Parse()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = *maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(*brokerList, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic: *topic,
		Value: sarama.StringEncoder(mongoStructAsBytes),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *topic, partition, offset)

	//Return
	return &pb.Reply{Message: buffer.String()}, nil
}

//Adivinar Numero
func Juego1(players int32) int32 {
	var winner int32
	winner = players - 1
	//Numero
	num := rand.Intn(20-1) + 1
	//Ejecutar el algoritmo
	for i := 1; i <= int(players); i++ {
		if i == num {
			winner = int32(i)
			break

		}
	}
	//Retonar ganador
	return winner
}

//Dados Mayor
func Juego2(players int32) int32 {
	var winner int32
	winner = players - 1
	//Numero
	num := rand.Intn(20-1) + 1
	num2 := rand.Intn(20-1) + 1
	//Ejecutar el algoritmo
	for i := 1; i <= int(players); i++ {
		if i == (num + num2) {
			winner = int32(i)
			break
		}
	}
	//Retonar ganador
	return winner
}

//Dados Iguales
func Juego3(players int32) int32 {
	var winner int32
	winner = players - 1
	//Numero
	num := rand.Intn(20-1) + 1
	num2 := rand.Intn(20-1) + 1
	//Ejecutar el algoritmo
	for i := 1; i <= int(players); i++ {
		if num == num2 {
			winner = int32(i)
			break
		}
	}
	//Retonar ganador
	return winner
}

//Numero Menor
func Juego4(players int32) int32 {
	var winner int32
	winner = players - 1
	//Numero
	num := rand.Intn(20-1) + 1
	num2 := rand.Intn(20-1) + 1
	//Ejecutar el algoritmo
	for i := 1; i <= int(players); i++ {
		if num < num2 {
			winner = int32(i)
			break
		}
	}
	//Retonar ganador
	return winner
}

// CerilloQuemado
func Juego5(players int32) int32 {
	var winner int32
	winner = players - 1
	//Numero
	num := rand.Intn(20-1) + 1
	num2 := rand.Intn(20-1) + 1
	//Ejecutar el algoritmo
	for i := 1; i <= int(players); i++ {
		if i == (num + num2) {
			winner = int32(i)
			break
		}
	}
	//Retonar ganador
	return winner
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterIniciarJuegoServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
