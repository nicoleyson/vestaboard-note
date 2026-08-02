BIN     := note
CMD     := ./cmd/note
DESTDIR := /usr/local/bin

VESTABOARD_DIR ?= $(shell pwd)
NOTE_BIN       ?= $(VESTABOARD_DIR)/$(BIN)
UNAME          := $(shell uname)
CRON_USER      := $(shell id -un)

.PHONY: build install uninstall run-weather run-clock run-calendar tidy \
        cron cron-init cron-update cron-uninstall _cron-render _cron-apply

# ── Build ────────────────────────────────────────────────────────────────────

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

# ── Cron ─────────────────────────────────────────────────────────────────────
# make cron          — smart: init on first run, update on subsequent runs
# make cron-init     — first-time install (aborts if already installed)
# make cron-update   — always replace with current template
# make cron-uninstall — remove vestaboard cron entries

# Render the template → /tmp/vestaboard-cron
# On Linux: add username column (required by /etc/cron.d/)
# On macOS: {{USER}} is blank (user crontab doesn't need it)
_cron-render:
ifeq ($(UNAME), Linux)
	@sed \
	  -e 's|{{VESTABOARD_DIR}}|$(VESTABOARD_DIR)|g' \
	  -e 's|{{NOTE_BIN}}|$(NOTE_BIN)|g' \
	  -e 's|{{USER}}|$(CRON_USER) |g' \
	  crontab.template > /tmp/vestaboard-cron
else
	@sed \
	  -e 's|{{VESTABOARD_DIR}}|$(VESTABOARD_DIR)|g' \
	  -e 's|{{NOTE_BIN}}|$(NOTE_BIN)|g' \
	  -e 's|{{USER}}||g' \
	  crontab.template > /tmp/vestaboard-cron
endif

# Apply /tmp/vestaboard-cron to the system
_cron-apply:
ifeq ($(UNAME), Linux)
	sudo cp /tmp/vestaboard-cron /etc/cron.d/vestaboard
	sudo chmod 644 /etc/cron.d/vestaboard
	@echo "Installed to /etc/cron.d/vestaboard"
else
	@( \
	  echo '# BEGIN VESTABOARD'; \
	  cat /tmp/vestaboard-cron; \
	  echo '# END VESTABOARD' \
	) | crontab -
	@echo "Installed to user crontab (macOS)"
endif

cron: _cron-render
ifeq ($(UNAME), Linux)
	@if [ -f /etc/cron.d/vestaboard ]; then \
	  $(MAKE) _cron-apply && echo "Updated existing /etc/cron.d/vestaboard"; \
	else \
	  $(MAKE) _cron-apply && echo "Created /etc/cron.d/vestaboard"; \
	fi
else
	@if crontab -l 2>/dev/null | grep -q '# BEGIN VESTABOARD'; then \
	  $(MAKE) _cron-apply && echo "Updated existing vestaboard crontab block"; \
	else \
	  $(MAKE) _cron-apply && echo "Added vestaboard block to crontab"; \
	fi
endif

cron-init: _cron-render
ifeq ($(UNAME), Linux)
	@if [ -f /etc/cron.d/vestaboard ]; then \
	  echo "Already installed at /etc/cron.d/vestaboard. Use 'make cron-update' to overwrite."; exit 1; \
	fi
	$(MAKE) _cron-apply
else
	@if crontab -l 2>/dev/null | grep -q '# BEGIN VESTABOARD'; then \
	  echo "Vestaboard crontab block already exists. Use 'make cron-update' to overwrite."; exit 1; \
	fi
	$(MAKE) _cron-apply
endif

cron-update: _cron-render _cron-apply

cron-uninstall:
ifeq ($(UNAME), Linux)
	sudo rm -f /etc/cron.d/vestaboard
	@echo "Removed /etc/cron.d/vestaboard"
else
	@crontab -r 2>/dev/null || true
	@echo "Removed vestaboard crontab"
endif
