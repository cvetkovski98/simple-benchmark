package psqlbenchmark

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

func GenerateTimestampedLoad(count int, ch chan<- *TimestampedModel) {
	for i := 0; i < count; i++ {
		ch <- &TimestampedModel{}
	}
}

func GenerateUUIDLoad(count int, ch chan<- *UUIDModel) {
	for i := 0; i < count; i++ {
		ch <- &UUIDModel{}
	}
}

func GenerateBigSerialLoad(count int, ch chan<- *BigSerialModel) {
	for i := 0; i < count; i++ {
		ch <- &BigSerialModel{}
	}
}

func InsertOneTimestamped(ctx context.Context, db *bun.DB, model *TimestampedModel) error {
	if _, err := db.NewInsert().Model(model).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func InsertOneUUID(ctx context.Context, db *bun.DB, model *UUIDModel) error {
	if _, err := db.NewInsert().Model(model).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func InsertOneBigSerial(ctx context.Context, db *bun.DB, model *BigSerialModel) error {
	if _, err := db.NewInsert().Model(model).Exec(ctx); err != nil {
		return err
	}
	return nil
}

// BigSerialInserter is a worker that receives a BigSerialModel through a channel and inserts it into the database
func BigSerialInserter(ctx context.Context, db *bun.DB, models <-chan *BigSerialModel, done chan<- bool) {
	fmt.Println("BigSerialInserter started...")
	for {
		select {
		case model := <-models:
			if err := InsertOneBigSerial(ctx, db, model); err != nil {
				panic(err)
			}
		case <-time.After(100 * time.Millisecond):
			done <- true
		}
	}
}

// UUIDInserter is a worker that receives a UUIDModel through a channel and inserts it into the database;
func UUIDInserter(ctx context.Context, db *bun.DB, models <-chan *UUIDModel, done chan<- bool) {
	fmt.Println("UUIDInserter started...")
	for {
		select {
		case model := <-models:
			if err := InsertOneUUID(ctx, db, model); err != nil {
				panic(err)
			}
		case <-time.After(100 * time.Millisecond):
			done <- true
		}
	}
}

// TimestampedInserter is a worker that receives a TimestampedModel through a channel and inserts it into the database;
func TimestampedInserter(ctx context.Context, db *bun.DB, models <-chan *TimestampedModel, done chan<- bool) {
	fmt.Println("TimestampedInserter started...")
	for {
		select {
		case model := <-models:
			if err := InsertOneTimestamped(ctx, db, model); err != nil {
				panic(err)
			}
		case <-time.After(100 * time.Millisecond):
			done <- true
		}
	}
}
