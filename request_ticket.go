package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func requestTicket(err error, client *nex.Client, callID uint32, userPID uint32, serverPID uint32) {
	if err != nil {
		fmt.Println(err)
		return
	}

	encryptedTicket, errorCode := generateKerberosTicket(userPID, serverPID, nexServer.KerberosKeySize())

	if errorCode != 0 {
		fmt.Println(errorCode)
		return
	}

	// Build the response body
	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt32LE(0x10001) // success
	rmcResponseStream.WriteBuffer(encryptedTicket)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.AuthenticationProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.AuthenticationMethodRequestTicket, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
