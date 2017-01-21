package datastore

import (
	"encoding/json"
	"errors"
	"log"
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
	out := Day{}

	// Start the transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return out, err
	}

	year := tx.Bucket([]byte(y))
	if year == nil {
		return out, errors.New("Bucket does not exist for: " + y)
	}

	month := year.Bucket([]byte(m))
	if month == nil {
		return out, errors.New("Bucket does not exist for: " + y)
	}

	day := month.Bucket([]byte(d))
	if day == nil {
		return out, errors.New("Bucket does not exist for: " + y)
	}

	if err := day.ForEach(func(k, v []byte) error {
		var tmp Event

		err := json.Unmarshal(v, &tmp)
		if err != nil {
			log.Println(err)
		}
		out = append(out, tmp)
		return nil
	}); err != nil {
		return nil, err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
