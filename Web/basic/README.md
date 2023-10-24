# Basic Golang Application

Contains basic example in order to deal with control a Golang application via Environment variables. The following repository serves to provide the basic application to deploy the same app into various different environments and hopefully have 0 issues between them.

Here is the full list of environment supported for this:

## Golang App on local workstation

Build the application and then run it

```bash
CGO_ENABLED=0 go build -o app main.go
```