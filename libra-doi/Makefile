GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet
BINNAME = cmd
DEPLOYNAME = bootstrap

build: cmdline

all: cmdline deployable

cmdline:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -tags cmdline -o bin/$(BINNAME)

deployable:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -tags lambda.norpc,lambda -o bin/$(DEPLOYNAME)
	cd bin; zip deployment.zip $(DEPLOYNAME)

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
