package psqlbenchmark

import (
	"context"
	"fmt"
	"math/rand"
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

var intToType = map[int]string{
	0: "a",
	1: "b",
	2: "c",
	3: "d",
	4: "e",
	5: "f",
	6: "g",
}

func GenerateIdentityLoad(count int, ch chan<- *IdentityModel) {
	typeInt := rand.Intn(7)
	modelType := intToType[typeInt]
	for i := 0; i < count; i++ {
		ch <- &IdentityModel{
			Type: modelType,
		}
	}
}

func GenerateBigSerialLoad(count int, ch chan<- *BigSerialModel) {
	typeInt := rand.Intn(7)
	modelType := intToType[typeInt]
	for i := 0; i < count; i++ {
		ch <- &BigSerialModel{
			Type: modelType,
		}
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

func InsertOneIdentity(ctx context.Context, db *bun.DB, model *IdentityModel) error {
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

// IdentityInserter is a worker that receives a IdentityModel through a channel and inserts it into the database
func IdentityInserter(ctx context.Context, db *bun.DB, models <-chan *IdentityModel, done chan<- bool) {
	fmt.Println("IdentityInserter started...")
	for {
		select {
		case model := <-models:
			if err := InsertOneIdentity(ctx, db, model); err != nil {
				panic(err)
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Println("IdentityInserter done")
			done <- true
		}
	}
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

// IdentityReader runs a read query on the database every second and queries the last 1000 rows
func IdentityReader(ctx context.Context, db *bun.DB, done <-chan bool) {
	fmt.Println("IdentityReader started...")
	for {
		select {
		case <-done:
			fmt.Println("IdentityReader done")
			return
		case <-time.After(1 * time.Second):
			var models []*IdentityModel
			if err := db.NewSelect().Model(&models).Order("id DESC").Limit(1000).Scan(ctx); err != nil {
				panic(err)
			}
			fmt.Println("IdentityReader read", len(models), "rows")
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
		case <-time.After(500 * time.Millisecond):
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
