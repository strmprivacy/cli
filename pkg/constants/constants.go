package constants

import "strings"

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

var OutputFormatFlagAllowedValues = []string{OutputFormatJson, OutputFormatJsonRaw, OutputFormatTable, OutputFormatPlain}
var OutputFormatFlagAllowedValuesText = strings.Join(OutputFormatFlagAllowedValues, ", ")

var UsageOutputFormatFlagAllowedValues = []string{OutputFormatCsv, OutputFormatJson, OutputFormatJsonRaw}
var UsageOutputFormatFlagAllowedValuesText = strings.Join(UsageOutputFormatFlagAllowedValues, ", ")
