package test

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"testing"
)

func TestBolt(t *testing.T) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte("answer"), []byte("42"))
		_ = b.Put([]byte("answer1"), []byte("1"))
		_ = b.Put([]byte("answer2"), []byte("2"))
		_ = b.Put([]byte("1"), []byte("2"))
		_ = b.Put([]byte("2"), []byte("2"))
		return err
	})

	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte("answer"))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})

	_ = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBucket"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("first key=%s, value=%s\n", k, v)
		}
		return nil
	})

	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		_ = b.Delete([]byte("answer"))
		_ = b.Delete([]byte("1"))
		fmt.Println("delete key=answer")
		fmt.Println("delete key=1")
		return err
	})

	_ = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBucket"))
		cc := b.Cursor()
		for k, v := cc.First(); k != nil; k, v = cc.Next() {
			fmt.Printf("second key=%s, value=%s\n", k, v)
		}
		return nil
	})

	defer db.Close()
}
