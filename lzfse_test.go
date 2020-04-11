package lzfse

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestSmall(t *testing.T) {
	DoDecomp("cmp.lz", "dec", "dec.err", t)
}

func TestMedium(t *testing.T) {
	DoDecomp("cmp2.lz", "dec2", "dec2.err", t)
}

func TestKern(t *testing.T) {
	DoDecomp("kernel.lzfse", "kernel.dec", "kernel.err", t)
}

func DoDecomp(compressed, original, errorOutputFile string, t *testing.T) {
	cmp, err := os.Open(compressed)
	if err != nil {
		t.Errorf("Couldn't open test file")
	}
	defer cmp.Close()

	dec, err := os.Open(original)
	if err != nil {
		t.Errorf("Couldn't open decompressed file")
	}
	defer dec.Close()

	decBytes, err := ioutil.ReadAll(dec)
	if err != nil {
		t.Errorf("Couldn't readall original")
	}

	var buffer bytes.Buffer

	d := NewReader(cmp)

	if n, err := io.Copy(&buffer, d); err != nil {
		t.Errorf("Error decompressing: %v [orig= %d new=%d]", err, len(decBytes), n)
	}

	if !bytes.Equal(buffer.Bytes(), decBytes) {
		t.Errorf("The outputs did not match")
		ioutil.WriteFile(errorOutputFile, buffer.Bytes(), 0644)
	}
}
