package node

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/madokast/distributed_system_learning/message"
)

type node interface {
	Name() string
	Send(name Namer, msg message.Message) bool
	handle(msg message.Message)
	Log(s string)
	Err(s string)
	Kill()
}

func New(name string, port uint16) node {

	n := &nodeImpl{
		name:       name,
		knownNames: make(map[string]*destination),
		server:     nil,
	}

	s := &http.Server{
		Addr:           fmt.Sprintf("localhost:%d", port),
		Handler:        httpHandler{n},
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go s.ListenAndServe()

	n.server = s
	n.knownNames[name] = &destination{
		Name: name,
		IP:   "localhost",
		Port: port,
	}
	n.Log("Starting")

	return n
}

type httpHandler struct {
	n node
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.n.handle(message.New(r.URL.Query().Get("msg")))
}

type nodeImpl struct {
	name       string
	knownNames map[string]*destination
	server     *http.Server
}

func (n *nodeImpl) Send(name Namer, msg message.Message) bool {
	n.Log(fmt.Sprintf("Sent %s to %s", msg.Content(), name.Name()))

	d, ex := n.knownNames[name.Name()]
	if !ex {
		n.Err(fmt.Sprintf("Unknown name %s whom sending %s", name.Name(), msg.Content()))
		return false
	}
	resp, err := http.Get(fmt.Sprintf("http://%s:%d?msg=%s", d.IP, d.Port, msg.Content()))
	if err != nil {
		n.Err(fmt.Sprintf("Error %s when sending %s to %s", err.Error(), msg.Content(), name.Name()))
		return false
	}
	resp.Body.Close()

	return true
}

func (n *nodeImpl) handle(msg message.Message) {
	ss := strings.Split(msg.Content(), "_")
	switch {
	case len(ss) == 0:
		n.Err("Receive empty message")
	case ss[0] == message.MSG_INFO && len(ss) == 2:
		n.Log(fmt.Sprintf("Receive %s", ss[1]))
	case ss[0] == message.MSG_RECORD_NODE:
		n.Log(fmt.Sprintf("Record node %v", ss[1:]))
		port, _ := strconv.Atoi(ss[2])
		n.knownNames[ss[1]] = &destination{
			Name: ss[1],
			IP:   "localhost",
			Port: uint16(port),
		}
	default:
		n.Err(fmt.Sprintf("Nnknow message %s", msg.Content()))
	}

}

func (n *nodeImpl) Name() string {
	return n.name
}

func (n *nodeImpl) Log(s string) {
	fmt.Printf("[%s] %s\n", n.Name(), s)
}

func (n *nodeImpl) Err(s string) {
	fmt.Errorf("[%s] %s\n", n.Name(), s)
}

func (n *nodeImpl) Kill() {
	fmt.Printf("[%s] Killed\n", n.Name())
	n.server.Close()
}
