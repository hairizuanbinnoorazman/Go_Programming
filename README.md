## Programming with Go

Code snippets with go code

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

Run the command to run check documentation
```bash
go doc strconv Atoi
```

Run a local webserver to look at web based documentation 
```bash
godoc -http=:8080
```

## Go Code Principles

- be articulate, concise.
- provide a simple API.
- have precise comments.
- be readable, **top-down code.**


## Go Anti patterns

- https://about.sourcegraph.com/go/idiomatic-go/

Please don't do the following:
- Tiny packages
- Premature Exportation (As much as possible don't expose too much of your code. This keeps the package simple and easy to understand.)
- Package util (Seriously don't do this. What is util supposed to do? It doesn't enhance readability etc)
- Config structs (Massive data structures that fly around. Making tracing hard. It makes the developer have to work hard to understand the code)
- Pointer all things? (Bad idea. It might be good for performance but measure before implementing it)
- No no to context.Value
- Dont panic for code in libraries
- No to blank interface{} (Why would you hide your implementation. It makes it harder to read. Instead, look around to see if there is any interface that you can potential use. Or create your own, but it has to be homegrown though)


## Random Advice

Random advice that makes sense:

- Do not test system software created via Go on production environment
- If you really need to run some go command line in production, make sure you have some sort of backup plan (in case sth catches fire)
- For functions, pass in interfaces but output structs. This would allow one to be explicit in what functions the person would need - e.g. I need a addAll() function and a removeAll() function to be created. Mantra: Accept native + interfaces but return structs
- Don't ever use the global variables unless you're thinking of commiting suicide via variable tracing - reason for not doing is that the variable can be manipulated by any function within the package which makes it really really dangerous



## Other Resources

- https://talks.golang.org/2014/readability.slide
- https://golang.org/doc/effective_go.html
- https://github.com/golang/go/wiki/CodeReviewComments#Package_Comments
- https://ahmet.im/blog/golang-json-decoder-pitfalls/ (Some resources, for some weird reason use the encoder/decoder rather marshal/unmarshal)
- https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702
- https://godoc.org/golang.org/x/tools/cmd

Program your next server in Go (Sameer Ajmani)
- https://www.youtube.com/watch?v=5bYO60-qYOI
- https://talks.golang.org/2016/applicative.slide#1
- https://github.com/golang/talks


