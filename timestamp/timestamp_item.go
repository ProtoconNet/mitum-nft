package timestamp

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	MaxProjectIDLen = 10
	MaxDataLen      = 1024
)

var TimeStampItemHint = hint.MustNewHint("mitum-timestamp-item-v0.0.1")

type TimeStampItem struct {
	hint.BaseHinter
	projectID         string
	requestTimeStamp  uint64
	responseTimeStamp uint64
	timestampID       uint64
	data              string
}

func NewTimeStampItem(
	pid string,
	reqTS,
	resTS,
	tID uint64,
	data string,
) TimeStampItem {
	return TimeStampItem{
		BaseHinter:        hint.NewBaseHinter(TimeStampItemHint),
		projectID:         pid,
		requestTimeStamp:  reqTS,
		responseTimeStamp: resTS,
		timestampID:       tID,
		data:              data,
	}
}

func (t TimeStampItem) IsValid([]byte) error {
	if len(t.projectID) < 1 || len(t.projectID) > MaxProjectIDLen {
		return errors.Errorf("invalid projectID length %v < 1 or > %v", len(t.projectID), MaxProjectIDLen)
	}

	if len(t.data) < 1 || len(t.data) > MaxDataLen {
		return errors.Errorf("invalid data length %v < 1 or > %v", len(t.data), MaxDataLen)
	}

	return nil
}

func (t TimeStampItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(t.projectID),
		util.Uint64ToBytes(t.requestTimeStamp),
		util.Uint64ToBytes(t.responseTimeStamp),
		util.Uint64ToBytes(t.timestampID),
		[]byte(t.data),
	)
}

func (t TimeStampItem) ProjectID() string {
	return t.projectID
}

func (t TimeStampItem) RequestTimeStamp() uint64 {
	return t.requestTimeStamp
}

func (t TimeStampItem) ResponseTimeStamp() uint64 {
	return t.responseTimeStamp
}

func (t TimeStampItem) TimestampID() uint64 {
	return t.timestampID
}

func (t TimeStampItem) Data() string {
	return t.data
}

func (t TimeStampItem) Equal(ct TimeStampItem) bool {
	if t.projectID != ct.projectID {
		return false
	}

	if t.requestTimeStamp != ct.requestTimeStamp {
		return false
	}

	if t.responseTimeStamp != ct.responseTimeStamp {
		return false
	}

	if t.timestampID != ct.timestampID {
		return false
	}

	if t.data != ct.data {
		return false
	}

	return true
}
