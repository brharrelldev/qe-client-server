package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

func responseHandler(conn net.Conn, errChan chan error, dataChan chan []byte) {

	buf := make([]byte, 1024)
	offset, err := conn.Read(buf)
	if err != nil {
		errChan <- err
	}

	if err := conn.Close(); err != nil {
		errChan <- err
	}

	data := buf[:offset]

	dataChan <- data

	fmt.Println(dataChan)
}

func StartServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error creating new listener %v", err)
	}

	for {

		errChan := make(chan error)
		dataChan := make(chan []byte, 1024)
		conn, err := lis.Accept()
		if err == io.EOF {
			break
		}

		fmt.Println(conn)

		if err != nil {
			return fmt.Errorf("error accepting connection %v", err)
		}

		go responseHandler(conn, errChan, dataChan)

		select {
		case data := <-dataChan:
			fmt.Println("incoming data", string(data))
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

	go func() {
		if err := StartServer(port); err != nil {
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
