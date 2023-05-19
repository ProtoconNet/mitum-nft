package timestamp

import (
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TimeStampItemJSONMarshaler struct {
	hint.BaseHinter
	ProjectID         string `json:"projectid"`
	RequestTimeStamp  uint64 `json:"request_timestamp"`
	ResponseTimeStamp uint64 `json:"response_timestamp"`
	TimeStampID       uint64 `json:"timestampid"`
	Data              string `json:"data"`
}

func (t TimeStampItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TimeStampItemJSONMarshaler{
		BaseHinter:        t.BaseHinter,
		ProjectID:         t.projectID,
		RequestTimeStamp:  t.requestTimeStamp,
		ResponseTimeStamp: t.responseTimeStamp,
		TimeStampID:       t.timestampID,
		Data:              t.data,
	})
}

type TimeStampItemJSONUnmarshaler struct {
	Hint              hint.Hint `json:"_hint"`
	ProjectID         string    `json:"projectid"`
	RequestTimeStamp  uint64    `json:"request_timestamp"`
	ResponseTimeStamp uint64    `json:"response_timestamp"`
	TimeStampID       uint64    `json:"timestampid"`
	Data              string    `json:"data"`
}

func (t *TimeStampItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of NFT")

	var u TimeStampItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	return t.unmarshal(enc, u.Hint, u.ProjectID, u.RequestTimeStamp, u.ResponseTimeStamp, u.TimeStampID, u.Data)
}
