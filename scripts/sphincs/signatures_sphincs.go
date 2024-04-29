package sphincs

import (
	"github.com/kasperdi/SPHINCSPLUS-golang/parameters"
	"github.com/kasperdi/SPHINCSPLUS-golang/sphincs"
)

// Define test data for SPHINCS+
var params = parameters.MakeSphincsPlusSHA256256fRobust(false)
var modeSphincs = "SPHINCS+ SHA256 256bit - Robust"
var nParams = 0

func GetCurrentSphincsMode() string {
	return modeSphincs
}

func ChangeSphincsMode() {
	if nParams == 4 {
		nParams = 0
	} else {
		nParams++
	}

	switch nParams {
	case 0:
		params = parameters.MakeSphincsPlusSHA256256fRobust(false)
		modeSphincs = "SPHINCS+ SHA256 256bit - Robust"
	case 1:
		params = parameters.MakeSphincsPlusSHA256256fSimple(false)
		modeSphincs = "SPHINCS+ SHA256 256bit - Simple"
	case 2:
		params = parameters.MakeSphincsPlusSHA256128fRobust(false)
		modeSphincs = "SPHINCS+ SHA256 128bit - Robust"
	case 3:
		params = parameters.MakeSphincsPlusSHA256128fSimple(false)
		modeSphincs = "SPHINCS+ SHA256 128bit - Simple"
	}
}

func GenerateKeySPHINCS() (*sphincs.SPHINCS_SK, *sphincs.SPHINCS_PK) {
	return sphincs.Spx_keygen(params)
}

func GenerateKeySPHINCSWrapper() {
	_, _ = sphincs.Spx_keygen(params)
}

func SignSPHINCS(hash []byte, sk *sphincs.SPHINCS_SK) *sphincs.SPHINCS_SIG {
	return sphincs.Spx_sign(params, hash, sk)
}

func VerifySignatureSPHINCS(hash []byte, signature *sphincs.SPHINCS_SIG, pk *sphincs.SPHINCS_PK) bool {
	return sphincs.Spx_verify(params, hash, signature, pk)
}

func KeysToBytes(sk *sphincs.SPHINCS_SK, pk *sphincs.SPHINCS_PK) ([]byte, []byte) {
	var sk_as_bytes []byte
	sk_as_bytes = append(sk_as_bytes, sk.SKseed...)
	sk_as_bytes = append(sk_as_bytes, sk.SKprf...)
	sk_as_bytes = append(sk_as_bytes, sk.PKseed...)
	sk_as_bytes = append(sk_as_bytes, sk.PKroot...)

	var pk_as_bytes []byte
	pk_as_bytes = append(pk_as_bytes, pk.PKseed...)
	pk_as_bytes = append(pk_as_bytes, pk.PKroot...)

	return sk_as_bytes, pk_as_bytes
}

func SignatureToBytes(s *sphincs.SPHINCS_SIG) []byte {
	var sig_as_bytes []byte
	sig_as_bytes = s.R
	for i := range s.SIG_FORS.Forspkauth {
		sig_as_bytes = append(sig_as_bytes, s.SIG_FORS.GetSK(i)...)
		sig_as_bytes = append(sig_as_bytes, s.SIG_FORS.GetAUTH(i)...)
	}
	for _, xmssSig := range s.SIG_HT.XMSSSignatures {
		sig_as_bytes = append(sig_as_bytes, xmssSig.GetWOTSSig()...)
		sig_as_bytes = append(sig_as_bytes, xmssSig.GetXMSSAUTH()...)
	}
	return sig_as_bytes
}
