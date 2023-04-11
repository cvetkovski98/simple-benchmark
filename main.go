package main

import (
	"context"
	"time"

	"github.com/cvetkovski98/psql-benchmark/psqlbenchmark"
	"github.com/cvetkovski98/psql-benchmark/psqlbenchmark/migrations"
)

func main() {
	var dsn = "postgres://root:root@localhost:5432/benchmarkdb?sslmode=disable"
	var ctx = context.Background()
	var db, err = psqlbenchmark.NewPgDb(dsn)
	if err != nil {
		panic(err)
	}
	if err := migrations.Run(ctx, db); err != nil {
		panic(err)
	}
	var N = 4
	var M = 1000
	var modelsPerWorker = 300

	// Create the necessary communication channels
	var timestampedModels = make(chan *psqlbenchmark.TimestampedModel)
	var uuidModels = make(chan *psqlbenchmark.UUIDModel)
	var identityModels = make(chan *psqlbenchmark.IdentityModel)
	var bigSerialModels = make(chan *psqlbenchmark.BigSerialModel)
	var nanosecondModels = make(chan *psqlbenchmark.NanosecondModel)
	var doneReaders = make(chan bool)
	var done = make(chan bool)

	// Start a reader that will read the models from the database
	// go psqlbenchmark.IdentityReader(ctx, db, doneReaders)
	go psqlbenchmark.BigSerialReader(ctx, db, doneReaders)
	// go psqlbenchmark.NanoReader(ctx, db, doneReaders)

	var start = time.Now()
	// Generate the load from M workers
	for i := 0; i < M; i++ {
		// go psqlbenchmark.GenerateTimestampedLoad(modelsPerWorker, timestampedModels)
		// go psqlbenchmark.GenerateUUIDLoad(modelsPerWorker, uuidModels)
		// go psqlbenchmark.GenerateIdentityLoad(modelsPerWorker, identityModels)
		go psqlbenchmark.GenerateBigSerialLoad(modelsPerWorker, bigSerialModels)
		// go psqlbenchmark.GenerateNanosecondLoad(modelsPerWorker, nanosecondModels)
	}

	// Start N workers that will insert the models into the database
	for i := 0; i < N; i++ {
		// go psqlbenchmark.TimestampedInserter(ctx, db, timestampedModels, done)
		// go psqlbenchmark.UUIDInserter(ctx, db, uuidModels, done)
		// go psqlbenchmark.IdentityInserter(ctx, db, identityModels, done)
		go psqlbenchmark.BigSerialInserter(ctx, db, bigSerialModels, done)
		// go psqlbenchmark.NanosecondInserter(ctx, db, nanosecondModels, done)
	}

	// Wait for all workers and readers to finish
	for i := 0; i < N; i++ {
		<-done
	}
	doneReaders <- true

	var elapsed = time.Since(start)
	println("Elapsed time:", elapsed)

	close(timestampedModels)
	close(uuidModels)
	close(identityModels)
	close(bigSerialModels)
	close(nanosecondModels)
	close(doneReaders)
	close(done)
}
