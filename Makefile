run api:
	go run server/cmd/main.go

seed:
	go run server/cmd/seed/seed.go

fseed:
	go run ./server/cmd/seed/ --flush
