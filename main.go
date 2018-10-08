package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
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

// makeAddress returns the address based on the public spend key and the public view key
func makeAddress(pubSpend, pubView PublicKey) []byte {
	// A Monero address 'mAddr' looks as follows:
	// c = netBytePrefix(0x12) | publicSpendKey | publicViewKey
	// mAddr = base58encode(c | checksum(c)[:4])
	const netBytePrefix = byte(18)
	buf := make([]byte, 0, 69)
	buf = append(buf, netBytePrefix)
	buf = append(buf, pubSpend...)
	buf = append(buf, pubView...)
	h := sha3.NewLegacyKeccak256()
	h.Write(buf)
	hash := h.Sum(nil)
	buf = append(buf, hash[:4]...)
	return base58encode(buf)
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

// Based on https://github.com/moneromooo-monero/monero-wallet-generator/blob/master/monero-wallet-generator.html
// base58encode converts data into Base58-format
func base58encode(data []byte) []byte {
	const (
		fullBlockSize        = 8
		fullEncodedBlockSize = 11
	)
	var (
		alphabet          = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
		encodedBlockSizes = []int{0, 2, 3, 5, 6, 7, 9, 10, 11}
	)

	encodeBlock := func(data, buf []byte, index int) []byte {
		lenAlphabet := big.NewInt(int64(len(alphabet)))
		num := big.NewInt(0).SetBytes(data)
		rem := big.NewInt(0)

		i := encodedBlockSizes[len(data)] - 1
		for num.Cmp(big.NewInt(0)) == 1 {
			num.QuoRem(num, lenAlphabet, rem)
			buf[index+i] = alphabet[int(rem.Int64())]
			i--
		}

		return buf
	}

	fullBlockCount := len(data) / fullBlockSize
	lastBlockSize := len(data) % fullBlockSize
	resSize := fullBlockCount*fullEncodedBlockSize + encodedBlockSizes[lastBlockSize]
	res := make([]byte, resSize)
	for i := 0; i < resSize; i++ {
		res[i] = alphabet[0]
	}

	for i := 0; i < fullBlockCount; i++ {
		res = encodeBlock(
			data[i*fullBlockSize:i*fullBlockSize+fullBlockSize],
			res,
			i*fullEncodedBlockSize,
		)
	}

	if lastBlockSize > 0 {
		res = encodeBlock(
			data[fullBlockCount*fullBlockSize:fullBlockCount*fullBlockSize+lastBlockSize],
			res,
			fullBlockCount*fullEncodedBlockSize,
		)
	}

	return res
}

func NewAddress() (*KeyPair, *KeyPair, []byte, error) {
	spendKeyPair, err := newSpendKeyPair()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create new spend key pair: %s", err.Error())
	}
	viewKeyPair := makeViewKeyPair(spendKeyPair.PrivateKey())
	address := makeAddress(spendKeyPair.PublicKey(), viewKeyPair.PublicKey())
	return spendKeyPair, viewKeyPair, address, nil
}

func run() error {
	spendKeyPair, viewKeyPair, address, err := NewAddress()
	if err != nil {
		return fmt.Errorf("failed to create new address: %s", err.Error())
	}

	/*
		Example output:

		Private Spend Key: dbcdb72ac43e2f3f9ca35c0b8fa8cee99759fce9e8d4fe84423186c39bb7260b
		Public Spend Key:  85b84a94d9d7152660c28afffb03c8707e45277c950b24275f2b19db04d4f737
		Private View Key:  6a5c667c9afd0b3256d9090b5aabbf83e592fc717d892ddf7df8275bb7a78400
		Public View Key:   634e9804e703a9c7d05a6a1fc6dd17b45b60e14774140d1a1c710e1be0ccd120
		Address:           46h1w3Z26Va7RKEY5SwD2XKpKsYQY7Qq97axQf2B3b8AAGLGUXr2FRAaRSok3pRHhQXAgvUcsvwJL5NK17egUqyS4euNvSp
	*/
	fmt.Println("Private Spend Key:", hex.EncodeToString(spendKeyPair.PrivateKey()))
	fmt.Println("Public Spend Key: ", hex.EncodeToString(spendKeyPair.PublicKey()))
	fmt.Println("Private View Key: ", hex.EncodeToString(viewKeyPair.PrivateKey()))
	fmt.Println("Public View Key:  ", hex.EncodeToString(viewKeyPair.PublicKey()))
	fmt.Println("Address:          ", string(address))

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
