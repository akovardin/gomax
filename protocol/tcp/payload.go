package tcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/vmihailenco/msgpack/v5"
)

type MsgpackPayloadCodec struct{}

func (c *MsgpackPayloadCodec) Encode(payload interface{}) ([]byte, error) {
	normalized := c.normalizeForMsgpack(payload)
	return msgpack.Marshal(normalized)
}

func (c *MsgpackPayloadCodec) normalizeForMsgpack(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		if hasMsgpackTag(rv.Type()) {
			return v
		}
		data, err := json.Marshal(v)
		if err != nil {
			return v
		}
		var result map[string]interface{}
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&result); err != nil {
			return v
		}
		for k, val := range result {
			result[k] = c.normalizeForMsgpack(val)
		}
		return result
	case reflect.Map:
		result := make(map[string]interface{})
		iter := rv.MapRange()
		for iter.Next() {
			key := fmt.Sprintf("%v", iter.Key().Interface())
			result[key] = c.normalizeForMsgpack(iter.Value().Interface())
		}
		return result
	case reflect.Slice, reflect.Array:
		if rv.Type() == reflect.TypeOf([]byte{}) {
			return v
		}
		result := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			result[i] = c.normalizeForMsgpack(rv.Index(i).Interface())
		}
		return result
	default:
		return convertJSONNumber(v)
	}
}

func (c *MsgpackPayloadCodec) Decode(payloadBytes []byte) (interface{}, error) {
	dec := msgpack.NewDecoder(bytes.NewReader(payloadBytes))
	dec.SetMapDecoder(func(d *msgpack.Decoder) (interface{}, error) {
		return d.DecodeUntypedMap()
	})
	result, err := dec.DecodeInterface()
	if err != nil {
		return nil, fmt.Errorf("msgpack decode: %w", err)
	}
	return result, nil
}

type TcpPayloadDecoder struct {
	Serializer       *MsgpackPayloadCodec
	Compression      *Lz4BlockCompression
	ZstdCompression  *ZstdCompression
}

func NewTcpPayloadDecoder(zstdCompression *ZstdCompression) *TcpPayloadDecoder {
	return &TcpPayloadDecoder{
		Serializer:       &MsgpackPayloadCodec{},
		Compression:      &Lz4BlockCompression{},
		ZstdCompression:  zstdCompression,
	}
}

func (d *TcpPayloadDecoder) Decode(payloadBytes []byte, flags int) (map[string]interface{}, error) {
	if len(payloadBytes) == 0 {
		return map[string]interface{}{}, nil
	}

	var decompressed []byte
	var err error

	if flags == 0xFF {
		if d.ZstdCompression == nil {
			return nil, fmt.Errorf("zstd compression not available")
		}
		decompressed, err = d.ZstdCompression.Decompress(payloadBytes, 0)
	} else if flags > 0 {
		decompressed, err = d.Compression.Decompress(payloadBytes, 5*1024*1024)
	} else {
		decompressed = payloadBytes
	}

	if err != nil {
		return nil, fmt.Errorf("decode decompress: %w", err)
	}

	decoded, err := d.Serializer.Decode(decompressed)
	if err != nil {
		head := decompressed
		if len(head) > 128 {
			head = head[:128]
		}
		tail := decompressed
		if len(tail) > 64 {
			tail = tail[len(tail)-64:]
		}
		return nil, fmt.Errorf("decode serialize: %w (flags=%d comp=%d decomp=%d head=%x tail=%x)",
			err, flags, len(payloadBytes), len(decompressed), head, tail)
	}

	normalized := normalizeKeys(decoded)
	if result, ok := normalized.(map[string]interface{}); ok {
		return result, nil
	}

	return nil, fmt.Errorf("decoded payload is not a map, got %T", decoded)
}

func hasMsgpackTag(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Field(i).Tag.Lookup("msgpack"); ok && tag != "" {
			return true
		}
	}
	return false
}

func convertJSONNumber(v interface{}) interface{} {
	if num, ok := v.(json.Number); ok {
		if i, err := num.Int64(); err == nil {
			return i
		}
		if f, err := num.Float64(); err == nil {
			return f
		}
		return string(num)
	}
	return v
}

func normalizeKeys(obj interface{}) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			result[key] = normalizeKeys(val)
		}
		return result
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			var keyStr string
			switch k := key.(type) {
			case string:
				keyStr = k
			case int:
				keyStr = strconv.Itoa(k)
			case int8:
				keyStr = strconv.FormatInt(int64(k), 10)
			case int16:
				keyStr = strconv.FormatInt(int64(k), 10)
			case int32:
				keyStr = strconv.FormatInt(int64(k), 10)
			case int64:
				keyStr = strconv.FormatInt(k, 10)
			case uint:
				keyStr = strconv.FormatUint(uint64(k), 10)
			case uint8:
				keyStr = strconv.FormatUint(uint64(k), 10)
			case uint16:
				keyStr = strconv.FormatUint(uint64(k), 10)
			case uint32:
				keyStr = strconv.FormatUint(uint64(k), 10)
			case uint64:
				keyStr = strconv.FormatUint(k, 10)
			case []byte:
				keyStr = string(k)
			default:
				keyStr = fmt.Sprintf("%v", key)
			}
			result[keyStr] = normalizeKeys(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = normalizeKeys(val)
		}
		return result
	default:
		return convertJSONNumber(v)
	}
}
