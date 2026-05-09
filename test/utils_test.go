package test

import (
	"slices"
	"testing"
)

func TestSlice(t *testing.T) {
	list := []int32{1, 2, 3}

	for i := len(list) - 1; i >= 0; i-- {
		println(list[i])
		list = slices.Delete(list, i, i+1)
	}

}
