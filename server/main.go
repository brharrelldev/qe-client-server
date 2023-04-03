package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brharrelldev/qe-client-server/pkg/types"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

type Method string

type Response struct {
	Message string
	Results []types.Payload
}

func responseHandler(conn net.Conn, errChan chan error, dataChan chan []byte) {

	buf := make([]byte, 1024)
	offset, err := conn.Read(buf)
	if err != nil {
		errChan <- err
	}

	data := buf[:offset]

	dataChan <- data

}

func actionDecider(db *DB, data []byte, resultChan chan *Response, errChan chan error) {

	var payload *types.Payload

	fmt.Println("incoming data", string(data))

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&payload); err != nil {

		errChan <- fmt.Errorf("error decoding incoming json %v", err)
		return
	}

	switch payload.Method {
	case "create":

		dataBytes, err := json.Marshal(payload.Data)
		if err != nil {
			errChan <- err
			return
		}

		bytesReader := bufio.NewReader(bytes.NewReader(dataBytes))
		result, err := db.Create(bytesReader)
		if err != nil {
			errChan <- fmt.Errorf("error creating new entry in the database")
			return

		}

		resultChan <- result
		return
	case "list":

		resp, err := db.List()
		if err != nil {
			errChan <- err
			return
		}

		resultChan <- resp
		return

	case "update":
		var payload types.Payload

		if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&payload); err != nil {
			errChan <- err
			return
		}

	case "delete":

		var payload types.Payload

		if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&payload); err != nil {
			errChan <- err
			return
		}

		if err := db.Delete(payload.Id); err != nil {
			errChan <- err
			return
		}

	default:
		fmt.Println("unrecognized method")
		errChan <- errors.New("unrecognized method")
		return
	}

}

func StartServer(qeLogger *zap.Logger, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error creating new listener %v", err)
	}

	db, err := NewDB()
	if err != nil {
		return fmt.Errorf("error creating new database %v", err)
	}

	for {

		errChan := make(chan error)
		dataChan := make(chan []byte, 1024)
		conn, err := lis.Accept()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error accepting connection %v", err)
		}

		go responseHandler(conn, errChan, dataChan)

		select {
		case data := <-dataChan:

			resultChan := make(chan *Response)
			go actionDecider(db, data, resultChan, errChan)

			select {
			case result := <-resultChan:
				fmt.Println("results", result)
				output, err := json.Marshal(result)
				if err != nil {
					continue
				}

				if _, err := conn.Write(output); err != nil {
					fmt.Println(err)
					continue
				}

				if err := conn.Close(); err != nil {
					fmt.Println(err)
					continue
				}
			case err := <-errChan:
				fmt.Println("err")
				if _, err := conn.Write([]byte(err.Error())); err != nil {
					continue
				}
				if err := conn.Close(); err != nil {
					continue
				}
			}
		case err = <-errChan:
			fmt.Println(err)

		}

	}

	return nil

}

func main() {

	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)

	port := ":3000"

	qelogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := StartServer(qelogger, port); err != nil {
			errChan <- err
		}
	}()

	signal.Notify(sigChan, os.Interrupt)

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-sigChan:
		log.Println("program terminated")
	}

}
