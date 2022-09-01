package main

import nex "github.com/PretendoNetwork/nex-go"

func passwordFromPID(pid uint32) (string, uint32) {
	user := getNEXAccountByPID(pid)

	if user == nil {
		return "", nex.Errors.RendezVous.InvalidUsername
	}

	return user["password"].(string), 0
}
