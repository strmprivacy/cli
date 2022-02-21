## dstrm auth print-access-token

Print your current access-token to stdout

### Synopsis


Print the current (JWT) access token to the terminal that can be used in a http header. Note that the token is printed on `stdout`, and the Expiry and billing-id are on
`stderr` so it’s easy to capture the token for scripting use with

```bash
export token=$(strm auth access-token)
```

Note that this token might be expired, so a refresh may be required.
Use token as follows:
'Authorization: Bearer &lt;token&gt;' 

### Usage

```
dstrm auth print-access-token [flags]
```

### Options

```
  -h, --help   help for print-access-token
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

* [dstrm auth](dstrm_auth.md)	 - Authentication command

