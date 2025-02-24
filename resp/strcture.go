package resp

import (
	"fmt"
	"strconv"
)

var (
	CRLF = "\r\n"
	NIL  = []byte("$-1\r\n")
)

type RedisData interface {
	ToBytes() []byte //返回可供解析的标准格式
	ByteData() []byte
	String() string
}

type StringData struct { //+
	data string
}

type BulkData struct { //$
	data []byte
}

type IntData struct { //:
	data int64
}

type ErrorData struct { //-
	data string
}

type ArrayData struct { //*
	data []RedisData
}

type PlainData struct {
	data string
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////StringData
func MakeStringData(data string) *StringData {
	return &StringData{
		data: data,
	}
}

func (r *StringData) Data() string {
	return r.data
}

func (r *StringData) ToBytes() []byte {
	return []byte("+" + r.data + CRLF)
}

func (r *StringData) ByteData() []byte {
	return []byte(r.data)
}

func (r *StringData) String() string {
	return r.data
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////BulkData
func MakeBulkData(data []byte) *BulkData {
	return &BulkData{
		data: data,
	}
}

func (r *BulkData) Data() []byte {
	return r.data
}

func (r *BulkData) ToBytes() []byte {
	if r.data == nil {
		return NIL
	}
	return []byte("$" + strconv.Itoa(len(r.data)) + CRLF + string(r.data) + CRLF)
}

func (r *BulkData) ByteData() []byte {
	return r.data
}

func (r *BulkData) String() string {
	return string(r.data)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////ErrorData
func MakeErrorData(data ...string) *ErrorData {
	errMsg := ""
	for _, v := range data {
		errMsg += v
	}
	return &ErrorData{
		data: errMsg,
	}
}

func (r *ErrorData) Data() string {
	return r.data
}

func (r *ErrorData) ToBytes() []byte {
	return []byte("-" + r.data + CRLF)
}

func (r *ErrorData) ByteData() []byte {
	return []byte(r.data)
}

func (r *ErrorData) String() string {
	return r.data
}

func MakeWrongNumberArgs(name string) *ErrorData {
	return &ErrorData{data: fmt.Sprintf("Ero wrong number of arguments for %s command", name)}
}

func MakeWrongType() *ErrorData {
	return &ErrorData{data: "Wrong type operation aginst a key holding the wrong kind of value"}
}
