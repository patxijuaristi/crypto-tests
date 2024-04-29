package ecdsa

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"reflect"

	"crypto/test/common"
	"crypto/test/common/hexutil"
	"crypto/test/common/math"
	"crypto/test/crypto"
)

var (
	testmsg     = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testsig     = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey  = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
)

func GenerateKeyECDSA() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func GenerateKeyECDSAWrapper() {
	_, _ = crypto.GenerateKey()
}

func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(d)
}

func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	return crypto.FromECDSA(priv)
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	return crypto.FromECDSAPub(pub)
}

func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	return crypto.UnmarshalPubkey(pub)
}

func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	return crypto.PubkeyToAddress(p)
}

func SaveECDSA(file string, key *ecdsa.PrivateKey) error {
	return crypto.SaveECDSA(file, key)
}

func LoadECDSA(file string) (*ecdsa.PrivateKey, error) {
	return crypto.LoadECDSA(file)
}

func SignECDSA(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	return crypto.Sign(hash, prv)
}

func VerifySignatureECDSA(pubkey, hash, signature []byte) bool {
	return crypto.VerifySignature(pubkey, hash, signature)
}

func TestEcrecover() {
	pubkey, err := crypto.Ecrecover(testmsg, testsig)
	if err != nil {
		fmt.Printf("recover error: %s", err)
	}
	if !bytes.Equal(pubkey, testpubkey) {
		fmt.Printf("pubkey mismatch: want: %x have: %x", testpubkey, pubkey)
	}
}

func TestVerifySignature() {
	sig := testsig[:len(testsig)-1] // remove recovery id
	if !crypto.VerifySignature(testpubkey, testmsg, sig) {
		fmt.Printf("can't verify signature with uncompressed key")
	}
	if !crypto.VerifySignature(testpubkeyc, testmsg, sig) {
		fmt.Printf("can't verify signature with compressed key")
	}

	if crypto.VerifySignature(nil, testmsg, sig) {
		fmt.Printf("signature valid with no key")
	}
	if crypto.VerifySignature(testpubkey, nil, sig) {
		fmt.Printf("signature valid with no message")
	}
	if crypto.VerifySignature(testpubkey, testmsg, nil) {
		fmt.Printf("nil signature valid")
	}
	if crypto.VerifySignature(testpubkey, testmsg, append(common.CopyBytes(sig), 1, 2, 3)) {
		fmt.Printf("signature valid with extra bytes at the end")
	}
	if crypto.VerifySignature(testpubkey, testmsg, sig[:len(sig)-2]) {
		fmt.Printf("signature valid even though it's incomplete")
	}
	wrongkey := common.CopyBytes(testpubkey)
	wrongkey[10]++
	if crypto.VerifySignature(wrongkey, testmsg, sig) {
		fmt.Printf("signature valid with with wrong public key")
	}
}

// This test checks that crypto.VerifySignature rejects malleable signatures with s > N/2.
func TestVerifySignatureMalleable() {
	sig := hexutil.MustDecode("0x638a54215d80a6713c8d523a6adc4e6e73652d859103a36b700851cb0e61b66b8ebfc1a610c57d732ec6e0a8f06a9a7a28df5051ece514702ff9cdff0b11f454")
	key := hexutil.MustDecode("0x03ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd3138")
	msg := hexutil.MustDecode("0xd301ce462d3e639518f482c7f03821fec1e602018630ce621e1e7851c12343a6")
	if crypto.VerifySignature(key, msg, sig) {
		fmt.Printf("crypto.VerifySignature returned true for malleable signature")
	}
}

func TestDecompressPubkey() {
	key, err := crypto.DecompressPubkey(testpubkeyc)
	if err != nil {
		fmt.Printf("Error")
	}
	if uncompressed := crypto.FromECDSAPub(key); !bytes.Equal(uncompressed, testpubkey) {
		fmt.Printf("wrong public key result: got %x, want %x", uncompressed, testpubkey)
	}
	if _, err := crypto.DecompressPubkey(nil); err == nil {
		fmt.Printf("no error for nil pubkey")
	}
	if _, err := crypto.DecompressPubkey(testpubkeyc[:5]); err == nil {
		fmt.Printf("no error for incomplete pubkey")
	}
	if _, err := crypto.DecompressPubkey(append(common.CopyBytes(testpubkeyc), 1, 2, 3)); err == nil {
		fmt.Printf("no error for pubkey with extra bytes at the end")
	}
}

func TestCompressPubkey() {
	key := &ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     math.MustParseBig256("0xe32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a"),
		Y:     math.MustParseBig256("0x0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652"),
	}
	compressed := crypto.CompressPubkey(key)
	if !bytes.Equal(compressed, testpubkeyc) {
		fmt.Printf("wrong public key result: got %x, want %x", compressed, testpubkeyc)
	}
}

func TestPubkeyRandom() {
	const runs = 200

	for i := 0; i < runs; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			fmt.Printf("iteration %d: %v", i, err)
		}
		pubkey2, err := crypto.DecompressPubkey(crypto.CompressPubkey(&key.PublicKey))
		if err != nil {
			fmt.Printf("iteration %d: %v", i, err)
		}
		if !reflect.DeepEqual(key.PublicKey, *pubkey2) {
			fmt.Printf("iteration %d: keys not equal", i)
		}
	}
}

func BenchmarkEcrecoverSignature() {
	for i := 0; i < 10; i++ {
		if _, err := crypto.Ecrecover(testmsg, testsig); err != nil {
			fmt.Printf("ecrecover error")
		}
	}
}

func BenchmarkVerifySignature() {
	sig := testsig[:len(testsig)-1] // remove recovery id
	for i := 0; i < 10; i++ {
		if !crypto.VerifySignature(testpubkey, testmsg, sig) {
			fmt.Printf("verify error")
		}
	}
}

func BenchmarkDecompressPubkey() {
	for i := 0; i < 10; i++ {
		if _, err := crypto.DecompressPubkey(testpubkeyc); err != nil {
			fmt.Printf("Error")
		}
	}
}
