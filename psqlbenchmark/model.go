package psqlbenchmark

import (
	"time"

	"github.com/uptrace/bun"
)

type UUIDModel struct {
	bun.BaseModel `bun:"uuid_model,alias:uuid_model"`
	ID            string `bun:"id,nullzero,notnull,pk,type:uuid,default:gen_random_uuid()"`
}

type IdentityModel struct {
	bun.BaseModel `bun:"identity_model,alias:identity_model"`
	ID            int64  `bun:"id,nullzero,notnull,pk,identity,type:bigint"`
	Type          string `bun:"type,nullzero,notnull,type:varchar(16)"`
}

type BigSerialModel struct {
	bun.BaseModel `bun:"bigserial_model,alias:bigserial_model"`
	ID            int64  `bun:"id,nullzero,notnull,type:bigserial"`
	Type          string `bun:"type,nullzero,notnull,type:varchar(16)"`
}

type TimestampedModel struct {
	bun.BaseModel `bun:"timestamped_model,alias:timestamped_model"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,type:timestamptz,default:current_timestamp"`
}

type NanosecondModel struct {
	bun.BaseModel    `bun:"nanosecond_model,alias:nanosecond_model"`
	CreatedAtSeconds int64 `bun:"created_at_seconds,notnull,type:bigint"`
	CreatedAtNanos   int   `bun:"created_at_nanos,notnull,type:bigint"`
}
