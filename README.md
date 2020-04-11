# lzfse

[![Go](https://github.com/aixiansheng/lzfseworkflows/Go/badge.svg?branch=master)](https://github.com/aixiansheng/lzfse/actions) [![GoDoc](https://godoc.org/github.com/aixiansheng/lzfse?status.svg)](https://pkg.go.dev/github.com/aixiansheng/lzfse)

> An LZFSE decompressor written in Go

```golang
import (
	"os"
	"gihub.com/aixiansheng/lzfse"
)

func main() {
	inf, err := os.Open("some.lzfse")
	outf, err := os.Create("some.file")
	d := lzfse.NewReader(fh)
	io.Copy(outf, inf)
}
```
