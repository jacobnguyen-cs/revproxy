# Go Reverse Proxy Learning Plan

Learn Go and backend fundamentals by building an HTTP reverse proxy, writing the forwarding logic by hand before using any library that does it.

The user is a data engineer (~2 years, mostly backend work) with no prior Go experience, learning Go while building this.

**Don't start with `httputil.ReverseProxy`.** The standard library's `httputil.NewSingleHostReverseProxy` gives a working proxy in about 10 lines, but you learn almost nothing from it. Write your own first. Later, read its source (`$(go env GOROOT)/src/net/http/httputil/reverseproxy.go`) and compare it to yours.

## Milestones

Each milestone ends in a working program. From milestone 4 on, run with `go run -race .`.

- [ ] **1. Hello server.** Use `net/http`, `http.HandleFunc` and `ListenAndServe`.
    - Learn: handlers, `http.ResponseWriter`, `*http.Request`.
- [ ] **2. Hand-written single-backend proxy.** Build the backend request with `http.NewRequestWithContext`, copy headers, send it with an `http.Client`, then copy the status, headers and body back with `io.Copy`.
    - Learn: error handling, `defer resp.Body.Close()`, streaming through `io.Reader`/`io.Writer` instead of buffering.
- [ ] **3. Get the details right.** Strip hop-by-hop headers (`Connection`, `Keep-Alive`, etc.), set `X-Forwarded-For`, add timeouts, return 502 when the backend is down, and pass context cancellation through when the client disconnects.
    - Learn: the HTTP details that separate real backend work from a toy.
- [ ] **4. Multiple backends with round-robin.** Pick the next backend on each request.
    - Learn: goroutines, `sync.Mutex` vs `sync/atomic`, data races.
- [ ] **5. Health checks.** A background goroutine on a `time.Ticker` marks backends up or down.
    - Learn: channels, `context` for shutdown, `select`.
- [ ] **6. Production basics.** Config from flags or a file, structured logging with `log/slog`, graceful shutdown with `server.Shutdown`, and middleware for request logging and metrics.
    - Learn: middleware composition, lifecycle management.
- [ ] **7. Tests.** Start fake backends inside tests with `net/http/httptest`.
    - Learn: table-driven tests, testing HTTP without real servers.

## Working with Claude

Default mode: the user writes the code and Claude reviews it. Don't write milestone code for the user unless they ask. If Claude writes the proxy, the user ends up with a proxy they don't understand.

| Mode | How it works | How much you learn |
| --- | --- | --- |
| You write, Claude reviews | You do a milestone; Claude critiques it and points out non-idiomatic Go | Most |
| Guided | Claude gives a skeleton with TODOs and hints per step | Medium |
| Explained | Claude writes each milestone with detailed comments; you change and extend it | Least, but fastest |
