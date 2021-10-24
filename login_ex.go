package main

import (
	"fmt"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func loginEx(err error, client *nex.Client, callID uint32, username string, authenticationInfo *nexproto.AuthenticationInfo) {
	// TODO: Verify auth info

	if err != nil {
		fmt.Println(err)
		return
	}

	userPID, _ := strconv.Atoi(username)

	serverPID := 1 // Quazal Rendez-Vous

	encryptedTicket, errorCode := generateKerberosTicket(uint32(userPID), uint32(serverPID), nexServer.KerberosKeySize())

	if errorCode != 0 {
		fmt.Println(errorCode)
		return
	}

	// Build the response body
	stationURL := "prudps:/address=66.177.0.8;port=60005;CID=1;PID=2;sid=1;stream=10;type=2"
	serverName := "Pretendo WiiU Chat"

	rvConnectionData := nex.NewRVConnectionData()
	rvConnectionData.SetStationURL(stationURL)
	rvConnectionData.SetSpecialProtocols([]byte{})
	rvConnectionData.SetStationURLSpecialProtocols("")
	serverTime := nex.NewDateTime(0)
	rvConnectionData.SetTime(serverTime.Now())

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt32LE(0x10001) // success
	rmcResponseStream.WriteUInt32LE(uint32(userPID))
	rmcResponseStream.WriteBuffer(encryptedTicket)
	rmcResponseStream.WriteStructure(rvConnectionData)
	rmcResponseStream.WriteString(serverName)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.AuthenticationProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.AuthenticationMethodLoginEx, rmcResponseBody)

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
