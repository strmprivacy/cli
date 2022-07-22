package common

import "strings"

var ConfigPath string

var GitSha = "dev"
var Version = "dev"
var BuiltOn = "unknown"

// The environment variable prefix of all environment variables bound to our command line flags.
// For example, --api-host is bound to STRM_API_HOST
const EnvPrefix = "STRM"

const ClientIdFlag = "client-id"
const ClientSecretFlag = "client-secret"
const DefaultConfigFilename = "config"
const DefaultConfigFileSuffix = ".yaml"

var DefaultConfigFileContents = []byte(`# The following configuration options are reflected in the CLI's flags
# save: true
# events-auth-url: https://sts.strmprivacy.io
# events-api-url: https://events.strmprivacy.io/event
# api-auth-url: https://accounts.strmprivacy.io
# api-host: api.strmprivacy.io:443
# kafka-bootstrap-hosts: export-bootstrap.kafka.strmprivacy.io:9092
`)

const SavedEntitiesDirectory = "saved-entities"

const GetCommandName = "get"
const ListCommandName = "list"
const CreateCommandName = "create"
const DeleteCommandName = "delete"
const ActivateCommandName = "activate"
const ArchiveCommandName = "archive"
const InviteCommandName = "invite"

const RecursiveFlagName = "recursive"
const RecursiveFlagUsage = "Retrieve entities and their dependents"
const RecursiveFlagShorthand = "r"

const OutputFormatJson = "json"
const OutputFormatJsonRaw = "json-raw"
const OutputFormatTable = "table"
const OutputFormatPlain = "plain"
const OutputFormatCsv = "csv"
const OutputFormatFilepath = "path"

const OutputFormatFlag = "output"
const OutputFormatFlagShort = "o"

var OutputFormatFlagAllowedValues = []string{OutputFormatJson, OutputFormatJsonRaw, OutputFormatTable, OutputFormatPlain}
var OutputFormatFlagAllowedValuesText = strings.Join(OutputFormatFlagAllowedValues, ", ")

var UsageOutputFormatFlagAllowedValues = []string{OutputFormatCsv, OutputFormatJson, OutputFormatJsonRaw}
var UsageOutputFormatFlagAllowedValuesText = strings.Join(UsageOutputFormatFlagAllowedValues, ", ")

var ContextOutputFormatFlagAllowedValues = []string{OutputFormatJson, OutputFormatJsonRaw, OutputFormatFilepath}
var ContextOutputFormatFlagAllowedValuesText = strings.Join(ContextOutputFormatFlagAllowedValues, ", ")

var ConfigOutputFormatFlagAllowedValues = []string{OutputFormatPlain, OutputFormatJson}
var ConfigOutputFormatFlagAllowedValuesText = strings.Join(ConfigOutputFormatFlagAllowedValues, ", ")

var AccountOutputFormatFlagAllowedValues = []string{OutputFormatPlain, OutputFormatJsonRaw}
var AccountOutputFormatFlagAllowedValuesText = strings.Join(AccountOutputFormatFlagAllowedValues, ", ")

var ProjectOutputFormatFlagAllowedValues = []string{OutputFormatPlain}
var ProjectOutputFormatFlagAllowedValuesText = strings.Join(ProjectOutputFormatFlagAllowedValues, ", ")
