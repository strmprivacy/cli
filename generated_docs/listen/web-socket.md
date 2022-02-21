## dstrm listen web-socket

Read events via the web-socket (not for production purposes)

### Synopsis

The global `listen` command is used for starting a Web Socket listener for a stream and output all events to the console.

This command can receive events from both Source Streams and Derived Streams.

### Usage

```
dstrm listen web-socket (stream-name) [flags]
```

### Options

```
      --client-id string       Client id to be used for receiving data
      --client-secret string   Client secret to be used for receiving data
  -h, --help                   help for web-socket
```

### Options inherited from parent commands

```
      --api-auth-url string            User authentication host (default "https://accounts.strmprivacy.io")
      --api-host string                API host and port (default "api.strmprivacy.io:443")
      --events-auth-url string         Event authentication host (default "https://sts.strmprivacy.io")
      --kafka-bootstrap-hosts string   Kafka bootstrap brokers, separated by comma (default "export-bootstrap.kafka.strmprivacy.io:9092")
  -o, --output string                  Output format [json, json-raw, table, plain] (default "table")
      --token-file string              Token file that contains an access token (default is $HOME/.config/strmprivacy/credentials-<api-auth-url>.json)
      --web-socket-url string          Websocket to receive events from (default "wss://websocket.strmprivacy.io/ws")
```

### SEE ALSO

* [dstrm listen](dstrm_listen.md)	 - Listen for events on a stream

