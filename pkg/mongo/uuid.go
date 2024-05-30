package mongo

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

var ErrMongoInternal = errors.New("mongo internal")

const uuidSubtype = byte(0x04)

var (
	tUUID = reflect.TypeOf(uuid.UUID{})

	mongoRegistry = bson.NewRegistry()
)

func init() {
	mongoRegistry.RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue))
	mongoRegistry.RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue))
}

func uuidEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	b, ok := val.Interface().(uuid.UUID)
	if !ok {
		return ErrMongoInternal
	}

	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var (
		data    []byte
		subtype byte
		err     error
	)

	switch vrType := vr.Type(); vrType {
	case bson.TypeBinary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("%w: unsupported binary subtype %v for UUID", ErrMongoInternal, subtype)
		}
	case bson.TypeNull:
		err = vr.ReadNull()
	case bson.TypeUndefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("%w: cannot decode %v into a UUID", ErrMongoInternal, vrType)
	}

	if err != nil {
		return err
	}

	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(uuid2))

	return nil
}
