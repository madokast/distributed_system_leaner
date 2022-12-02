package message

import "fmt"

var MSG_INFO string = "info"
var MSG_RECORD_NODE string = "recordnode"

type Message interface {
	Content() string
}

func New(msg string) Message {
	return &strMsg{msg}
}

func Info(msg string) Message {
	return &strMsg{MSG_INFO + "_" + msg}
}

func RecordNode(name string, port uint16) Message {
	return &strMsg{MSG_RECORD_NODE + "_" + name + "_" + fmt.Sprintf("%d", port)}
}

type strMsg struct {
	Msg string
}

func (m *strMsg) Content() string {
	return m.Msg
}
