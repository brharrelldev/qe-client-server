package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brharrelldev/qe-client-server/pkg/types"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "qe-client"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "action",
			Aliases: []string{"a"},
		},
		&cli.StringFlag{
			Name: "path",
		},
	}

	app.Action = sendRequest

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func sendRequest(c *cli.Context) error {

	buf := make([]byte, 1024)

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("error dialing server %v", err)
	}

	var inventoryData *types.Data
	switch c.String("action") {
	case "get":

		if c.String("path") == "" {
			return errors.New("path is required for get request")
		}

		f, err := os.Open(c.String("path"))
		if err != nil {
			return fmt.Errorf("error opening up file %v", err)
		}

		defer f.Close()

		offset, err := f.Read(buf)
		if err != nil {
			return fmt.Errorf("error reading file %v", err)
		}

		data := buf[:offset]

		if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&inventoryData); err != nil {
			return fmt.Errorf("error decoding input data %v", err)
		}

		if inventoryData == nil {
			return errors.New("input not serialized correctly")
		}

		reqId := uuid.NewV4()
		req := types.Payload{
			RequestID: reqId.String(),
			Method:    "get",
			Data:      inventoryData,
		}

		msg, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("error marshaling request json %v", err)
		}

		if _, err := conn.Write(msg); err != nil {
			return fmt.Errorf("error sending data to server %v", err)
		}

	case "create":
		return errors.New("not implemented")

	case "update":
		return errors.New("not implemented")
	case "delete":
		return errors.New(" not implemented")
	default:
		return errors.New("error unrecognized option")
	}

	return nil

}
