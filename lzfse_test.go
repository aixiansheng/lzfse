package lzfse

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"testing"
)

func allocInStreams(n int) ([]*inStream, error){
	ret := make([]*inStream, n)
	payload := make([]byte, 64)
	rand.Read(payload)

	for i := 0; i < n; i++ {
		var err error
		r := bytes.NewReader(payload)
		ret[i], err = newInStream(0, r)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func BenchmarkFsePull(b *testing.B) {
	inStreams, err := allocInStreams(b.N)
	if err != nil {
		b.FailNow()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := inStreams[i]
		for s.accum_nbits >= 1 {
			s.pull(1)
		}
	}
}

// Run make -C test/ to generate the data.
func TestVariousSizes(t *testing.T) {
	if testFile, err := os.Open("test/test.list"); err == nil {
		defer testFile.Close()
		scanner := bufio.NewScanner(testFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			for _, compressedInput := range strings.Fields(scanner.Text()) {
				decompressedInput := strings.TrimSuffix(compressedInput, ".cmp")
				errorFile := decompressedInput + ".err"
				t.Run(compressedInput, func(t *testing.T) {
					DoDecomp(compressedInput, decompressedInput, errorFile, t)
				})
			}
		}
	}
}

func DoDecomp(compressed, original, errorOutputFile string, t *testing.T) {
	ff, err := syscall.Open(compressed, syscall.O_RDONLY, 0)
	if err != nil {
		t.Errorf("couldn't open file")
	}
	var stat syscall.Stat_t
	err = syscall.Fstat(ff, &stat)
	if err != nil {
		t.Errorf("Couldn't fstat")
	}
	fd, err := syscall.Mmap(ff, 0, int(stat.Size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		t.Errorf("Couldn't mmap test file")
	}
	cmp := bytes.NewReader(fd)
//	cmp, err := os.Open(compressed)
//	if err != nil {
//		t.Errorf("Couldn't open test file")
//	}
//	defer cmp.Close()

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
