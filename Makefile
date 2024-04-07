migrate-up:
	goose -dir ./migrations postgres "host=localhost user=postgres password=password port=5432 dbname=postgres sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=localhost user=postgres password=password port=5432 dbname=postgres sslmode=disable" down
