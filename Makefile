BIN := note
CMD := ./cmd/note

.PHONY: build run-weather run-clock run-calendar tidy

build:
	go build -o $(BIN) $(CMD)

run-weather: build
	./$(BIN) weather

run-clock: build
	./$(BIN) clock

run-calendar: build
	./$(BIN) calendar

tidy:
	go mod tidy
