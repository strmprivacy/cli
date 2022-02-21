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