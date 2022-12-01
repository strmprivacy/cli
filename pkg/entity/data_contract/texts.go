package data_contract

import "strmprivacy/strm/pkg/util"

var createExample = util.DedentTrim(`

Here's an example of a Simple Schema file

{
	"name": "GDDemo",
	"nodes": [
		{ "type": "STRING", "name": "transactionId", "required": true },
		{ "type": "STRING", "name": "userId"},
		{ "type": "STRING", "name": "email"},
		{ "type": "INTEGER", "name": "age"},
		{ "type": "FLOAT", "name": "height"},
		{ "type": "STRING", "name": "size"},
		{ "type": "INTEGER", "name": "transactionAmount"},
		{ "type": "STRING", "name": "items"},
		{ "type": "STRING", "name": "hairColor"},
		{ "type": "INTEGER", "name": "itemCount"},
		{ "type": "STRING", "name": "date"},
		{ "type": "INTEGER", "name": "purpose"}
	],
}
Here's an example of Data Contract definition file'

{
	"ref": { "handle": "strmprivacy", "name": "GDDemo", "version": "1.0.10" },
	"keyField": "transactionId",
	"validations": [
		{ "field": "email", "type": "regex", "value": "..." }
	],
	"metadata": {
		"title": "Schema used for the GDDemo",
		"description": "Somewhat valid e-commerce data",
		"industries": [ "e-commerce" ]
	},
	"fieldMetadata": [
	{
		"fieldName": "email",
		"personalDataConfig": { "isPii": true, "isQuasiId": true, "purposeLevel": 1 }
	},
	{
		"fieldName": "userId",
		"personalDataConfig": { "isPii": true, "isQuasiId": true, "purposeLevel": 1 }
	}
	]
}
`)
