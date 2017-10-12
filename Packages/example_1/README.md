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


# Overall Learnings over the whole example


### General

- There is a difference between `go build` and `go install`. `go build` rebuilds the application without any caching for any of the packages. `go install` actually rebuilds the application with caching to ensure compile time is still fast even though there are numerous modules within the service.
- Apparently, if you did not put in GOBIN in the path, it won't properly know where to build and dump the binary file. It will keep complaining. You can resolve this by cd-ing into the same folder as the file that contains main and run the build from there.
- Go apparently introduced a `-o` flag. (Not exactly sure of which version). That would allow you to identify where to dump the generated bin folders.
- Go Dep is an official experimental tool that would sooner or later be merged (hopefully) be merged into the Go toolchain. It would be good to start experimenting with it in order to do package management within go.
- It is possible to generate more than 1 kind of output from a set of packages. If the top contains package main, it would use that as entry and compile accordingly
- You cannot put files that is supposed to belong to other packages in the same folder. E.g. A file with package miao and another file with package heh in a miao folder. Compiler will complain about heh.
- If file ends with test, it is not part of compilation step. It will be ignored
- A general good advice is to actually go read the source code for golang on github - can see how its being written
- For packages with subpackages, the imports still take from the root of the project... e.g If there is a subpackage hey in the miao package, to import it, need the statement: `import "miao/hey"`. The functions within the subpackage is referred to via hey namespace instead.


### Getting documentation

- Run the following command: `godoc -http=127.0.0.1:6060` in order to go to a locally hosted go doc website. Within it will contain some documentation on packages that you have not published on golang website.
- Running the command `go doc miao` allows you to immediately have a bird's eye view on whats available in the package. Struct is hidden but if you run `go doc miao.<Struct_name>`, the documention for it will pop out.


### Learnings on miao package?

- If you define a interface in the package, it will be available everywhere in the whole package but not across package
- A package is defined as everything in the folder (also defined by the package by the value at the top).