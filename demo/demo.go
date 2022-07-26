// @author cold bin
// @date 2022/7/26

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func main() {

	gZipBytes := GZipBytes([]byte("你好，我来自外ssssssssssssssssssssssaa\n" +
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa面..."))
	fmt.Println(gZipBytes, len(gZipBytes))

	ugZipBytes := UGZipBytes(gZipBytes)

	fmt.Println(string(ugZipBytes), len(ugZipBytes))
}

//压缩
func GZipBytes(data []byte) []byte {
	var input bytes.Buffer
	g, _ := gzip.NewWriterLevel(&input, gzip.BestCompression)
	//g := gzip.NewWriter(&input) //面向api编程调用压缩算法的一个api
	//参数就是指向某个数据缓冲区默认压缩等级是DefaultCompression 在这里还有另一个api可以调用调整压缩级别
	//gzip.NewWirterLevel(&in,gzip.BestCompression) NoCompression（对应的int 0）、
	//BestSpeed（1）、DefaultCompression（-1）、HuffmanOnly（-2）BestCompression（9）这几个级别也可以
	//这样写gzip.NewWirterLevel(&in,0)
	//这里的异常返回最好还是处理下，我这里纯属省事
	g.Write(data)
	g.Close()
	return input.Bytes()
}

//解压
func UGZipBytes(data []byte) []byte {
	var out bytes.Buffer
	var in bytes.Buffer
	in.Write(data)
	r, _ := gzip.NewReader(&in)
	r.Close() //这句放在后面也没有问题，不写也没有任何报错
	//机翻注释：关闭关闭读者。它不会关闭底层的io.Reader。为了验证GZIP校验和，读取器必须完全使用，直到io.EOF。

	io.Copy(&out, r) //这里我看了下源码不是太明白，
	//我个人想法是这样的，Reader本身就是go中表示一个压缩文件的形式，r转化为[]byte就是一个符合压缩文件协议的压缩文件

	return out.Bytes()
}
