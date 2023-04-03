package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brharrelldev/qe-client-server/pkg/types"
	"os"
)

func fileToJson(path string) (*types.Data, error) {

	buf := make([]byte, 1024)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening up file %v", err)
	}

	defer f.Close()

	offset, err := f.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading file %v", err)
	}

	data := buf[:offset]

	var inventoryData *types.Data

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&inventoryData); err != nil {
		return nil, fmt.Errorf("error decoding input data %v", err)
	}

	return inventoryData, nil

}
