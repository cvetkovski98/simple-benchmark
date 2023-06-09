DROP INDEX IF EXISTS "idx_uuid_model_id";
DROP INDEX IF EXISTS "idx_bigserial_model_id";
DROP INDEX IF EXISTS "idx_bigserial_model_type";
DROP INDEX IF EXISTS "idx_identity_model_id";
DROP INDEX IF EXISTS "idx_identity_model_type";
DROP INDEX IF EXISTS "idx_timestamped_model_created_at";
DROP INDEX IF EXISTS "idx_nanosecond_model_created_at";
DROP TABLE IF EXISTS "uuid_model";
DROP TABLE IF EXISTS "bigserial_model";
DROP TABLE IF EXISTS "identity_model";
DROP TABLE IF EXISTS "timestamped_model";
DROP TABLE IF EXISTS "nanosecond_model";
