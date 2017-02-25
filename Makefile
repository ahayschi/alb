GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

BIN=./bin
BALANCE_MAIN=github.com/parallelworks/alb/cmd/balance


OUT=parallelworks-test-candidate

all: build

build:
	env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN)/balance $(BALANCE_MAIN)

test:
	$(GOTEST) $(MAIN)

clean:
	rm -r $(BIN)
