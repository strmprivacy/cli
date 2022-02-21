A Sink is a STRM Privacy configuration object for a remote file storage.
For now, AWS S3 and Google Cloud Storage Buckets are supported. By
itself, a Sink does nothing. A Batch Exporter needs to be connected to a
Sink and a Stream to start outputting events.

Upon creation, STRM Privacy validates whether or not the Bucket exists
and if it is accessible with the given credentials.

### Usage