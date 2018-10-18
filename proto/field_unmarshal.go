package proto

import (
	"math"
	"reflect"
)

func UnmarshalBool(b *Buffer, wt WireType, ret *bool) error {
	switch wt {
	case WireVarint:
		v, err := b.DecodeVarint()
		if err != nil {
			return err
		}

		if v != 0 {
			*ret = true
		} else {
			*ret = false
		}
	default:
		return ErrBadWireType
	}

	return nil
}

func UnmarshalInt64(b *Buffer, wt WireType, ret *int64) error {

	switch wt {
	case WireVarint:
		v, err := b.DecodeVarint()
		if err != nil {
			return err
		}
		*ret = int64(v)

	case WireZigzag64:
		v, err := b.DecodeZigzag64()
		if err != nil {
			break
		}

		*ret = int64(v)
	default:
		return ErrBadWireType
	}

	return nil
}

func UnmarshalInt32(b *Buffer, wt WireType, ret *int32) error {

	switch wt {
	case WireVarint:
		v, err := b.DecodeVarint()
		if err != nil {
			return err
		}
		*ret = int32(v)

	case WireZigzag32:
		v, err := b.DecodeZigzag32()
		if err != nil {
			break
		}

		*ret = int32(v)
	default:
		return ErrBadWireType
	}

	return nil
}

func UnmarshalUInt32(b *Buffer, wt WireType, ret *uint32) error {

	switch wt {
	case WireVarint:
		v, err := b.DecodeVarint()
		if err != nil {
			return err
		}
		*ret = uint32(v)

	case WireZigzag32:
		v, err := b.DecodeZigzag32()
		if err != nil {
			break
		}

		*ret = uint32(v)
	default:
		return ErrBadWireType
	}

	return nil
}

func UnmarshalUInt64(b *Buffer, wt WireType, ret *uint64) error {

	switch wt {
	case WireVarint:
		v, err := b.DecodeVarint()
		if err != nil {
			return err
		}
		*ret = v

	case WireZigzag64:
		v, err := b.DecodeZigzag64()
		if err != nil {
			break
		}

		*ret = v
	default:
		return ErrBadWireType
	}

	return nil
}

func UnmarshalFloat32(b *Buffer, wt WireType, ret *float32) error {

	switch wt {
	case WireZigzag32:
		v, err := b.DecodeFixed32()
		if err != nil {
			return err
		}

		*ret = math.Float32frombits(uint32(v))
	default:
		return ErrBadWireType
	}

	return nil
}

func UnmarshalFloat64(b *Buffer, wt WireType, ret *float64) error {

	switch wt {
	case WireZigzag64:
		v, err := b.DecodeFixed64()
		if err != nil {
			return err
		}

		*ret = math.Float64frombits(uint64(v))
	default:
		return ErrBadWireType
	}

	return nil
}

func UnmarshalString(b *Buffer, wt WireType, ret *string) error {
	switch wt {
	case WireBytes:
		v, err := b.DecodeStringBytes()
		if err != nil {
			return err
		}

		*ret = v
	default:
		return ErrBadWireType
	}

	return nil
}

func rawUnmarshalStruct(b *Buffer, msg Struct) error {

	for {
		wireTag, err := b.DecodeVarint()

		if err != nil {
			return err
		}

		fieldIndex, wt := parseWireTag(wireTag)

		err = msg.Unmarshal(b, fieldIndex, wt)

		if err != nil {
			return err
		}
	}

	return nil
}

func UnmarshalStruct(b *Buffer, wt WireType, msgPtr interface{}) error {

	switch wt {
	case WireBytes:
		size, err := b.DecodeVarint()
		if err != nil {
			return err
		}

		limitBuffer := NewBuffer(b.ConsumeBytes(int(size)))

		// MyType
		msgType := reflect.TypeOf(msgPtr).Elem().Elem()

		msgIns := reflect.New(msgType)
		msgValuePtr := reflect.ValueOf(msgPtr)
		msgValuePtr.Elem().Set(msgIns)

		return rawUnmarshalStruct(limitBuffer, msgIns.Interface().(Struct))

	default:
		return ErrBadWireType
	}
}