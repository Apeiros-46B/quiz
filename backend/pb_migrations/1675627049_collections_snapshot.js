migrate((db) => {
  const snapshot = [
    {
      "id": "90khgnzbjh4enu5",
      "created": "2023-02-04 02:27:29.860Z",
      "updated": "2023-02-04 23:31:37.107Z",
      "name": "attempts",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "znxm0s05",
          "name": "userid",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "ztx8wg9p",
          "name": "time",
          "type": "number",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null
          }
        },
        {
          "system": false,
          "id": "nxtxyxtt",
          "name": "total_time",
          "type": "number",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null
          }
        },
        {
          "system": false,
          "id": "kty4zbs3",
          "name": "correct",
          "type": "bool",
          "required": true,
          "unique": false,
          "options": {}
        }
      ],
      "listRule": null,
      "viewRule": null,
      "createRule": null,
      "updateRule": null,
      "deleteRule": null,
      "options": {}
    },
    {
      "id": "24z5x4dlxzzy8dh",
      "created": "2023-02-04 02:38:41.764Z",
      "updated": "2023-02-04 06:16:59.054Z",
      "name": "questions",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "foolsn6x",
          "name": "index",
          "type": "number",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null
          }
        },
        {
          "system": false,
          "id": "6blgxsgx",
          "name": "question",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "xg7x2oga",
          "name": "choices",
          "type": "json",
          "required": true,
          "unique": false,
          "options": {}
        }
      ],
      "listRule": "",
      "viewRule": "",
      "createRule": null,
      "updateRule": null,
      "deleteRule": null,
      "options": {}
    },
    {
      "id": "u7bbinfv3ml8c4i",
      "created": "2023-02-04 02:39:39.444Z",
      "updated": "2023-02-04 06:15:46.728Z",
      "name": "correct",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "k9gt0rqe",
          "name": "question",
          "type": "relation",
          "required": true,
          "unique": false,
          "options": {
            "maxSelect": 1,
            "collectionId": "24z5x4dlxzzy8dh",
            "cascadeDelete": true
          }
        },
        {
          "system": false,
          "id": "uveynsou",
          "name": "correct",
          "type": "number",
          "required": false,
          "unique": false,
          "options": {
            "min": null,
            "max": null
          }
        }
      ],
      "listRule": null,
      "viewRule": null,
      "createRule": null,
      "updateRule": null,
      "deleteRule": null,
      "options": {}
    },
    {
      "id": "6k75mdzun955ne3",
      "created": "2023-02-04 02:49:20.345Z",
      "updated": "2023-02-04 02:58:10.981Z",
      "name": "settings",
      "type": "base",
      "system": false,
      "schema": [
        {
          "system": false,
          "id": "zksygskl",
          "name": "key",
          "type": "text",
          "required": true,
          "unique": false,
          "options": {
            "min": null,
            "max": null,
            "pattern": ""
          }
        },
        {
          "system": false,
          "id": "uzechwpv",
          "name": "value",
          "type": "json",
          "required": true,
          "unique": false,
          "options": {}
        }
      ],
      "listRule": "key != \"successKW\"",
      "viewRule": "key != \"successKW\"",
      "createRule": null,
      "updateRule": null,
      "deleteRule": null,
      "options": {}
    }
  ];

  const collections = snapshot.map((item) => new Collection(item));

  return Dao(db).importCollections(collections, true, null);
}, (db) => {
  return null;
})
