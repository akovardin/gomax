package tcp

import (
	"fmt"
	"sync"

	lz4 "github.com/pierrec/lz4/v4"
	"github.com/klauspost/compress/zstd"
)

type Lz4BlockCompression struct{}

func (c *Lz4BlockCompression) Decompress(src []byte, maxOutput int) ([]byte, error) {
	if maxOutput <= 0 {
		maxOutput = 5 * 1024 * 1024
	}
	dst := make([]byte, maxOutput)
	n, err := lz4.UncompressBlock(src, dst)
	if err != nil {
		return nil, fmt.Errorf("lz4: %w", err)
	}
	return dst[:n], nil
}

func (c *Lz4BlockCompression) Compress(src []byte) []byte {
	bound := lz4.CompressBlockBound(len(src))
	dst := make([]byte, bound)
	n, err := lz4.CompressBlock(src, dst, nil)
	if err != nil {
		return src
	}
	return dst[:n]
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
