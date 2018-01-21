# About

This project is intended to process information about release versions of
Go from the following API:

https://storage.googleapis.com/golang/

In general it's best to build and run the latest official release of Go as
it may contain security fixes. The intent of this project is two parts:

  * Provide an easy API for accessing releases.
  * Wrap these APIs in a binary that can be used as part of a build process.


## API

This package contains a self-contained API to retrieve information about
Golang versions.
