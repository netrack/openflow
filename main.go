package main

import (
	ofp "./protocol/ofp13"
	"bytes"
	"fmt"
)

func main() {
	var b bytes.Buffer
	hello := ofp.NewHello()

	e := hello.Write(&b)
	fmt.Println(b.Len(), b.Bytes(), e)

	var h ofp.Hello
	e = h.Read(bytes.NewBuffer([]byte{4, 0, 0, 8, 0, 0, 0, 0}))
	fmt.Println(h, e)
}
