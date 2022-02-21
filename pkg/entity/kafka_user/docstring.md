A Kafka User is a user on a Kafka Exporter, that can be used for
authentication when connecting to a Kafka Exporter. By default, every
Kafka Exporter gets one Kafka User upon creation, but these can be
added/removed later.

In the current data model, the user does not have a assignable name; it
is assigned upon creation. It’s still very low level. See the end of
this page for an example.

### Usage