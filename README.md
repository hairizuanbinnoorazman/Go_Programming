## Programming with Go

Code code snippets with go code

## List of useful commands

Run the Go formatter
```bash
# Prints it out on the stdout
gofmt hw.go

# Echo it into a file
gofmt hw.go > hw2.go
```

Run the command to get help docs on go env variables
```bash
go help environment
```

Run the command to build 
```bash
# Build an executable called hw in the same folder
go build hw.go

# To run the executable
./hw
```

Run the command to run the go file without running a compile step
```bash
go run hw.go
```
