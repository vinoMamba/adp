BINARY  := adp
PREFIX  := $(HOME)/.local/bin

SKILL_SRC  := adp-skill
SKILL_DATA := internal/skills/data

.PHONY: build install clean skills skillsCheck check

build:
	go build -o $(BINARY) .

install: build
	@mkdir -p $(PREFIX)
	cp $(BINARY) $(PREFIX)/
	@echo "installed → $(PREFIX)/$(BINARY)"

clean:
	rm -f $(BINARY)

# Sync the human-edited adp-skill/ (source of truth) into the embeddable
# internal/skills/data/ mirror. go:embed cannot reach adp-skill/ from
# internal/skills/ (no `..`), so the mirror is the embed source. The mirror is
# tracked in git so a plain `go build` works without running this first.
skills:
	@rm -rf $(SKILL_DATA)
	@mkdir -p $(SKILL_DATA)
	@cp -r $(SKILL_SRC)/. $(SKILL_DATA)/
	@find $(SKILL_DATA) -name '.DS_Store' -delete
	@echo "synced $(SKILL_SRC) → $(SKILL_DATA)"

# Fail if the mirror drifted from the source. Run before opening a PR.
skillsCheck:
	@diff -rq --exclude=.DS_Store $(SKILL_SRC) $(SKILL_DATA) >/dev/null 2>&1 || { \
		echo "skills mirror out of sync: run 'make skills'"; \
		diff -rq --exclude=.DS_Store $(SKILL_SRC) $(SKILL_DATA); \
		exit 1; }
	@echo "skills mirror in sync"

# Local pre-PR gate (mirrors CI).
check: skillsCheck
	go vet ./...
	go test ./...
