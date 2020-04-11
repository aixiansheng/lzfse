package lzfse

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestSmall(t *testing.T) {
	DoDecomp("cmp.lz", "dec", t)
}

func TestMedium(t *testing.T) {
	DoDecomp("cmp2.lz", "dec2", t)
}

func DoDecomp(compressed, original string, t *testing.T) {
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

	outBytes := make([]byte, len(decBytes))
	outBuffer := bytes.NewBuffer(outBytes)

	d := NewReader(cmp)

	n, err := io.Copy(outBuffer, d)
	if int(n) != len(outBytes) {
		t.Errorf("len(outBytes) != n:  %d != %d err (%v)", len(outBytes), int(n), err)
	}

	if err != nil {
		t.Errorf("io.Copy should have returned EOF, instead it returned %v", err)
	}
}
