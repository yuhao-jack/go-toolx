package net

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	DefaultProtoc  = "json"
	DefaultVersion = "1.0.0"

	DefaultMaxMsgSize = 1024 * 1024 * 16 //16MB
)

type message struct {
	protocLen  uint8  //协议长度
	versionLen uint32 //版本长度
	bodyLen    uint32 //包体长度

	protoc  []byte //协议
	version []byte //版本
	body    []byte //包体
}

func NewDefaultMessage(msg []byte) *message {
	message := message{protoc: []byte(DefaultProtoc), version: []byte(DefaultVersion), body: msg}
	message.protocLen = uint8(len(message.protoc))
	message.versionLen = uint32(len(message.version))
	message.bodyLen = uint32(len(message.body))
	return &message
}
func NewMessage(protoc, version, msg []byte) {
	message := message{protoc: protoc, version: version, body: msg}
	message.protocLen = uint8(len(message.protoc))
	message.versionLen = uint32(len(message.version))
	message.bodyLen = uint32(len(message.body))
}

type DataPack struct {
	io.ReadWriter
}

// Pack
//
//	@Description:
//	@receiver p
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-01 10:05:56
//	@param msg
//	@return error
func (p *DataPack) Pack(msg []byte) error {
	message := message{protoc: []byte(DefaultProtoc), version: []byte(DefaultVersion), body: msg}
	message.protocLen = uint8(len(message.protoc))
	message.versionLen = uint32(len(message.version))
	message.bodyLen = uint32(len(message.body))

	return p.PackMessage(&message)
}

// PackMessage
//
//	@Description:
//	@receiver p
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-01 10:10:17
//	@param message
//	@return error
func (p *DataPack) PackMessage(message *message) error {
	if err := binary.Write(p, binary.BigEndian, message.protocLen); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.protoc); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.versionLen); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.version); err != nil {
		return err
	}
	if message.bodyLen > DefaultMaxMsgSize {
		return errors.New("too large msg ")
	}
	if err := binary.Write(p, binary.BigEndian, message.bodyLen); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.body); err != nil {
		return err
	}
	return nil
}

// UnPackMessage
//
//	@Description:
//	@receiver p
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-01 12:32:00
//	@return *message
//	@return error
func (p *DataPack) UnPackMessage() (*message, error) {
	message := message{}
	if err := binary.Read(p, binary.BigEndian, &message.protocLen); err != nil {
		return nil, err
	}
	message.protoc = make([]byte, message.protocLen)
	if _, err := p.Read(message.protoc); err != nil {
		return nil, err
	}
	if err := binary.Read(p, binary.BigEndian, &message.versionLen); err != nil {
		return nil, err
	}
	message.version = make([]byte, message.versionLen)
	if _, err := p.Read(message.version); err != nil {
		return nil, err
	}

	if err := binary.Read(p, binary.BigEndian, &message.bodyLen); err != nil {
		return nil, err
	}
	if message.bodyLen > DefaultMaxMsgSize {
		return nil, errors.New("too large msg ")
	}
	message.body = make([]byte, message.bodyLen)
	if _, err := p.Read(message.body); err != nil {
		return nil, err
	}
	return &message, nil
}
