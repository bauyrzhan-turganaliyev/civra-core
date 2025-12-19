up:
	docker compose up --build

down:
	docker compose down

test:
	go test ./...

run-gateway:
	go run ./cmd/gateway

run-kingdom:
	go run ./cmd/kingdom

run-economy:
	go run ./cmd/economy

run-market:
	go run ./cmd/market
