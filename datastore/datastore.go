package datastore

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

// Day is a custom Slice of events
type Day []Event

// Event is a struct holding all the data for an event.
//
// This will be encoded into the BoltDB database for storage.
type Event struct {
	ID        uint64
	Name      string
	StartTime string
	EndTime   string
}

//EventToDb stores an event struct in the DB based on the date.
func EventToDb(db *bolt.DB, y string, m string, d string, e Event) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(y))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		year := tx.Bucket([]byte(y))
		_, err := year.CreateBucketIfNotExists([]byte(m))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		month := tx.Bucket([]byte(y + "/" + m))
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
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func EventsForDay(db *bolt.DB, y string, m string, d string) (Day, error) {
	out := Day{}

	// Get events from day.
	//
	// Start the transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	year := tx.Bucket([]byte(y))
	month := year.Bucket([]byte(m))
	day := month.Bucket([]byte(d))

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

	return out, nil
}
