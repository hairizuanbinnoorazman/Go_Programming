# Example of using own packages

This is an example of using your own packages to modularize and split up your codebase.

You will need to store all of your codebase into the src folder of the project.


# Naive Command Lines

1. Export GOPATH. Go tooling depends quite heavily on the environment variables being set. Try running the below command
```bash
export GOPATH="$(pwd)"
```

2. Go to the directory that contains the `hello.go` file. Run the `go install` command. That's kind of it.

If you want to take things further...

3. Set the GOBIN as well.
```bash
export GOBIN="$(pwd)/bin"
```

4. Run the following set of commands.
```bash
go install src/hello/hello.go
go install src/meh/meh.go
```
It will create 2 binary executables in the bin folder. hello and meh.


# Learnings

- There is a difference between `go build` and `go install`. `go build` rebuilds the application without any caching for any of the packages. `go install` actually rebuilds the application with caching to ensure compile time is still fast even though there are numerous modules within the service.
- Apparently, if you did not put in GOBIN in the path, it won't properly know where to build and dump the binary file. It will keep complaining. You can resolve this by cd-ing into the same folder as the file that contains main and run the build from there.
- Go apparently introduced a `-o` flag. (Not exactly sure of which version). That would allow you to identify where to dump the generated bin folders.
- Go Dep is an official experimental tool that would sooner or later be merged (hopefully) be merged into the Go toolchain.