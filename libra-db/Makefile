GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet
BINNAME = cmd
COMMON = ../lambda-common
DEPLOYNAME = bootstrap

build: common cmdline

linux: common deployable

all: common cmdline deployable

cmdline:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -tags cmdline -o bin/$(BINNAME)

deployable:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -tags lambda.norpc,lambda -o bin/$(DEPLOYNAME)
	cd bin; zip -r deployment.zip $(DEPLOYNAME) data

common: download

download:
	cd bin && test -f migrate || curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && chmod +x migrate


clean:
	$(GOCLEAN)
	rm -rf bin

dep:
	$(GOGET) -u
	$(GOMOD) tidy
	$(GOMOD) verify

fmt:
	$(GOFMT)

vet:
	$(GOVET)
