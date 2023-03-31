package psqlbenchmark

import (
	"time"

	"github.com/uptrace/bun"
)

type UUIDModel struct {
	bun.BaseModel `bun:"uuid_model,alias:uuid_model"`
	ID            string `bun:"id,nullzero,notnull,pk,type:uuid,default:gen_random_uuid()"`
}

type BigSerialModel struct {
	bun.BaseModel `bun:"big_serial_model,alias:big_serial_model"`
	ID            int64 `bun:"id,nullzero,notnull,pk,identity,type:bigint"`
}

type TimestampedModel struct {
	bun.BaseModel `bun:"timestamped_model,alias:timestamped_model"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,type:timestamptz,default:current_timestamp"`
}
