package codec

import "io"

//一个典型的rpc调用 err = client.Call("Arith.Multiply", args, &reply)

// Header header信息
type Header struct {
	ServiceMethod string //Arith.Multiply
	Seq           uint64 // 请求序号
	Error         string // 服务端返回的错误信息
}

// Codec 抽象对消息体的编码和解码
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// 仿照工厂设计模式，通过类型返回构造函数

type NewCodecFunc func(reader io.ReadWriteCloser) Codec
type Type string

const (
	GobType Type = "application/gob"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
