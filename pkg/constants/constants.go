package constants

import "strings"

const DefaultConfigFilename = "strm"
const DefaultConfigFilenameLength = len(DefaultConfigFilename) + 5 // suffix .json or .yaml

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
