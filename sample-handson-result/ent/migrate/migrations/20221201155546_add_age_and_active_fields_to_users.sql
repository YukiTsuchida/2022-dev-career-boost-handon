-- +goose Up
-- modify "users" table
ALTER TABLE "users" ADD COLUMN "age" bigint NULL, ADD COLUMN "active" boolean NOT NULL DEFAULT true;

-- +goose Down
-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "active", DROP COLUMN "age";
