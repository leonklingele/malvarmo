package address

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

type fixture struct {
	privSpendHex, pubSpendHex,
	privViewHex, pubViewHex,
	address string
}

var (
	fixtures = []fixture{
		{
			"4a078e76cd41a3d3b534b83dc6f2ea2de500b653ca82273b7bfad8045d85a400",
			"7849297236cd7c0d6c69a3c8c179c038d3c1c434735741bb3c8995c3c9d6f2ac",
			"e514321d6163c9c222f22eb9f43dd1421aee455bb87adb9e0aee138aa8b4b806",
			"c3fb70733f47f076a70766bfc3ff074e7b7c2663e65394790cc214549458d28e",
			"46BVM4CnrP53FE2gcT3LJjAWJ6fGWq8t8YKRqwwit8vmVu3TJhqmYeKLr5VaNKENaJE8Nt1kdzpeFMFLS6aaePC5H35CgTN",
		},
		{
			"75a84a0ec08f795474eb4952b40aec6648ffbad90a5cc4bec3a9964fc6ee1c01",
			"f7b84112e3d36b774bcf01e63218439335171562c0d1b8917897b656cfe9ffad",
			"699443b7ba8a0744b54b5a99b0197f0471c2d19027307fac3315d3b67ede640b",
			"b6fa07731fe6a0045dfd323aea7ee85abc4f9729028d28e43e92c3878894b424",
			"4B1ahC2k4bcLxKfnBpViWWRd6a1VjeFdqRLFbEHhoFciW4FYNUAT45D1jNGq4YKejHGBFSE2ktZRqfBFu3tHaLGT5Aug3jk",
		},
	}
)

func TestPrivateSpend2PublicSpend(t *testing.T) {
	foreachFixture(func(fx fixture) error {
		priv := h2b(fx.privSpendHex)
		if got := private2Public(priv); b2h(got) != fx.pubSpendHex {
			return fmt.Errorf("got incorrect public spend key: %s", got)
		}
		return nil
	}, t)
}

func TestPrivateSpend2ViewPair(t *testing.T) {
	foreachFixture(func(fx fixture) error {
		privSpend := h2b(fx.privSpendHex)
		vk := makeViewKeyPair(PrivateKey(privSpend))
		if got := vk.PrivateKey(); b2h(got) != fx.privViewHex {
			return fmt.Errorf("got incorrect private view key: %s", got)
		}
		if got := vk.PublicKey(); b2h(got) != fx.pubViewHex {
			return fmt.Errorf("got incorrect public view key: %s", got)
		}
		return nil
	}, t)
}

func TestMakeAddress(t *testing.T) {
	foreachFixture(func(fx fixture) error {
		if got := makeAddress(h2b(fx.pubSpendHex), h2b(fx.pubViewHex)); string(got) != fx.address {
			return fmt.Errorf("got incorrect address: %s", got)
		}
		return nil
	}, t)
}

func TestNewAddress(t *testing.T) {
	spendKeyPair, viewKeyPair, _, err := New()
	if err != nil {
		t.Fatalf("failed to create new address: %s", err.Error())
	}

	if got := private2Public(spendKeyPair.PrivateKey()); !bytes.Equal(got, spendKeyPair.PublicKey()) {
		t.Fatalf("got incorrect public spend key from secret spend key: %s", b2h(got))
	}
	if got := private2Public(viewKeyPair.PrivateKey()); !bytes.Equal(got, viewKeyPair.PublicKey()) {
		t.Fatalf("got incorrect public view key from secret view key: %s", b2h(got))
	}

	vk := makeViewKeyPair(spendKeyPair.PrivateKey())
	if got := vk.PrivateKey(); !bytes.Equal(got, viewKeyPair.PrivateKey()) {
		t.Fatalf("got incorrect private view key: %s", b2h(got))
	}
	if got := vk.PublicKey(); !bytes.Equal(got, viewKeyPair.PublicKey()) {
		t.Fatalf("got incorrect public view key: %s", b2h(got))
	}
}

func foreachFixture(f func(fixture) error, t *testing.T) {
	for i, fx := range fixtures {
		if err := f(fx); err != nil {
			t.Fatalf("failed at fixture %d: %s", i, err.Error())
		}
	}
}

func h2b(h string) []byte {
	dec, _ := hex.DecodeString(h)
	return dec
}

func b2h(b []byte) string {
	return hex.EncodeToString(b)
}
