## dstrm get key-stream

Get key stream by name

### Synopsis

Key Streams are a restricted feature. For now, enabling and disabling
key streams can not be done through the Console or CLI. Contact us for
more information.

A Key Stream can be enabled on a stream and it contains all encryption
keys that are used on this stream. Normally STRM Privacy fully manages
and stores the encryption keys that are used, but with a key stream,
clients can get access to the keys and decrypt events themselves.

Usage of key streams places a lot more responsibility in the hands of
the client, so this feature requires careful consideration before using.

With regard to the data flow, STRM Privacy generates a new encryption
key whenever an event with a new "key link" (which can be seen as a
"session" concept, in that it links separate events together) is
received. This encryption key is stored internally, rotated after a
certain period and, if there is a key stream, put on the Key Stream
approximately at the same time as the event is put on the input stream.

### Usage

```
dstrm get key-stream [name] [flags]
```

### Options

```
  -h, --help   help for key-stream
```

### Options inherited from parent commands

```
      --api-auth-url string            User authentication host (default "https://accounts.strmprivacy.io")
      --api-host string                API host and port (default "api.strmprivacy.io:443")
      --events-auth-url string         Event authentication host (default "https://sts.strmprivacy.io")
      --kafka-bootstrap-hosts string   Kafka bootstrap brokers, separated by comma (default "export-bootstrap.kafka.strmprivacy.io:9092")
  -o, --output string                  Output format [json, json-raw, table, plain] (default "table")
  -r, --recursive                      Retrieve entities and their dependents
      --token-file string              Token file that contains an access token (default is $HOME/.config/strmprivacy/credentials-<api-auth-url>.json)
      --web-socket-url string          Websocket to receive events from (default "wss://websocket.strmprivacy.io/ws")
```

### SEE ALSO

* [dstrm get](dstrm_get.md)	 - Get an entity

