# Full observability App

This Golang application is meant to cover Logging, Metrics, Distributed Tracing and Profiling.

# Features of this app

- [x] Distributed Tracing
- [x] Logging linked to Distributed Tracing
- [x] Metrics
- [ ] Metrics with examplar support
- [x] Profiling support

Configuration - maybe use HPA to create large number of pods? 
- [ ] Generation of Large amount of logs
- [ ] Generation of Large amount of metrics
- [ ] Generation of Large amount of traces

# Commands

```bash
# Profiling application
# Ensure graphviz is installed as well
go tool pprof -http=localhost:8500 http://localhost:6060/debug/pprof/heap
```
