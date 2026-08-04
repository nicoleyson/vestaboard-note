BIN     := note
CMD     := ./cmd/note
DESTDIR := /usr/local/bin

VESTABOARD_DIR ?= $(shell pwd)
NOTE_BIN       ?= $(VESTABOARD_DIR)/$(BIN)
GIT_SHA        := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)

.PHONY: build install uninstall run-weather run-clock run-calendar tidy lint status \
        cron cron-init cron-update cron-uninstall _cron-render _cron-apply

# ── Build ────────────────────────────────────────────────────────────────────

build:
	go build -ldflags "-X main.version=$(GIT_SHA)" -o $(BIN) $(CMD)

install: build
	@echo "Installing $(BIN) to $(DESTDIR)/$(BIN)..."
	install -m 755 $(BIN) $(DESTDIR)/$(BIN)
	@echo ""
	@echo "Done. Add shell completion by running one of:"
	@echo "  Bash:  note completion bash  >> ~/.bash_profile"
	@echo "  Zsh:   note completion zsh   >> ~/.zshrc"
	@echo "  Fish:  note completion fish  > ~/.config/fish/completions/note.fish"
	@echo ""
	@echo "Set up cron jobs with:"
	@echo "  make cron"

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

lint:
	go vet ./...
	@which staticcheck > /dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

status: build
	./$(BIN) status

# ── Cron ─────────────────────────────────────────────────────────────────────
# make cron          — smart: init on first run, update on subsequent runs
# make cron-init     — first-time install (aborts if already installed)
# make cron-update   — always replace with current template
# make cron-uninstall — remove vestaboard cron entries

# Render the template → /tmp/vestaboard-cron-render
# {{USER}} is unused (user crontab doesn't need username column)
_cron-render:
	@sed \
	  -e 's|{{VESTABOARD_DIR}}|$(VESTABOARD_DIR)|g' \
	  -e 's|{{VESTABOARD_LOG}}|$(VESTABOARD_DIR)/vestaboard.log|g' \
	  -e 's|{{NOTE_BIN}}|$(NOTE_BIN)|g' \
	  -e 's|{{USER}}||g' \
	  crontab.template > /tmp/vestaboard-cron-render

# Apply /tmp/vestaboard-cron-render to the user crontab (works on both macOS and Linux)
_cron-apply:
	@( \
	  crontab -l 2>/dev/null \
	  | awk '/# BEGIN VESTABOARD/{found=1} !found{print} /# END VESTABOARD/{found=0}'; \
	  echo '# BEGIN VESTABOARD'; \
	  cat /tmp/vestaboard-cron-render; \
	  echo '# END VESTABOARD' \
	) | crontab -
	@echo "Installed to user crontab"

cron: _cron-render
	@if crontab -l 2>/dev/null | grep -q '# BEGIN VESTABOARD'; then \
	  $(MAKE) _cron-apply && echo "Updated existing vestaboard crontab block"; \
	else \
	  $(MAKE) _cron-apply && echo "Added vestaboard block to crontab"; \
	fi

cron-init: _cron-render
	@if crontab -l 2>/dev/null | grep -q '# BEGIN VESTABOARD'; then \
	  echo "Vestaboard crontab block already exists. Use 'make cron-update' to overwrite."; exit 1; \
	fi
	$(MAKE) _cron-apply

cron-update: _cron-render _cron-apply

cron-uninstall:
	@( crontab -l 2>/dev/null \
	  | awk '/# BEGIN VESTABOARD/{found=1} !found{print} /# END VESTABOARD/{found=0}' \
	) | crontab -
	@echo "Removed vestaboard block from crontab"
