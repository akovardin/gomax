package types

import (
	"encoding/json"
	"strconv"

	"github.com/vmihailenco/msgpack/v5"
)

type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexInt(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*f = FlexInt(i)
		return nil
	}
	return &json.UnmarshalTypeError{Value: string(data)}
}

func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(f))
}

func (f *FlexInt) DecodeMsgpack(dec *msgpack.Decoder) error {
	v, err := dec.DecodeInterface()
	if err != nil {
		return err
	}
	switch n := v.(type) {
	case int64:
		*f = FlexInt(n)
	case int32:
		*f = FlexInt(n)
	case int16:
		*f = FlexInt(n)
	case int8:
		*f = FlexInt(n)
	case int:
		*f = FlexInt(n)
	case uint64:
		*f = FlexInt(n)
	case float64:
		*f = FlexInt(int64(n))
	case string:
		i, err := strconv.Atoi(n)
		if err != nil {
			return err
		}
		*f = FlexInt(i)
	default:
		*f = 0
	}
	return nil
}

func (f FlexInt) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeInt(int64(f))
}

func (f FlexInt) Int() int {
	return int(f)
}

func (f FlexInt) String() string {
	return strconv.Itoa(int(f))
}
