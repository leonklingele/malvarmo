package main

import (
	"encoding/hex"
	"testing"
)

// Fixture
const (
	privSpendHex = "75a84a0ec08f795474eb4952b40aec6648ffbad90a5cc4bec3a9964fc6ee1c01"
	pubSpendHex  = "f7b84112e3d36b774bcf01e63218439335171562c0d1b8917897b656cfe9ffad"

	privViewHex = "699443b7ba8a0744b54b5a99b0197f0471c2d19027307fac3315d3b67ede640b"
	pubViewHex  = "b6fa07731fe6a0045dfd323aea7ee85abc4f9729028d28e43e92c3878894b424"
)

func TestPrivateSpend2PublicSpend(t *testing.T) {
	priv := h2b(privSpendHex)
	if got := private2Public(priv); b2h(got) != pubSpendHex {
		t.Fatalf("got incorrect public spend key: %s", got)
	}
}

func TestPrivateSpend2ViewPair(t *testing.T) {
	privSpend := h2b(privSpendHex)
	vk := makeViewKeyPair(PrivateKey(privSpend))
	if got := vk.PrivateKey(); b2h(got) != privViewHex {
		t.Fatalf("got incorrect private view key: %s", got)
	}
	if got := vk.PublicKey(); b2h(got) != pubViewHex {
		t.Fatalf("got incorrect public view key: %s", got)
	}
}

func h2b(h string) []byte {
	dec, _ := hex.DecodeString(h)
	return dec
}

func b2h(b []byte) string {
	return hex.EncodeToString(b)
}
