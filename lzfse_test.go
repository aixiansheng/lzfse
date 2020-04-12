package lzfse

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Run make -f test.mk to generate the data.
func TestVariousSizes(t *testing.T) {
	if testFile, err := os.Open("test.list"); err == nil {
		defer testFile.Close()
		scanner := bufio.NewScanner(testFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			for _, compressedInput := range strings.Fields(scanner.Text()) {
				decompressedInput := strings.TrimSuffix(compressedInput, ".cmp")
				errorFile := decompressedInput + ".err"
				DoDecomp(compressedInput, decompressedInput, errorFile, t)
			}
		}
	}
}

func DoDecomp(compressed, original, errorOutputFile string, t *testing.T) {
	t.Logf("Testing lzfse on %s -> %s (error will be in %s)", compressed, original, errorOutputFile)

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
