A stream is the central resource in STRM Privacy. Clients can connect to
a stream to send and to receive events. A stream can be either an "input
stream" or a "derived stream".

Events are always sent to an input stream. Sending events to a derived
stream is not possible. After validation and encryption of all PII
fields, STRM Privacy sends all events to the input stream. Clients
consuming from the input stream will see all events, but with all PII
fields encrypted.

Derived streams can be made on top of a input stream. A derived stream
is configured with one or more consent levels and it only receives
events with matching consent levels (see details about this matching
process here). The PII fields with matching consent levels are decrypted
and sent to the derived stream. Clients connecting to the derived stream
will only receive the events on this stream.

Every stream has its own set of access tokens. Connecting to an input
stream requires different credentials than when connecting to a derived
stream.

### Usage