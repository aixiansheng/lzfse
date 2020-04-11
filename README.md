# A LZFSE decompressor written in Go

```
import (
	"os"
	"gihub.com/aixiansheng/go-pure-lzfse"
)

inf, err := os.Open("some.lzfse")
outf, err := os.Create("some.file")
d := lzfse.NewReader(fh)
io.Copy(outf, inf)
```
