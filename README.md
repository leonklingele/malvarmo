# Malvarmo — A Monero cold storage wallet generator

[![Build Status](https://travis-ci.org/leonklingele/malvarmo.svg?branch=master)](https://travis-ci.org/leonklingele/malvarmo)

Malvarmo — meaning "cold" in Esperanto — is a tiny cold storage wallet generator written in Go. It is intended for educational purposes only. __Please don't use it to store any Moneroj.__

## Installation

```sh
$ go get -u github.com/leonklingele/malvarmo
```

## Usage

```sh
$ malvarmo
Private Spend Key: dbcdb72ac43e2f3f9ca35c0b8fa8cee99759fce9e8d4fe84423186c39bb7260b
Public Spend Key:  85b84a94d9d7152660c28afffb03c8707e45277c950b24275f2b19db04d4f737
Private View Key:  6a5c667c9afd0b3256d9090b5aabbf83e592fc717d892ddf7df8275bb7a78400
Public View Key:   634e9804e703a9c7d05a6a1fc6dd17b45b60e14774140d1a1c710e1be0ccd120
Address:           46h1w3Z26Va7RKEY5SwD2XKpKsYQY7Qq97axQf2B3b8AAGLGUXr2FRAaRSok3pRHhQXAgvUcsvwJL5NK17egUqyS4euNvSp
```

To specify an address prefix:

```sh
$ malvarmo -prefix abc
Private Spend Key: ecc492d7c487b0f7fd2db401c9e227187e5668d824156edcde515c720f604f08
Public Spend Key:  d1f6d9d6d6614f031935c50a3615de798a276c2a98386453253ec997b14b3570
Private View Key:  9117ef625cacb933d6f6453df5706155676fb4c8886ce89db9d0d21c5102400b
Public View Key:   1e2e616e3ca5ca2084c5f85c5b26fd797a480206034ad7acf6dabcbbbaf2cf53
Address:           49abc4DSxUn1X4RqB26NDTML63dPsGUsDEuchCee3PorKkgg56eakCh6SUMTYwqev8MKVAAcN2VDQVvyGXakCnR4APvavkM
```
