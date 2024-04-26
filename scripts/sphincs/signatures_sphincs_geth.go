package sphincs

import (
	"crypto/ecdsa"
	"fmt"

	"crypto/test/common/hexutil"
	"crypto/test/common/math"
	"crypto/test/crypto"

	"github.com/kasperdi/SPHINCSPLUS-golang/parameters"
	"github.com/kasperdi/SPHINCSPLUS-golang/sphincs"
)

// Define test data for SPHINCS+
var (
	testmsg_sphincs2     = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testsig_sphincs2     = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey_sphincs2  = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc_sphincs2 = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
)

func Sign(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(digestHash) != crypto.DigestLength {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", crypto.DigestLength, len(digestHash))
	}
	/*
		seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
		defer zeroBytes(seckey)
		return secp256k1.Sign(digestHash, seckey)*/
	signature := SignWithSphincs(digestHash, prv)
	return signature.R, nil
}

func SignWithSphincs(digestHash []byte, prv *ecdsa.PrivateKey) *sphincs.SPHINCS_SIG {
	params := parameters.MakeSphincsPlusSHA256256fRobust(false)
	seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
	defer zeroBytes(seckey)

	sk, err := sphincs.DeserializeSK(params, seckey)

	if err != nil {
		return sphincs.Spx_sign(params, digestHash, sk)
	} else {
		return nil
	}
}

// VerifySignature checks that the given public key created signature over digest.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature(pubkey, digestHash, signature []byte) bool {
	//return secp256k1.VerifySignature(pubkey, digestHash, signature)
	params := parameters.MakeSphincsPlusSHA256256fRobust(false)
	sig, err := sphincs.DeserializeSignature(params, signature)
	if err != nil {
		pk, err := sphincs.DeserializePK(params, pubkey)
		if err != nil {
			return sphincs.Spx_verify(params, digestHash, sig, pk)
		} else {
			return false
		}
	} else {
		return false
	}
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
