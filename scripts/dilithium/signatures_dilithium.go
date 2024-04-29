package dilithium

import (
	"github.com/cloudflare/circl/sign/dilithium"
)

var modes = dilithium.ModeNames()

// var modes = [...]string{"Dilithium3-AES", "Dilithium5", "Dilithium5-AES", "Dilithium2", "Dilithium2-AES", "Dilithium3"}
var nMode = 0
var mode = dilithium.ModeByName(modes[nMode])

func GetCurrentDilithiumMode() string {
	return modes[nMode]
}

func ChangeDilithiumMode() {
	if nMode == 5 {
		nMode = 0
	} else {
		nMode++
	}
	mode = dilithium.ModeByName(modes[nMode])
}

func GenerateKeyDilithium() (dilithium.PublicKey, dilithium.PrivateKey, error) {
	return mode.GenerateKey(nil)
}

func GenerateKeyDilithiumWrapper() {
	_, _, _ = mode.GenerateKey(nil)
}

func SignDilithium(sk dilithium.PrivateKey, msg []byte) []byte {
	return mode.Sign(sk, msg)
}

func VerifySignatureDilithium(pk dilithium.PublicKey, msg []byte, signature []byte) bool {
	return mode.Verify(pk, msg, signature)
}
