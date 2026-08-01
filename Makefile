BIN     := note
CMD     := ./cmd/note
DESTDIR := /usr/local/bin

.PHONY: build install uninstall run-weather run-clock run-calendar tidy

build:
	go build -o $(BIN) $(CMD)

install: build
	@echo "Installing $(BIN) to $(DESTDIR)/$(BIN)..."
	install -m 755 $(BIN) $(DESTDIR)/$(BIN)
	@echo ""
	@echo "Done. Add shell completion by running one of:"
	@echo "  Bash:  note completion bash  >> ~/.bash_profile"
	@echo "  Zsh:   note completion zsh   >> ~/.zshrc"
	@echo "  Fish:  note completion fish  > ~/.config/fish/completions/note.fish"
	@echo ""
	@echo "Suggested crontab entries (run: crontab -e):"
	@echo "  */30 * * * * $(DESTDIR)/$(BIN) weather"
	@echo "  0    * * * * $(DESTDIR)/$(BIN) calendar"
	@echo "  15,45 * * * * $(DESTDIR)/$(BIN) clock"

uninstall:
	rm -f $(DESTDIR)/$(BIN)

run-weather: build
	./$(BIN) weather

run-clock: build
	./$(BIN) clock

run-calendar: build
	./$(BIN) calendar

tidy:
	go mod tidy
