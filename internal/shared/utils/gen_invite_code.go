package utils

import (
	"crypto/rand"
	"math/big"
	
)

func GenerateInviteCode(n int) *string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, n)

	for i,_:=range code{
		index,err:=rand.Int(rand.Reader,big.NewInt(int64(len(charset))))
		if err!=nil{
			return nil
		}
		code[i]=charset[index.Int64()]
	}

	newCode:=string(code)
	return &newCode

}