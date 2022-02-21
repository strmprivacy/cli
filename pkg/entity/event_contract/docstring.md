An Event Contract defines the rules that are to be applied to events.

The Event Contract defines:

-   the Schema to use via a full Schema reference (handle/name/version)

-   the key field

-   the PII fields

-   any validations on fields (e.g. a regex to validate an email
    address)

Like Schemas, Event Contracts can be private or public, allowing them to
be found and used by others than the owning client. Be careful, public
Event Contracts cannot be deleted.

Also like Schemas, Event Contracts are versioned using a versioning
scheme that can be fully determined by the client. The only restrictions
are that version numbers:

-   MUST follow the semantic version format (major/minor/patch),

-   MUST always be ascending

An Event Contract is uniquely identified by its Event Contract
reference, in the format (organization handle/event contract
name/version).

### Usage