package zbar

import (
	"testing"
	"fmt"
)

func TestZBarImageCreate(t *testing.T) {
	var img = ZBarImageCreate()
	fmt.Println(ZBarImageGetWidth(img))
	fmt.Println(ZBarImageGetHeight(img))
}