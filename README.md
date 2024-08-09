# Easy Go Proxy

A quick and easy proxy server for development purposes.

**Note:** This is not for production use. It exists only to help circumvent CORS issues during development.

---
## Releases
Windows binaries are available in the releases section.

---
To build the binary, run `go build -o easy-go-proxy.exe main.go` if on Windows

After setting the GOOS and GOARCH environment variables,

`GOOS=linux GOARCH=amd64 go build -o easy-go-proxy main.go` if on Linux

see [here](https://golang.org/doc/install/source#environment) for more information on these options.

---

## Usage
In the terminal of your choice, run the following command:
```
./easy-go-proxy -port=[port] -target=[url] -lport=[localPort]
```

### Flags:
- `-port`: The port on which the proxy server will run. Default: 8781.
- `-target`: The target URL to which the proxy server will forward requests. Default: "http://localhost"
- `-lport`: The local port on which the proxy server will run. Default: 8782.

