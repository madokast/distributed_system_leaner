package node

type Namer interface {
	Name() string
}

type destination struct {
	Name string
	IP   string
	Port uint16
}
