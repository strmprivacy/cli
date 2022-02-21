## dstrm list kafka-users

List Kafka users

### Synopsis

A Kafka User is a user on a Kafka Exporter, that can be used for
authentication when connecting to a Kafka Exporter. By default, every
Kafka Exporter gets one Kafka User upon creation, but these can be
added/removed later.

In the current data model, the user does not have a assignable name; it
is assigned upon creation. It’s still very low level. See the end of
this page for an example.

### Usage

```
dstrm list kafka-users [kafka-exporter-name] [flags]
```

### Options

```
  -h, --help   help for kafka-users
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

* [dstrm list](dstrm_list.md)	 - List entities

