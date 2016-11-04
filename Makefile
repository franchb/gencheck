COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
endif

all: generate build test cover install
install-deps:
	glide install
build: generate
	if [ ! -d bin ]; then mkdir bin; fi
	go build -v -o bin/gencheck ./gencheck
test: generate gen-test
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -v ./generator -race -cover -coverprofile=$(COVERAGEDIR)/generator.coverprofile
	go test -v ./ -race -cover -coverprofile=$(COVERAGEDIR)/gencheck.coverprofile
	# go test -v ./internal/gkexample -race -cover -coverprofile=$(COVERAGEDIR)/gkexample.coverprofile
cover:
	go tool cover -html=$(COVERAGEDIR)/generator.coverprofile -o $(COVERAGEDIR)/generator.html
	go tool cover -html=$(COVERAGEDIR)/gencheck.coverprofile -o $(COVERAGEDIR)/gencheck.html
	go tool cover -html=$(COVERAGEDIR)/gencheck.coverprofile -o $(COVERAGEDIR)/gkexample.html
tc: test cover
coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN)
clean:
	go clean
	rm -f bin/gencheck
	rm -rf coverage/

# go-bindata will take all of the template files and create readable assets from them in the executable.
# This way the templates are readable in atom (or another text editor with golang template language support)
generate:
	go-bindata -o generator/assets.go -pkg=generator template/*.tmpl

gen-test: build
	./bin/gencheck -f=./internal/example/example.go
	
install:
	go install ./gencheck

reinstall: build
	go install ./gencheck


phony: clean tc build
