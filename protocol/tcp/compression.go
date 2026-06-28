package tcp

import (
	"fmt"
	"sync"

	"github.com/klauspost/compress/zstd"
)

type Lz4BlockCompression struct{}

func (c *Lz4BlockCompression) Decompress(src []byte, maxOutput int) ([]byte, error) {
	if maxOutput <= 0 {
		maxOutput = 5 * 1024 * 1024
	}

	dst := make([]byte, 0, len(src)*2)
	pos := 0

	for pos < len(src) {
		token := int(src[pos])
		pos++

		litLen := token >> 4
		if litLen == 15 {
			for pos < len(src) {
				b := int(src[pos])
				pos++
				litLen += b
				if b != 255 {
					break
				}
			}
		}

		if litLen > 0 {
			if pos+litLen > len(src) {
				return nil, fmt.Errorf("lz4: literal length out of bounds")
			}
			dst = append(dst, src[pos:pos+litLen]...)
			pos += litLen
			if len(dst) > maxOutput {
				return nil, fmt.Errorf("lz4: output too large")
			}
		}

		if pos >= len(src) {
			break
		}

		if pos+1 >= len(src) {
			return nil, fmt.Errorf("lz4: incomplete offset")
		}
		offset := int(src[pos]) | (int(src[pos+1]) << 8)
		pos += 2
		if offset == 0 {
			return nil, fmt.Errorf("lz4: zero offset")
		}

		matchLen := (token & 0x0F) + 4
		if (token & 0x0F) == 0x0F {
			for pos < len(src) {
				b := int(src[pos])
				pos++
				matchLen += b
				if b != 255 {
					break
				}
			}
		}

		matchPos := len(dst) - offset
		if matchPos < 0 {
			return nil, fmt.Errorf("lz4: match out of bounds")
		}

		for i := 0; i < matchLen; i++ {
			dst = append(dst, dst[matchPos+(i%offset)])
		}

		if len(dst) > maxOutput {
			return nil, fmt.Errorf("lz4: output too large")
		}
	}

	return dst, nil
}

func (c *Lz4BlockCompression) Compress(src []byte) []byte {
	dst := make([]byte, 0, len(src))
	pos := 0

	for pos < len(src) {
		litStart := pos
		for pos < len(src) && (pos-litStart) < 15 {
			pos++
		}
		litLen := pos - litStart
		token := (litLen << 4) & 0xF0

		matchOffset := 0
		matchLen := 0

		searchStart := litStart - 65535
		if searchStart < 0 {
			searchStart = 0
		}
		for i := searchStart; i < litStart; i++ {
			j := i
			k := litStart
			for j < i+65535 && k < len(src) && src[j] == src[k] {
				j++
				k++
			}
			if j-i > matchLen {
				matchOffset = litStart - i
				matchLen = j - i
			}
		}

		if matchLen >= 4 {
			token |= (matchLen - 4) & 0x0F
			dst = append(dst, byte(token))
			dst = append(dst, src[litStart:litStart+litLen]...)
			dst = append(dst, byte(matchOffset&0xFF), byte((matchOffset>>8)&0xFF))
			pos += matchLen
		} else {
			token |= litLen & 0x0F
			dst = append(dst, byte(token))
			dst = append(dst, src[litStart:litStart+litLen]...)
		}
	}

	return dst
}

type ZstdCompression struct {
	encoder *zstd.Encoder
	decoder *zstd.Decoder
	mu      sync.Mutex
}

func NewZstdCompression() (*ZstdCompression, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, fmt.Errorf("zstd: failed to create encoder: %w", err)
	}
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("zstd: failed to create decoder: %w", err)
	}
	return &ZstdCompression{
		encoder: encoder,
		decoder: decoder,
	}, nil
}

func (z *ZstdCompression) Decompress(src []byte, maxOutput int) ([]byte, error) {
	z.mu.Lock()
	defer z.mu.Unlock()

	out, err := z.decoder.DecodeAll(src, nil)
	if err != nil {
		return nil, fmt.Errorf("zstd decompress: %w", err)
	}
	if maxOutput > 0 && len(out) > maxOutput {
		return nil, fmt.Errorf("zstd: output too large (%d > %d)", len(out), maxOutput)
	}
	return out, nil
}

func (z *ZstdCompression) Compress(src []byte) ([]byte, error) {
	z.mu.Lock()
	defer z.mu.Unlock()

	return z.encoder.EncodeAll(src, nil), nil
}

func (z *ZstdCompression) Close() {
	z.encoder.Close()
	z.decoder.Close()
}
