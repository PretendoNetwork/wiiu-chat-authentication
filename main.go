package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(2)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetAccessKey("e7a47214")

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==WiiU Chat - Auth==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	authenticationServer := nexproto.NewAuthenticationProtocol(nexServer)

	authenticationServer.LoginEx(loginEx)
	authenticationServer.RequestTicket(requestTicket)

	nexServer.Listen(":60004")
}
