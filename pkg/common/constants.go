package common

import "strings"

var ConfigPath string

// The environment variable prefix of all environment variables bound to our command line flags.
// For example, --api-host is bound to STRM_API_HOST
const EnvPrefix = "STRM"

const AuthSuccessHTML = `<html><head><meta http-equiv="refresh" content="0; url=https://streammachine.io"/></head><body></body></html>`

const DefaultConfigFilename = "strm"
const DefaultConfigFileSuffix = ".yaml"

var DefaultConfigFileContents = []byte(`# The following configuration options are reflected in the CLI's flags
# save: true
# event-auth-host: https://auth.strm.services
# events-gateway: https://in.strm.services/event
# api-auth-host: https://accounts.streammachine.io
# api-host: apis.streammachine.io:443
`)

const SavedEntitiesDirectory = "saved-entities"

const GetCommandName = "get"
const ListCommandName = "list"
const CreateCommandName = "create"
const DeleteCommandName = "delete"

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

var ConfigOutputFormatFlagAllowedValues = []string{OutputFormatPlain}
var ConfigOutputFormatFlagAllowedValuesText = strings.Join(ConfigOutputFormatFlagAllowedValues, ", ")
