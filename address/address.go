package address

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

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

func New() (*KeyPair, *KeyPair, []byte, error) {
	spendKeyPair, err := newSpendKeyPair()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create new spend key pair: %s", err.Error())
	}
	viewKeyPair := makeViewKeyPair(spendKeyPair.PrivateKey())
	address := makeAddress(spendKeyPair.PublicKey(), viewKeyPair.PublicKey())
	return spendKeyPair, viewKeyPair, address, nil
}
