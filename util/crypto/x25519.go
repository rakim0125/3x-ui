package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
)

// GenerateRealityX25519KeyPair returns a REALITY / Xray-compatible X25519 key pair
// using base64.RawURLEncoding, matching `xray x25519` output (PrivateKey / public material).
func GenerateRealityX25519KeyPair() (privateKeyRawURL, publicKeyRawURL string, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", "", err
	}
	raw[0] &= 248
	raw[31] &= 127
	raw[31] |= 64

	key, err := ecdh.X25519().NewPrivateKey(raw)
	if err != nil {
		return "", "", err
	}
	pub := key.PublicKey().Bytes()
	enc := base64.RawURLEncoding.EncodeToString
	return enc(raw), enc(pub), nil
}
