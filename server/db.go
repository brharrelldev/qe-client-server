package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brharrelldev/qe-client-server/pkg/types"
	"github.com/dgraph-io/badger/v4"
	uuid "github.com/satori/go.uuid"
)

type DB struct {
	db *badger.DB
}

func NewDB() (*DB, error) {

	db, err := badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		return nil, fmt.Errorf("error opening up badger DB %v", err)
	}

	return &DB{db: db}, nil

}

func (db *DB) Create(data *bufio.Reader) (*types.Response, error) {

	buf := make([]byte, 1024)

	n, err := data.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading data into buffer %v", err)
	}

	id := uuid.NewV4().String()

	dataInBytes := buf[:n]

	fmt.Println("data in bytes", string(dataInBytes))

	if err := db.db.Update(func(txn *badger.Txn) error {

		entry := badger.NewEntry([]byte(id), dataInBytes)
		if err := txn.SetEntry(entry); err != nil {
			return fmt.Errorf("error writing to database %v", err)
		}

		return nil

	}); err != nil {

		return nil, fmt.Errorf("error creating new entry in inventory %v", err)

	}

	fmt.Println("new record successfully created!")

	return &types.Response{
		Message: fmt.Sprintf("%s successfull created", id),
	}, nil
}

func (db *DB) List() (*types.Response, error) {

	fmt.Println("testing")
	var dataList []types.Payload

	if err := db.db.View(func(txn *badger.Txn) error {

		fmt.Println("inside iterator")

		iter := txn.NewIterator(badger.DefaultIteratorOptions)

		defer iter.Close()

		for iter.Rewind(); iter.Valid(); iter.Next() {
			fmt.Println("inside loop")
			entry := iter.Item()

			key := entry.Key()

			val, err := entry.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("error copying value %v", err)
			}

			fmt.Println(string(key))

			fmt.Println("value", string(val))

			var data types.Data

			if err := json.NewDecoder(bytes.NewBuffer(val)).Decode(&data); err != nil {
				fmt.Println("error deserializing data", err)
				continue
			}

			payload := types.Payload{
				Id:   string(key),
				Data: &data,
			}

			dataList = append(dataList, payload)

		}

		return nil

	}); err != nil {

		return nil, fmt.Errorf("error retrieving data %v", err)

	}

	resp := &types.Response{
		Message: "list request successful",
		Results: dataList,
	}

	return resp, nil

}

func (db *DB) Update(id string, data *types.Data) (*types.Response, error) {

	if err := db.db.Update(func(txn *badger.Txn) error {

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error serializing data %v", err)
		}

		if err := txn.Set([]byte(id), dataBytes); err != nil {
			return fmt.Errorf("error updating entry %v", err)
		}

		return nil

	}); err != nil {

		return nil, fmt.Errorf("error updating new entry %v", err)

	}

	return &types.Response{
		Message: "update successful",
	}, nil

}

func (db *DB) Delete(id string) error {

	if err := db.db.Update(func(txn *badger.Txn) error {

		iID := []byte(id)
		if err := txn.Delete(iID); err != nil {
			return fmt.Errorf("error deleting the key id %s %v", iID, err)
		}
		return nil

	}); err != nil {
		return fmt.Errorf("error updating entry %v", err)
	}

	return nil

}
