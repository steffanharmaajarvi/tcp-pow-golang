package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"wordofwisdom/internal/pow"
	"wordofwisdom/internal/protocol"
	"wordofwisdom/internal/storage"
)

type ServerConfig struct {
	Port int
	Host string
}

type Storage interface {
	Add(string, int64) error
	Get(string) (bool, error)
	Delete(string)
}

func getStorage() Storage {
	return storage.InitInMemoryStorage(time.Now())
}

func Run(ctx context.Context, config ServerConfig) error {
	fmt.Println(ctx.Value("storage"))
	serverAddress := fmt.Sprintf("%s:%d", config.Host, config.Port)

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Println("Running server at ", listener.Addr())

	for {
		// Listen for an incoming connection.
		connection, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}
		// Handle connections in a new goroutine.
		go handleConnection(ctx, connection)
	}
}

func handleConnection(ctx context.Context, connection net.Conn) {
	fmt.Println("Handling connection from ", connection.RemoteAddr())

	defer connection.Close()

	reader := bufio.NewReader(connection)

	for {
		request, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading connection:", err)
			return
		}

		msg, err := processRequest(ctx, request, connection.RemoteAddr().String())
		if err != nil {
			fmt.Println("Error processing request:", err)
			return
		}
		if msg != nil {
			err := sendMessage(*msg, connection)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	}
}

func processRequest(ctx context.Context, messageStr string, clientInfo string) (*protocol.Message, error) {
	msg, err := protocol.ParseMessage(messageStr)
	if err != nil {
		return nil, err
	}

	switch msg.Header {
	case protocol.Close:
		return nil, fmt.Errorf("close connection")
	case protocol.RequestChallenge:
		return processRequestChallenge(ctx, clientInfo)
	case protocol.RequestResource:
		return processRequestResource(ctx, msg, clientInfo)
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

func processRequestChallenge(ctx context.Context, clientInfo string) (*protocol.Message, error) {
	fmt.Printf("client %s requests challenge\n", clientInfo)

	storageInst := ctx.Value("storage").(Storage)
	storageExpiration := ctx.Value("storageExpiration").(int64)

	baseValue := strconv.Itoa(rand.Int())

	err := storageInst.Add(baseValue, storageExpiration)
	if err != nil {
		return nil, fmt.Errorf("Error adding to storage: %v", err)
	}

	hashcash := pow.NewStd(baseValue, clientInfo)
	hashcashJson, err := toJson(hashcash)
	if err != nil {
		return nil, fmt.Errorf("Error serializing json: %v", err)
	}

	msg := protocol.Message{
		Header:  protocol.ResponseChallenge,
		Payload: hashcashJson,
	}
	fmt.Println("Sending challenge:", msg)

	return &msg, nil
}

func processRequestResource(ctx context.Context, message *protocol.Message, clientInfo string) (*protocol.Message, error) {
	fmt.Printf("client %s requests resource with payload %s\n", clientInfo, message.Payload)
	// parse client's solution
	var hashcash pow.Hashcash
	err := json.Unmarshal([]byte(message.Payload), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err unserialized hashcash: %w", err)
	}

	// validate hashcash params
	if hashcash.Client != clientInfo {
		return nil, fmt.Errorf("hashcash client mismatch")
	}

	storageInst := ctx.Value("storage").(Storage)

	exists, err := storageInst.Get(hashcash.BaseValue)
	if err != nil {
		return nil, fmt.Errorf("err get rand from storageInst: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("challenge expired or not sent")
	}

	//to prevent indefinite computing on server if client sent hashcash with 0 counter
	maxIter := hashcash.Counter
	if maxIter == 0 {
		maxIter = 1
	}
	_, err = hashcash.ComputeHashcash(maxIter)
	if err != nil {
		return nil, fmt.Errorf("invalid hashcash")
	}

	//get random quote
	fmt.Printf("client %s succesfully computed hashcash %s\n", clientInfo, message.Payload)

	quote := storage.GetRandomQuote()

	msg := protocol.Message{
		Header:  protocol.ResponseResource,
		Payload: quote,
	}

	storageInst.Delete(hashcash.BaseValue)
	return &msg, nil
}

// sendMsg - send protocol message to connection
func sendMessage(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToString())
	_, err := conn.Write([]byte(msgStr))
	return err
}

func toJson(object any) (string, error) {
	str, err := json.Marshal(object)
	if err != nil {
		return "", fmt.Errorf("Error converting to json")
	}

	strTrimmed := strings.ReplaceAll(string(str), "\n", "")

	return strTrimmed, nil
}
