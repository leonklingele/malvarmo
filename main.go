package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/leonklingele/malvarmo/address"
)

func run(prefix []byte, numWorkers int) error {
	var (
		spendKeyPair, viewKeyPair *address.KeyPair
		addr                      []byte
		err                       error
	)
	if bytes.Equal([]byte{}, prefix) {
		spendKeyPair, viewKeyPair, addr, err = address.New()
	} else {
		spendKeyPair, viewKeyPair, addr, err = address.NewWithPrefix(prefix, numWorkers)
	}
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
	fmt.Println("Address:          ", string(addr))

	return nil
}

func main() {
	prefix := flag.String("prefix", "", "optional, the address prefix to search for")
	numWorkers := flag.Int("workers", runtime.GOMAXPROCS(-1), "optional, the number of workers to use for prefix search")
	showHelp := flag.Bool("help", false, "show help and exit")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if err := run([]byte(*prefix), *numWorkers); err != nil {
		log.Fatal(err)
	}
}
