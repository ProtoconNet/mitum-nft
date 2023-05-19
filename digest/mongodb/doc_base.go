package mongodbstorage

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v2/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type Doc interface {
	ID() interface{}
}

type BaseDoc struct {
	id          interface{}
	encoderHint hint.Hint
	data        interface{}
	isHinted    bool
}

func NewBaseDoc(id, data interface{}, enc encoder.Encoder) (BaseDoc, error) {
	_, isHinted := data.(hint.Hinter)

	return BaseDoc{
		id:          id,
		encoderHint: enc.Hint(),
		isHinted:    isHinted,
		data:        data,
	}, nil
}

func (do BaseDoc) ID() interface{} {
	return do.id
}

func (do BaseDoc) M() (bson.M, error) {
	m := bson.M{
		"_e":      do.encoderHint.String(),
		"_hinted": do.isHinted,
		"d":       do.data,
	}

	if do.id != nil {
		m["_id"] = do.id
	}

	return m, nil
}

type BaseDocBSONUnMarshaler struct {
	I bson.Raw      `bson:"_id,omitempty"`
	E string        `bson:"_e"`
	D bson.RawValue `bson:"d"`
	H bool          `bson:"_hinted"`
}

func LoadDataFromDoc(b []byte, encs *encoder.Encoders) (bson.Raw /* id */, interface{} /* data */, error) {
	var bd BaseDocBSONUnMarshaler
	if err := bsonenc.Unmarshal(b, &bd); err != nil {
		return nil, nil, err
	}

	ht, err := hint.ParseHint(bd.E)
	if err != nil {
		return nil, nil, err
	}

	enc := encs.Find(ht)
	if enc == nil {
		return nil, nil, util.ErrNotFound.Errorf("encoder not found for %q", bsonenc.BSONEncoderHint)
	}

	if !bd.H {
		return bd.I, bd.D, nil
	}

	doc, ok := bd.D.DocumentOK()
	if !ok {
		return nil, nil, errors.Errorf("hinted should be mongodb Document")
	}

	var data interface{}
	if i, err := enc.Decode([]byte(doc)); err != nil {
		return nil, nil, err
	} else {
		data = i
	}

	return bd.I, data, nil
}

type BaseManifestDocBSONUnMarshaler struct {
	I bson.Raw      `bson:"_id,omitempty"`
	E string        `bson:"_e"`
	D bson.RawValue `bson:"d"`
	H bool          `bson:"_hinted"`
	O uint64        `bson:"operations"`
}

func LoadManifestDataFromDoc(b []byte, encs *encoder.Encoders) (bson.Raw /* id */, interface{} /* data */, uint64 /* operations */, error) {

	var bd BaseManifestDocBSONUnMarshaler
	if err := bsonenc.Unmarshal(b, &bd); err != nil {
		return nil, nil, 0, err
	}

	ht, err := hint.ParseHint(bd.E)
	if err != nil {
		return nil, nil, 0, err
	}

	enc := encs.Find(ht)
	if enc == nil {
		return nil, nil, 0, util.ErrNotFound.Errorf("encoder not found for %q", bsonenc.BSONEncoderHint)
	}

	if !bd.H {
		return bd.I, bd.D, 0, nil
	}

	doc, ok := bd.D.DocumentOK()
	if !ok {
		return nil, nil, 0, errors.Errorf("hinted should be mongodb Document")
	}

	var data interface{}
	if i, err := enc.Decode([]byte(doc)); err != nil {
		return nil, nil, 0, err
	} else {
		data = i
	}

	return bd.I, data, bd.O, nil
}
