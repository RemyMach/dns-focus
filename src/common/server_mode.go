package common

type ServerMode int

const (
	Proxy ServerMode = iota
	Server
)