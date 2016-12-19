package datastore

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

type Day []Event

type Event struct {
	ID        uint64
	Name      string
	StartTime string
	EndTime   string
}

func EventToDb(db *bolt.DB, y string, m string, d string, e Event) error {
	// Start the transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	year, err := tx.CreateBucketIfNotExists([]byte(y))
	if err != nil {
		return err
	}

	month, err := year.CreateBucketIfNotExists([]byte(m))
	if err != nil {
		return err
	}

	day, err := month.CreateBucketIfNotExists([]byte(d))
	if err != nil {
		return err
	}

	// Generate an ID for the new event.
	eventID, err := day.NextSequence()
	if err != nil {
		return err
	}
	e.ID = eventID

	// Marshal and save the encoded event.
	if buf, err := json.Marshal(e); err != nil {
		return err
	} else if err := day.Put([]byte(strconv.FormatUint(e.ID, 10)), buf); err != nil {
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func EventsForDay(db *bolt.DB, y string, m string, d string) (Day, error) {

	// TODO: create out Day var to hold
	// events from day

	// Start the transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	year := tx.Bucket([]byte(y))
	month := year.Bucket([]byte(m))
	day := month.Bucket([]byte(d))

	day.ForEach(func(k, v []byte) error {

		// TODO: unmarshal value and store in Event struct

		fmt.Printf("key=%s, value=%s\n", k, v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
