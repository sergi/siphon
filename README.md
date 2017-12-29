### Siphon - Stream commands output to your browser

Siphon is an application and a library consisting of a client and a server that communicate through UDP. The client `siph` accepts stdin and streams it to the server. There can be as many clients streaming to the server, each with its own id.

The server emits data via websockets to any listener that connects to it

### Usage

You can simply use the siphone CLI application, or use it programmatically as a library.

#### Command Line Usage

```sh
$ # Start a streaming client that sends command stdout output to a server.
$ # the id and adderss flags are both optional, with the id defaulting to a
$ # random string and the default server address defaulting to "127.0.0.1:1200"
$ yourprogram | siph client id="myClient" address="127.0.0.1:1200"

$ # Start a UDP/WebSockets server. The defaults for udp and websockets port are
$ # 1200 and 3000, respectively.
$ siph server --udp-port=1200 --ws-port=3000
```

To see other options and flags use `siph --help` or `siph <command> --help`.

#### As a library

To run as a server:

```go
// NewServer signature: NewServer(udpPort int, wsPort int) *Server
server := siphon.NewServer(1200, 3000)
err := server.Init()
if err != nil {
    fmt.Println(err)
    os.Exit(ExitCodeError)
}
```

To run as a client:
```go
// Init signature: Init(address string, id string, stream *bufio.Reader, emitOutput bool) error
err := siphon.Init("127.0.0.1:1200", "myclient", bufio.NewReader(os.Stdin), true)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(ExitCodeError)
}
```
### How it works

### License

See LICENSE
