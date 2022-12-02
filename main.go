package main

import (
	"time"

	"github.com/madokast/distributed_system_learning/message"
	"github.com/madokast/distributed_system_learning/node"
)

func main() {
	n1 := node.New("n1", 20021)

	n1.Send(n1, message.Info("hello"))

	n2 := node.New("n2", 20022)
	n1.Send(n1, message.RecordNode("n2", 20022))
	n1.Send(n2, message.Info("hellon2"))
	n1.Send(n2, message.RecordNode("n1", 20021))

	n2.Send(n1, message.Info("iknowyoun1"))

	// n1.Send(n2, message.Info("newfriend"))

	time.Sleep(time.Second)

	n1.Kill()
	n2.Kill()
}
