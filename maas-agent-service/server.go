package main

import (
	_ "github.com/netcracker/qubership-core-maas-agent/maas-agent-service/v2/config"
	service "github.com/netcracker/qubership-core-maas-agent/maas-agent-service/v2/lib"
)

func main() {
	service.RunServer()
}
