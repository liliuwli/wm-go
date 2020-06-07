package utils

import (
	"bytes"
	"strings"
	"unsafe"
)

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// example data []byte("12345") len 3
// return  []byte("123") []byte("45")
func BytesSlice(data []byte,len int) ([]byte,[]byte) {

	rd := bytes.NewBuffer(data)
	buf := make([]byte,len)

	rd.Read(buf)

	return  buf,rd.Bytes()
}

func BytesFind(data []byte,findstr string) int {
	data_str := Self_bytes2str(data)
	return strings.IndexAny(data_str,findstr)
}

func Self_bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Self_str2bytes(s string) []byte {
  	x := (*[2]uintptr)(unsafe.Pointer(&s))
  	h := [3]uintptr{x[0], x[1], x[1]}
  	return *(*[]byte)(unsafe.Pointer(&h))
}

func ArrMinKeyFind(arr []int) int{
	minVal := arr[0]
	minIndex := 0

	for i:= 1;i<len(arr);i++{
		if minVal > arr[i] {
			minVal = arr[i]
			minIndex = i
		}
	}
	return minIndex
}