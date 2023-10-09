package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
	"wordofwisdom/internal/pow"
	"wordofwisdom/internal/protocol"
)

func Run(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	fmt.Println("connected to ", address)
	defer conn.Close()

	for {
		message, err := handleConnection(ctx, conn, conn)
		if err != nil {
			return err
		}
		fmt.Println("quote result:", message)
		time.Sleep(5 * time.Second)
	}
}

func handleConnection(ctx context.Context, readerConnection io.Reader, writerConnection io.Writer) (string, error) {
	reader := bufio.NewReader(readerConnection)

	// 1. requesting challenge
	err := sendMessage(protocol.Message{
		Header: protocol.RequestChallenge,
	}, writerConnection)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	// reading and parsing response
	msgStr, err := readConnMsg(reader)
	fmt.Println("Reading challenge:", msgStr)
	if err != nil {
		return "", fmt.Errorf("error reading msg: %w", err)
	}

	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("error parsing msg: %w", err)
	}

	var hashcash *pow.Hashcash
	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return "", fmt.Errorf("error parsing hashcash: %w", err)
	}
	fmt.Println("got hashcash: ", hashcash)

	hashcash, err = hashcash.ComputeHashcash(1000000)
	if err != nil {
		return "", fmt.Errorf("error computing hashcash: %w", err)
	}

	fmt.Println("hashcash computed:", hashcash)
	// marshal solution to json
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return "", fmt.Errorf("Error decoding json hashcash: %w", err)
	}

	// 3. send challenge solution back to server
	err = sendMessage(protocol.Message{
		Header:  protocol.RequestResource,
		Payload: string(byteData),
	}, writerConnection)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	fmt.Println("challenge sent to server")

	// 4. get result quote from server
	msgStr, err = readConnMsg(reader)
	if err != nil {
		return "", fmt.Errorf("error reading msg: %w", err)
	}
	msg, err = protocol.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("error parsing msg: %w", err)
	}
	return msg.Payload, nil
}

func readConnMsg(connectionReader *bufio.Reader) (string, error) {
	return connectionReader.ReadString('\n')
}

func sendMessage(msg protocol.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.ToString())
	_, err := conn.Write([]byte(msgStr))
	return err
}
