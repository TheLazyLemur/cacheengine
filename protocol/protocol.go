package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command byte

const (
	CMDNone Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
	CmdAll
)

type CommandJoin struct{}

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

type CommandGet struct {
	Key []byte
}

type CommandDel struct {
	Key []byte
}

type CommandAll struct{}

func (c *CommandAll) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, CmdAll)

	return buf.Bytes()
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, CmdSet)

	keyLen := int32(len(c.Key))
	_ = binary.Write(buf, binary.LittleEndian, keyLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Key)

	valLen := int32(len(c.Value))
	_ = binary.Write(buf, binary.LittleEndian, valLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Value)

	_ = binary.Write(buf, binary.LittleEndian, int32(c.TTL))

	return buf.Bytes()
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, CmdGet)

	keyLen := int32(len(c.Key))
	_ = binary.Write(buf, binary.LittleEndian, keyLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func (c *CommandDel) Bytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, CmdDel)

	keyLen := int32(len(c.Key))
	_ = binary.Write(buf, binary.LittleEndian, keyLen)
	_ = binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case CmdSet:
		return parseSetCommand(r), nil
	case CmdGet:
		return parseGetCommand(r), nil
	case CmdDel:
		return parseDelCommand(r), nil
	case CmdJoin:
		return &CommandJoin{}, nil
	case CmdAll:
		return parseAllCommand(r), nil
	default:
		return nil, fmt.Errorf("unknown command: %d", cmd)
	}
}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	_ = binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	_ = binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)

	return cmd
}

func parseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Key)

	return cmd
}

func parseDelCommand(r io.Reader) *CommandDel {
	cmd := &CommandDel{}

	var keyLen int32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, &cmd.Key)

	return cmd
}

func parseAllCommand(_ io.Reader) *CommandAll {
	cmd := &CommandAll{}
	return cmd
}
