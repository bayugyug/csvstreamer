BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
BUILD_TIME := $(shell date +%Y%m%d.%H%M%S)
BUILD_HASH := $(shell git log -1 2>/dev/null| head -n 1 | cut -d ' ' -f 2)
BUILD_NAME := csvstreamer

all: build

build :
	go get -v
	CGO_ENABLED=0 GOOS=linux go build -o $(BUILD_NAME) -a -tags netgo -installsuffix netgo -installsuffix cgo -v -ldflags "-X main.BuildTime=$(BUILD_TIME) " .

test : build
	go test ./... > testrun.txt
	golint > lint.txt
	go tool vet -v . > vet.txt
	gocov test github.com/bayugyug/csvstreamer | gocov-xml > coverage.xml
	go test ./... -bench=. -test.benchmem -v 2>/dev/null | gobench2plot > benchmarks.xml

testrun : clean test
	time go test -v -bench=. -benchmem -dummy >> testrun.txt 2>&1

prepare : build

preptest: test

clean:
	rm -f $(BUILD_NAME)
	rm -f benchmarks.xml coverage.xml vet.txt lint.txt testrun.txt

re: clean all

