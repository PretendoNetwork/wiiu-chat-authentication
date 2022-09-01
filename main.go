package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-common-go/authentication"
)

var nexServer *nex.Server

func main() {
	/*
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
	*/

	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(2)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	nexServer.SetAccessKey("e7a47214")

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==WiiU Chat - Auth==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	authenticationProtocol := authentication.NewCommonAuthenticationProtocol(nexServer)

	secureStationURL := nex.NewStationURL("")
	secureStationURL.SetScheme("prudps")
	secureStationURL.SetAddress(os.Getenv("SECURE_SERVER_LOCATION"))
	secureStationURL.SetPort(os.Getenv("SECURE_SERVER_PORT"))
	secureStationURL.SetCID("1")
	secureStationURL.SetPID("2")
	secureStationURL.SetSID("1")
	secureStationURL.SetStream("10")
	secureStationURL.SetType("2")

	authenticationProtocol.SetSecureStationURL(secureStationURL)
	authenticationProtocol.SetBuildName("Pretendo WiiU Chat Auth")
	authenticationProtocol.SetPasswordFromPIDFunction(passwordFromPID)

	nexServer.Listen(":60004")
}
