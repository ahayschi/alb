GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

BIN=./bin
BALANCE_MAIN=./cmd/balance


mac:
	env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN)/balance $(BALANCE_MAIN)

linux:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN)/balance $(BALANCE_MAIN)

test:
	$(GOTEST) $(MAIN)

clean:
	rm -r $(BIN)
