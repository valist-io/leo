package p2p

import (
	"github.com/libp2p/go-libp2p-core/crypto"
)

func GenerateKey() (crypto.PrivKey, crypto.PubKey, error) {
	return crypto.GenerateKeyPair(crypto.Ed25519, -1)
}

func EncodeKey(priv crypto.PrivKey) (string, error) {
	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", err
	}
	return crypto.ConfigEncodeKey(data), nil
}

func DecodeKey(enc string) (crypto.PrivKey, error) {
	data, err := crypto.ConfigDecodeKey(enc)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(data)
}
