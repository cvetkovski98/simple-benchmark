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

func GenerateNanosecondLoad(count int, ch chan<- *NanosecondModel) {
	for i := 0; i < count; i++ {
		t := time.Now()
		m := &NanosecondModel{
			CreatedAtSeconds: t.Unix(),
			CreatedAtNanos:   t.Nanosecond(),
		}
		ch <- m
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

func InsertOneNanosecond(ctx context.Context, db *bun.DB, model *NanosecondModel) error {
	if _, err := db.NewInsert().Model(model).Exec(ctx); err != nil {
		fmt.Printf("InsertOneNanosecond failed with: %v for model: %v\n", err, model)
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
			fmt.Println("BigSerialInserter done")
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

func NanosecondInserter(ctx context.Context, db *bun.DB, models <-chan *NanosecondModel, done chan<- bool) {
	fmt.Println("NanosecondInserter started...")
	for {
		select {
		case model := <-models:
			if err := InsertOneNanosecond(ctx, db, model); err != nil {
				fmt.Printf("Error inserting nanosecond model: %v, model: %v", err, model)
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Println("NanosecondInserter done")
			done <- true
		}
	}
}

// BigSerialReader runs a read query on the database every second and queries the last 1000 rows
func BigSerialReader(ctx context.Context, db *bun.DB, done <-chan bool) {
	fmt.Println("BigSerialReader started...")
	for {
		select {
		case <-done:
			fmt.Println("BigSerialReader done")
			return
		case <-time.After(1 * time.Second):
			var models []*BigSerialModel
			if err := db.NewSelect().Model(&models).Order("id DESC").Limit(1000).Scan(ctx); err != nil {
				panic(err)
			}
			fmt.Println("BigSerialReader read", len(models), "rows")
		}
	}
}

// NanoReader runs a read query on the database every second and queries the last 1000 rows
func NanoReader(ctx context.Context, db *bun.DB, done <-chan bool) {
	fmt.Println("NanoReader started...")
	for {
		select {
		case <-done:
			fmt.Println("NanoReader done")
			return
		case <-time.After(500 * time.Millisecond):
			var models []*NanosecondModel
			if err := db.NewSelect().
				Model(&models).
				Order("created_at_seconds DESC", "created_at_nanos DESC").
				Limit(1000).
				Scan(ctx); err != nil {
				panic(err)
			}
			fmt.Println("NanoReader read", len(models), "rows")
		}
	}
}
