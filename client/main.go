package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brharrelldev/qe-client-server/pkg/types"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "qe-client"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "action",
			Aliases:  []string{"a"},
			Required: true,
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

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("error dialing server %v", err)
	}

	defer conn.Close()

	switch c.String("action") {
	case "list":

		req := types.Payload{
			Method: "list",
		}

		msg, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("error serializing request %v", err)
		}

		fmt.Println("list request sent to server")

		if _, err := conn.Write(msg); err != nil {
			return fmt.Errorf("error sending request to server %v", err)
		}

		output, err := io.ReadAll(conn)
		if err != nil {
			return fmt.Errorf("error reading data from server %v", err)
		}

		fmt.Println(string(output))

	case "create":
		if c.String("path") == "" {
			return errors.New("path is required for get request")
		}

		inventoryData, err := fileToJson(c.String("path"))
		if err != nil {
			return fmt.Errorf("error getting inventorty data from file %v", err)
		}

		fmt.Printf("%v\n", inventoryData)
		iID := uuid.NewV4()
		req := types.Payload{
			Id:     iID.String(),
			Method: "create",
			Data:   inventoryData,
		}

		msg, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("error marshaling request json %v", err)
		}

		if _, err := conn.Write(msg); err != nil {
			return fmt.Errorf("error sending data to server %v", err)
		}

		output, err := io.ReadAll(conn)
		if err != nil {
			return fmt.Errorf("erro reading connection %v", err)
		}

		fmt.Println(string(output))

	case "update":
		return errors.New("not implemented")
	case "delete":
		return errors.New(" not implemented")
	default:
		return errors.New("error unrecognized option")
	}

	return nil

}
