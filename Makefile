BINARY  := adp
PREFIX  := $(HOME)/.local/bin

.PHONY: build install clean

build:
	go build -o $(BINARY) .

install: build
	@mkdir -p $(PREFIX)
	cp $(BINARY) $(PREFIX)/
	@echo "installed → $(PREFIX)/$(BINARY)"

clean:
	rm -f $(BINARY)
