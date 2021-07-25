# libts

libts is a pure go wrapper for the TeamSpeak Serverquery protocol. Its goal is to ease the pain of interfacing with a TeamSpeak Server from a program.

## What it can do

libts is supposed to be a complete wrapper around the TeamSpeak Serverquery, so you don't have to actually deal with it. It implements most of the functionality of the Serverquery and is *much* easier to use as a developer than actually trying to deal with the Serverquery and its weird responses.

### Connecting

There are currently three ways to connect to a TeamSpeak Serverquery:

* Telnet
* SSH
* HTTP/HTTPS

Which one you choose is pretty much irrelevent for the most part as the Queries are nearly feature-equal. However there are some differences and a clear recommendation from my side.

#### Telnet

Telnet is that ancient protocol noone really uses anymore as it has no security whatsoever. It is, afaik, the predessesor of SSH and provides a remote terminal to another machine.

| Pros                                            | Cons                     |
|-------------------------------------------------|--------------------------|
| easy to use                                     | insecure                 |
| supports the whole functionality of Serverquery | plain-text communication |

```go
serverquery, err := NewServerQuery("my-cool-teamspeak.com", 10011, "serveradmin", "mySuperSecurePassword")
if err != nil {
    panic(err)
}
```

#### SSH

SSH was later added to TeamSpeak Server as a secure alternative to Telnet. It supports all the functionality Telnet does.

| Pros                                            | Cons |
|-------------------------------------------------|------|
| easy to use                                     |      |
| supports the whole functionality of Serverquery |      |
| secure                                          |      |

```go
serverquery, err := NewSSHQuery("my-cool-teamspeak.com", 10022, "serveradmin", "mySuperSecurePassword")
if err != nil {
    panic(err)
}
```

#### HTTP/HTTPS

As of late 2020 TeamSpeak Server provides a somewhat RESTful interface. That means you can query a TeamSpeak Server with pretty much everything, as long as it is connected to the Internet. Unfortunately since HTTP isn't a statful protocol, events are currently not supported. Also some other commands, which are obsolete when using HTTP/HTTPS, are also not supported such as `use` or `login`.

| Pros                              | Cons                 |
|-----------------------------------|----------------------|
| easy to use                       | not feature complete |
| can be queried by nearly anything | not really RESTful   |
| can be secured                    |                      |

```go
webClient := &http.Client{...}
serverquery := webquery.WebQuery{
		Host: "my-cool-teamspeak.com",
		Port: 10088,
		Key: "myAPIKey",
		TLS: false,
		HTTPClient: webClient,
	}
```

### Querying

After connecting to the TeamSpeak server you can start querying information from the TeamSpeak server.
You can either just send commands,

```go
globalMessage := libts.Request{
    Command: "gm",
    Args: map[string]interface{}{
        "msg": "Hello World",
    },
}
response, err := serverquery.DoRaw(globalMessage)
if err != nil {
    panic(err)
}
```

send them and process the response,

```go
hostinfo := libts.Request{
    Command: "hostinfo",
}
host := &struct {
    Uptime int `mapstructure:"instance_uptime"`
}
err := serverquery.Do(hostinfo, host)
if err != nil {
    panic(err)
}
fmt.Printf("Server is %ds up", host.Uptime)
```

or use the helper methods provided in the query package:

```go
queryAgent := query.Agent{
    Query: serverquery,
}
hostinfo, err := queryAgent.Host()
if err != nil {
    panic(err)
}
fmt.Printf("There are %d clients in %d channel on %d virtual server",
    hostinfo.ClientsOnlineTotal,
    hostinfo.ChannelOnlineTotal,
    hostinfo.VirtualServerRunningTotal)
```

