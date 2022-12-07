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
  "keyField": "transactionId",
  "dataSubjectField": "userId",
  "validations": [{"field": "email", "type": "regex", "value": "..."}],
  "fieldMetadata": [
    {
      "fieldName": "email",
      "personalDataConfig": {"isPii": true, "isQuasiId": true, "purposeLevel": 1}
    },
    {
      "fieldName": "userId",
      "personalDataConfig": {"isPii": false, "isQuasiId": true}
    },
    {
      "fieldName": "userAgeGroup",
      "personalDataConfig": {"isPii": false, "isQuasiId": true},
      "statisticalDataType": "ORDINAL",
      "ordinalValues": ["child", "teenager", "adult", "senior"],
      "nullHandlingConfig": {"type": "DEFAULT_VALUE", "defaultValue": "adult"}
    }
  ]
}
`)
