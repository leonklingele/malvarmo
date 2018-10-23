package address

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"

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

// nextSpendKeyPairMaker returns a func to generate
// a new key pair using an already existing one.
// The previous key pair will be overwritten.
func nextSpendKeyPairMaker(p *KeyPair) func() {
	one := big.NewInt(1)
	bn := big.NewInt(0)
	return func() {
		// TODO(leon): Reusing the old private key will give us a huge
		// speed increase: p.pub = ed25519GeScalarMult(oldPriv, 2)

		// Increase private key by one
		bn.SetBytes(p.priv)
		_ = bn.Add(bn, one)

		var newPriv [32]byte
		// bn.Bytes() might return a slice of length < 32 bytes
		copy(newPriv[:], bn.Bytes())
		newPriv[0] &= 248
		newPriv[31] &= 127
		newPriv[31] |= 64

		p.priv = reduce(newPriv[:])
		// TODO(leon): Reusing the old private key will give us a huge
		// speed increase: p.pub = ed25519GeScalarMult(oldPriv, 2)
		p.pub = private2Public(p.priv)
	}
}

// nextSpendKeyPair generates a new key pair using
// an already existing one. The previous key pair
// will be overwritten.
func nextSpendKeyPair(p *KeyPair) {
	priv := p.priv
	privUint64 := binary.BigEndian.Uint64(priv)
	privUint64++

	newPriv := priv
	binary.BigEndian.PutUint64(newPriv[:32], privUint64)
	newPriv[0] &= 248
	newPriv[31] &= 127
	newPriv[31] |= 64

	p.priv = reduce(newPriv)
	p.pub = private2Public(p.priv)
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
