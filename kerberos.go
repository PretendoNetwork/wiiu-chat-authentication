package main

import (
	nex "github.com/PretendoNetwork/nex-go"
)

func generateKerberosTicket(userPID uint32, serverPID uint32, keySize int) ([]byte, int) {
	user := getNEXAccountByPID(userPID)
	if user == nil {
		return []byte{}, 0x80030064 // RendezVous::InvalidUsername
	}

	userPassword := user["password"].(string)
	serverPassword := "password"

	// Create session key and ticket keys
	sessionKey := make([]byte, keySize)

	ticketInfoKey := make([]byte, 16)                   // key for encrypting the internal ticket info. Only used by server. TODO: Make this random!
	userKey := deriveKey(userPID, []byte(userPassword)) // Key for encrypting entire ticket. Used by client and server
	serverKey := deriveKey(serverPID, []byte(serverPassword))
	finalKey := nex.MD5Hash(append(serverKey, ticketInfoKey...))

	//rand.Read(sessionKey) // Create a random session key

	//fmt.Println("Using Session Key: " + hex.EncodeToString(sessionKey))

	////////////////////////////////
	// Build internal ticket info //
	////////////////////////////////

	expiration := nex.NewDateTime(0)
	ticketInfoStream := nex.NewStreamOut(nexServer)

	ticketInfoStream.WriteUInt64LE(expiration.Now())
	ticketInfoStream.WriteUInt32LE(userPID)
	ticketInfoStream.Grow(int64(keySize))
	ticketInfoStream.WriteBytesNext(sessionKey)

	// Encrypt internal ticket info

	ticketInfoEncryption := nex.NewKerberosEncryption(nex.MD5Hash(finalKey))
	encryptedTicketInfo := ticketInfoEncryption.Encrypt(ticketInfoStream.Bytes())

	///////////////////////////////////
	// Build ticket data New Version //
	///////////////////////////////////

	ticketDataStream := nex.NewStreamOut(nexServer)

	ticketDataStream.WriteBuffer(ticketInfoKey)
	ticketDataStream.WriteBuffer(encryptedTicketInfo)

	///////////////////////////
	// Build Kerberos Ticket //
	///////////////////////////

	ticketStream := nex.NewStreamOut(nexServer)

	// Write session key
	ticketStream.Grow(int64(keySize))
	ticketStream.WriteBytesNext(sessionKey)
	ticketStream.WriteUInt32LE(serverPID)
	ticketStream.WriteBuffer(ticketDataStream.Bytes())

	// Encrypt the ticket
	ticketEncryption := nex.NewKerberosEncryption(userKey)
	encryptedTicket := ticketEncryption.Encrypt(ticketStream.Bytes())

	return encryptedTicket, 0
}

func deriveKey(pid uint32, password []byte) []byte {
	for i := 0; i < 65000+int(pid)%1024; i++ {
		password = nex.MD5Hash(password)
	}

	return password
}
