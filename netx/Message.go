package netx

import (
	"encoding/binary"
	"errors"
	"net"
	"strings"
)

const (
	DefaultProtoc  = "json"
	DefaultVersion = "1.0.0"

	DefaultMaxMsgSize = 1024 * 1024 * 16 //16MB
)

type IMessage interface {
	GetProtoc() []byte  //协议
	GetVersion() []byte //版本
	GetCommand() []byte //命令
	GetBody() []byte    //包体

	GetProtocLen() uint8   //协议
	GetVersionLen() uint32 //版本
	GetCommandLen() uint32 //命令
	GetBodyLen() uint32    //包体

	SetProtoc([]byte)  //协议
	SetVersion([]byte) //版本
	SetCommand([]byte) //命令
	SetBody([]byte)    //包体

	String() string //转string
}

type message struct {
	protocLen  uint8  //协议长度
	versionLen uint32 //版本长度
	commandLen uint32 //命令
	bodyLen    uint32 //包体长度

	protoc  []byte //协议
	version []byte //版本
	command []byte //命令
	body    []byte //包体
}

func (m *message) GetProtoc() []byte {
	return m.protoc
}

func (m *message) GetVersion() []byte {
	return m.version
}

func (m *message) GetBody() []byte {
	return m.body
}

func (m *message) GetCommand() []byte {
	return m.command
}

func (m *message) SetProtoc(protoc []byte) {
	m.protoc = protoc
}
func (m *message) SetVersion(version []byte) {
	m.version = version
}
func (m *message) SetCommand(command []byte) {
	m.command = command
}
func (m *message) SetBody(body []byte) {
	m.body = body
}

func (m *message) GetProtocLen() uint8 {
	return uint8(len(m.protoc))
}

func (m *message) GetVersionLen() uint32 {
	return uint32(len(m.version))
}

func (m *message) GetCommandLen() uint32 {
	return uint32(len(m.command))
}

func (m *message) GetBodyLen() uint32 {
	return uint32(len(m.body))
}

func (m *message) String() string {
	builder := strings.Builder{}
	builder.WriteString("protoc:")
	builder.WriteString(string(m.protoc))
	builder.WriteString("\t ")
	builder.WriteString("version:")
	builder.WriteString(string(m.version))
	builder.WriteString("\t ")
	builder.WriteString("command:")
	builder.WriteString(string(m.command))
	builder.WriteString("\t ")
	builder.WriteString("body:")
	builder.WriteString(string(m.body))
	return builder.String()

}
func NewDefaultMessage(command, msg []byte) IMessage {
	message := message{protoc: []byte(DefaultProtoc), version: []byte(DefaultVersion), command: command, body: msg}
	message.protocLen = uint8(len(message.protoc))
	message.versionLen = uint32(len(message.version))
	message.commandLen = uint32(len(message.command))
	message.bodyLen = uint32(len(message.body))
	return &message
}
func NewMessage(protoc, version, command, msg []byte) IMessage {
	message := message{protoc: protoc, version: version, command: command, body: msg}
	message.protocLen = uint8(len(message.protoc))
	message.versionLen = uint32(len(message.version))
	message.commandLen = uint32(len(message.command))
	message.bodyLen = uint32(len(message.body))
	return &message
}

type DataPack struct {
	net.Conn
}

// Pack
//
//	@Description:
//	@receiver p
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-11-01 10:05:56
//	@param msg
//	@return error
func (p *DataPack) Pack(command, msg []byte) error {
	message := message{protoc: []byte(DefaultProtoc), version: []byte(DefaultVersion), command: command, body: msg}
	message.protocLen = uint8(len(message.protoc))
	message.versionLen = uint32(len(message.version))
	message.commandLen = uint32(len(message.command))
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
func (p *DataPack) PackMessage(message IMessage) error {
	if err := binary.Write(p, binary.BigEndian, message.GetProtocLen()); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.GetProtoc()); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.GetVersionLen()); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.GetVersion()); err != nil {
		return err
	}

	if err := binary.Write(p, binary.BigEndian, message.GetCommandLen()); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.GetCommand()); err != nil {
		return err
	}

	if message.GetBodyLen() > DefaultMaxMsgSize {
		return errors.New("too large msg ")
	}
	if err := binary.Write(p, binary.BigEndian, message.GetBodyLen()); err != nil {
		return err
	}
	if err := binary.Write(p, binary.BigEndian, message.GetBody()); err != nil {
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
func (p *DataPack) UnPackMessage() (IMessage, error) {
	message := message{}
	if err := binary.Read(p, binary.BigEndian, &message.protocLen); err != nil {
		return nil, err
	}
	if message.protocLen > 0 {
		message.protoc = make([]byte, message.protocLen)
		if _, err := p.Read(message.protoc); err != nil {
			return nil, err
		}
	}

	if err := binary.Read(p, binary.BigEndian, &message.versionLen); err != nil {
		return nil, err
	}
	message.version = make([]byte, message.versionLen)
	if message.versionLen > 0 {
		if _, err := p.Read(message.version); err != nil {
			return nil, err
		}
	}

	if err := binary.Read(p, binary.BigEndian, &message.commandLen); err != nil {
		return nil, err
	}
	if message.commandLen > 0 {
		message.command = make([]byte, message.commandLen)
		if _, err := p.Read(message.command); err != nil {
			return nil, err
		}
	}

	if err := binary.Read(p, binary.BigEndian, &message.bodyLen); err != nil {
		return nil, err
	}
	if message.bodyLen > DefaultMaxMsgSize {
		return nil, errors.New("too large msg ")
	}
	if message.bodyLen > 0 {
		message.body = make([]byte, message.bodyLen)
		if _, err := p.Read(message.body); err != nil {
			return nil, err
		}
	}
	return &message, nil
}
