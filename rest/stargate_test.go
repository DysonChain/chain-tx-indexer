package rest_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/likecoin/likecoin-chain-tx-indexer/db"
	. "github.com/likecoin/likecoin-chain-tx-indexer/rest"
	. "github.com/likecoin/likecoin-chain-tx-indexer/test"
)

type Response struct {
	Pagination   interface{}
	Txs          []interface{}
	Tx_responses []interface{}
}

func TestStargate(t *testing.T) {
	b := db.NewBatch(Conn, 10000)
	b.Batch.Queue(
		"INSERT INTO txs (height, tx_index, tx, events) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", 1, 1,
		[]byte(`
{
  "height": "1",
  "txhash": "AAAAAA",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "iscn_record",
          "attributes": [
            {
              "key": "iscn_id",
              "value": "iscn://testing/AAAAAA/1"
            },
            {
              "key": "iscn_id_prefix",
              "value": "iscn://testing/AAAAAA"
            },
            {
              "key": "owner",
              "value": "like1qyqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqewmlu9"
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            { "key": "action", "value": "create_iscn_record" },
            {
              "key": "sender",
              "value": "like1qyqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqewmlu9"
            },
            { "key": "module", "value": "iscn" },
            {
              "key": "sender",
              "value": "like1qyqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqewmlu9"
            }
          ]
        }
      ]
    }
  ],
  "tx": {
    "@type": "/cosmos.tx.v1beta1.Tx",
    "body": {
      "messages": [
        {
          "@type": "/likechain.iscn.MsgCreateIscnRecord",
          "from": "like1qyqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqewmlu9",
          "record": {
            "recordNotes": "",
            "contentFingerprints": [],
            "stakeholders": [],
            "contentMetadata": {}
          }
        }
      ],
      "memo": "AAAAAA",
      "timeout_height": "0",
      "extension_options": [],
      "non_critical_extension_options": []
    },
    "auth_info": { "fee": {} },
    "signatures": [""]
  },
  "timestamp": "2022-01-01T00:00:00Z",
  "events": []
}
`),
		[]string{`iscn_record.iscn_id="iscn://testing/AAAAAA/1"`},
	)
	err := b.Flush()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = CleanupTestData(Conn) }()

	req := httptest.NewRequest(
		"GET",
		STARGATE_ENDPOINT+"?events=iscn_record.iscn_id='iscn://testing/AAAAAA/1'", nil)
	res, body := request(req)
	if res.StatusCode != 200 {
		t.Fatal(body)
	}
	var result Response
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Txs) == 0 {
		t.Fatal("No response:", result)
	}
}
