## dstrm get schema-code

Get schema code archive by schema-ref

### Synopsis

In order to simplify sending correctly serialized data to STRM Privacy
it is recommended to use generated source code that defines a
class/object structure in a certain programming language, that knows
(with help of some open-source libraries) how to serialize objects.

The result of a `get schema-code` is a zip file with some source code
files for getting started with sending events in a certain programming
language. Generally this will be code where you’ll have to do some sort
of `build` step in order to make this fully operational in your
development setting (using a JDK, a Python or a Node.js environment).

### Usage

```
dstrm get schema-code (schema-ref) [flags]
```

### Options

```
  -h, --help                 help for schema-code
      --language string      which programming language
      --output-file string   Destination zip file location
      --overwrite            do we allow overwriting an existing file
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

