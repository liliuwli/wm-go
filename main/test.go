package main

import (
	"bytes"
	"fmt"
	"github.com/study/utils"
)

func main(){
	/*
	byte1 := []byte("12345");
	byte2 := []byte("222");

	fmt.Printf("byte1 : %s \n", byte1)
	fmt.Printf("byte2 : %s \n", byte2)

	rd := bytes.NewBufferString("Hello World!")
	buf := make([]byte, 6)
	// 读出一部分数据，看看切片有没有变化
	rd.Read(buf)
	// 获取数据切片
	b := rd.Bytes()
	fmt.Printf("%s\n", rd.String()) // World!
	fmt.Printf("%s\n\n", b)
	fmt.Printf("%s\n\n", buf)
	byte1 := []byte("12345")

	byte2,byte3 := utils.BytesSlice(byte1,3)
	fmt.Printf("byte2 : %s \n", byte2)
	fmt.Printf("byte3 : %s \n", byte3)
	*/


	buffer := make([]byte,4096)
	rd := bytes.NewBufferString("1111")
	rd.Read(buffer)
	buffer1,_ := utils.BytesSlice(buffer,4)
	fmt.Printf("%s : %d", buffer1,len(buffer1))
}