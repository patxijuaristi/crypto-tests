package sphincs

import (
	"crypto/test/common/hexutil"
	"encoding/hex"
	"fmt"

	"github.com/kasperdi/SPHINCSPLUS-golang/parameters"
	"github.com/kasperdi/SPHINCSPLUS-golang/sphincs"
)

// Define test data for SPHINCS+
var (
	testmsg_sphincs     = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testsig_sphincs     = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey_sphincs  = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc_sphincs = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
	params              = parameters.MakeSphincsPlusSHA256256fSimple(false)
)

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

// TestSphincs is a function to test SPHINCS+ signature scheme
func TestSphincs() {
	fmt.Printf("Test SPHINCS+\n\n")

	// Generate parameters for SPHINCS+
	params := parameters.MakeSphincsPlusSHA256256fRobust(false)

	// Generate key pair
	sk, pk := sphincs.Spx_keygen(params)

	// Generate signature
	signature := sphincs.Spx_sign(params, testmsg_sphincs, sk)

	// Print private and public keys
	fmt.Printf("\nPrivate key (seed): %x\n", sk.SKseed)
	fmt.Printf("Private key: (prf) %x\n", sk.SKprf)
	fmt.Printf("Private key: (pkseed) %x\n", sk.PKseed)
	fmt.Printf("Private key: (pkroot) %x\n", sk.PKroot)
	fmt.Printf("Public key (seed): %x\n", pk.PKseed)
	fmt.Printf("Public key (root): %x\n", pk.PKroot)

	// Verify the signature
	if sphincs.Spx_verify(params, testmsg_sphincs, signature, pk) {
		fmt.Printf("\nSignature verified\n")
	} else {
		fmt.Printf("\nFailure to verify signature\n")
	}

	// Attempt to verify a modified signature (should fail)
	signature.R[0] ^= 1
	if sphincs.Spx_verify(params, testmsg_sphincs, signature, pk) {
		fmt.Printf("\nBad signature has been verified\n")
	} else {
		fmt.Printf("\nBad signature failed to verify\n")
	}

	// Print the modified signature
	Rhex := hex.EncodeToString(signature.R)
	fmt.Printf("R = %s (Length=%d)\n", Rhex, len(Rhex))

	// Print parts of the signature for analysis
	/*fmt.Printf("Showing [SK, Auth]\n")
	for i := 0; i < params.K; i++ {
		fmt.Printf("[%s, %s]\n", hex.EncodeToString(signature.SIG_FORS.GetSK(i)[:16]), hex.EncodeToString(signature.SIG_FORS.GetAUTH(i)[:16]))
	}*/
}
