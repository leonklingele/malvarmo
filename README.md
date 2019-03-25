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
Private Spend Key: 14cf51a9dfe17dc4752642b3ed224075038cfadd7ac0e6fe4f73ced1b7944209
Public Spend Key:  b7937472c7143cdf19f1b3614456b943446795c97413e1102ec7943c5338baf4
Private View Key:  13fe2eafacba62eca72519e8b02c7127beb8aa34f4e672a629438d925270580e
Public View Key:   ceaae0a32aea0ad93a1cdef1bed2479a0c0dfebd2db92713272112cbb67b45f9
Address:           48abce5GhYXeKN2UeGfNxGCFaRC3Y4u1i3hzaiFkQpiDhwwNUb7g6ZXdLNhGWFXFpzSmT5sy3MtAr4ConUWzjFHnVBz3855
```
