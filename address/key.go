package address

import (
	"crypto/rand"
	"fmt"

	"github.com/agl/ed25519/edwards25519"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

type PrivateKey []byte
type PublicKey []byte

type KeyPair struct {
	priv PrivateKey
	pub  PublicKey
}

// PrivateKey returns the private part of the key pair
func (p *KeyPair) PrivateKey() PrivateKey {
	return p.priv
}

// PublicKey returns the public part of the key pair
func (p *KeyPair) PublicKey() PublicKey {
	return p.pub
}

// newSpendKeyPair generates a new spend key pair
func newSpendKeyPair() (*KeyPair, error) {
	// Generate a new random Ed25519 key
	_, k, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 key pair: %s", err.Error())
	}
	// Cut out private key part
	priv := PrivateKey(reduce(k[:32]))
	// .. and generate an associated public key
	pub := private2Public(priv)
	return &KeyPair{priv, pub}, nil
}

// makeViewKeyPair returns a view key pair based on a private spend key
func makeViewKeyPair(p PrivateKey) *KeyPair {
	// Hash private spend key using Keccak-256
	h := sha3.NewLegacyKeccak256()
	h.Write(p)
	// Important: Reduce to stay in finite field
	priv := reduce(h.Sum(nil))
	// Turn private into public key
	pub := private2Public(priv)
	return &KeyPair{priv, pub}
}

// Based on golang.org/x/crypto/ed25519
// private2Public converts a private key into the associated public key
func private2Public(priv PrivateKey) PublicKey {
	var A edwards25519.ExtendedGroupElement
	var scalar [32]byte
	copy(scalar[:], priv[:])
	edwards25519.GeScalarMultBase(&A, &scalar)

	var pub [32]byte
	A.ToBytes(&pub)
	return pub[:]
}

// reduce ensures we stay in the Ed25519 finite field
func reduce(scalar []byte) []byte {
	var in [64]byte
	copy(in[:32], scalar[:])
	var out [32]byte
	edwards25519.ScReduce(&out, &in)
	return out[:]
}
