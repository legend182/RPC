package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	// decode与marshal的区别，decode直接从流中进行解码，marshal是从内存中读取的，所以decode用于http连接
	// 与socket连接中读取与写入，或者文件读取，unmarshal用于byte的输入
	dec *gob.Decoder
	enc *gob.Encoder
}

// 此时如果GobCodec没有实现Codec就会报错
var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		// 创建一个从conn中读取的解码器
		dec: gob.NewDecoder(conn),
		// 创建一个编码到buf的编码器
		enc: gob.NewEncoder(buf),
	}
}
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc codec gob error encoding header", err)
		return err
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc codec gob error encoding body", err)
		return err
	}
	return nil
}
func (c *GobCodec) Close() error {
	return c.conn.Close()
}
