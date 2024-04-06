package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	conn *pgxpool.Pool
}

func NewDatabase() *Database {
	url := fmt.Sprintf("postgres://system:%s@localhost:5432/system?sslmode=disable", os.Getenv("PGPASSWORD"))

	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}

	_, err = conn.Exec(context.Background(), `
		CREATE ANALYZER three_gram_analyzer (
			TOKENIZER three_gram with ( type = 'ngram', min_gram  = 3, max_gram = 3)
		);

		CREATE TABLE if not exists vehicles (
			id string primary key,
			vin text,
			description text,
			use text,
			INDEX vin_grams using fulltext(vin) with (analyzer = 'three_gram_analyzer'),
			INDEX text_fields USING FULLTEXT(description, use)
		);

		CREATE TABLE if not exists vehicle_measurements (
			vehicle_id string,
			measurement_type string,
			measurement_time timestamp,
			measurement_value double,
			measurement_partition as DATE_BIN('1 minute'::INTERVAL, measurement_time, 0)
		) CLUSTERED BY (measurement_time) INTO 4 SHARDS PARTITIONED BY (measurement_partition);

		CREATE TABLE if not exists vehicle_current_locations (
			vehicle_id string primary key,
			vehicle_location geo_point INDEX USING "geohash" WITH (precision='5m'),
			received_at timestamp
		);

		CREATE TABLE if not exists countries (
			name text,
			country_code text primary key,
			shape geo_shape INDEX USING "geohash" WITH (precision='100m'),
			capital text,
			capital_location geo_point
		) WITH (number_of_replicas=0);
	`)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to init schema: %v\n", err)
		panic(err)
	}

	return &Database{conn: conn}
}

func (d *Database) Close() {
	d.conn.Close()
}
