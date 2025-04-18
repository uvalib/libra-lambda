# libra-doi

Lambda function for allocating and updating Datacite DOIs for Libra.

## Local Setup

`make build` will use main-cmdline.go and takes the flags below.

Required environment variables are located in config.go

``` bash
cd libra-doi
make build
cd bin
./cmd -messageid=123 \
  -source=librabus \
  -eventname=metadata/edit \
  -namespace=libraetd \
  -objid=oid:co3e8h0ki3ukr53r60ng
```

### Requirements

* Go 1.21+
