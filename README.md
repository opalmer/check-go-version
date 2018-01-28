# About

This project is intended to process information about release versions of
Go from the following API:

https://storage.googleapis.com/golang/

In general it's best to build and run the latest official release of Go as
it may contain security fixes. The intent of this project is two parts:

  * Provide an easy API for accessing releases.
  * Wrap these APIs in a binary that can be used as part of a build process.

## Usage

### Command Line

```
$ go get github.com/opalmer/check-go-version
$ check-go-version
 latest: Version{Name: go1.9.3.linux-amd64, Version: 1.9.3, Platform: linux, Architecture: amd64}
running: Version{Name: go1.9.3.linux-amd64, Version: 1.9.3, Platform: linux, Architecture: amd64}
```

### API

This package contains a self-contained API to retrieve information about
Golang versions. For more information, see godoc:

https://godoc.org/github.com/opalmer/check-go-version/api

### Caching

Note, this project provides a local cache of the response retrieved from
the bucket that's updated if the last query was > 30 minutes go. This can
behavior can be changed however by modifying one of the `BucketCache*`
variables in the api package.
