CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "gender" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "room_id" int
);

CREATE TABLE "room" (
  "room_id" bigserial PRIMARY KEY,
  "member" varchar[],
  "date" int NOT NULL,
  "time" int NOT NULL
);

CREATE INDEX ON "users" ("name");

CREATE INDEX ON "room" ("room_id");

ALTER TABLE "users" ADD FOREIGN KEY ("room_id") REFERENCES "room" ("room_id");
