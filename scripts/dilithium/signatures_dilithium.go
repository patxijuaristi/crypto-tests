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

func SignDilithium() bool {
	return false
}

func VerifySignatureDilithium() bool {
	return false
}
