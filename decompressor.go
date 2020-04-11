package lzfse

import (
	"encoding/binary"
	"io"
)

type Magic uint32

const (
	LZFSE_NO_BLOCK_MAGIC             Magic = 0
	LZFSE_ENDOFSTREAM_BLOCK_MAGIC    Magic = 0x24787662
	LZFSE_UNCOMPRESSED_BLOCK_MAGIC   Magic = 0x2d787662
	LZFSE_COMPRESSEDV1_BLOCK_MAGIC   Magic = 0x31787662
	LZFSE_COMPRESSEDV2_BLOCK_MAGIC   Magic = 0x32787662
	LZFSE_COMPRESSEDLZVN_BLOCK_MAGIC Magic = 0x6e787662
	INVALID                                = 0xdeadbeef
)

type decompressor struct {
	r            *cachedReader
	pipeR        *io.PipeReader
	pipeW        *io.PipeWriter
	handlerError error
}

func decodeUncompressedBlock(r *cachedReader, w io.Writer) error {
	var n_raw_bytes uint32
	if err := binary.Read(r, binary.LittleEndian, &n_raw_bytes); err != nil {
		return err
	}

	if _, err := io.CopyN(w, r, int64(n_raw_bytes)); err != nil {
		return err
	}

	return nil
}

func readBlockMagic(r io.Reader) (magic Magic, err error) {
	err = binary.Read(r, binary.LittleEndian, &magic)
	return
}

type blockHandler func(*cachedReader, io.Writer) error

func (d *decompressor) handleBlock(handler blockHandler) (Magic, error) {
	if err := handler(d.r, d.pipeW); err != nil {
		return INVALID, err
	}

	return readBlockMagic(d.r)
}

func (d *decompressor) Read(b []byte) (int, error) {
	return d.pipeR.Read(b)
}

func NewReader(r io.Reader) *decompressor {
	pipeR, pipeW := io.Pipe()
	d := &decompressor{
		r:     newCachedReader(r),
		pipeR: pipeR,
		pipeW: pipeW,
	}

	go func() {
		var err error
		magic := LZFSE_NO_BLOCK_MAGIC

		for nil == err {
			switch magic {
			case LZFSE_NO_BLOCK_MAGIC:
				magic, err = readBlockMagic(d.r)
			case LZFSE_UNCOMPRESSED_BLOCK_MAGIC:
				magic, err = d.handleBlock(decodeUncompressedBlock)
			case LZFSE_COMPRESSEDV1_BLOCK_MAGIC:
				magic, err = d.handleBlock(decodeCompressedV1Block)
			case LZFSE_COMPRESSEDV2_BLOCK_MAGIC:
				magic, err = d.handleBlock(decodeCompressedV2Block)
			case LZFSE_COMPRESSEDLZVN_BLOCK_MAGIC:
				magic, err = d.handleBlock(decodeLZVNBlock)
			case LZFSE_ENDOFSTREAM_BLOCK_MAGIC:
				magic = LZFSE_ENDOFSTREAM_BLOCK_MAGIC
				err = io.EOF
			default:
				panic("Bad magic")
			}
		}

		d.handlerError = err
		d.pipeW.Close()
	}()

	return d
}
