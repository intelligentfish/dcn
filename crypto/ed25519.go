package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"github.com/intelligentfish/dcn/io"
	"io/ioutil"
	"os"
)

// ED25519 ed25519 wrapper
type ED25519 struct {
	pub ed25519.PublicKey
	prv ed25519.PrivateKey
}

// New factory method
func New() (object *ED25519, err error) {
	keyPath := os.ExpandEnv("$HOME/.dcn.key")
	object = &ED25519{}
	if io.PathExists(keyPath) {
		err = object.fromFile(keyPath)
		return
	}
	if object.pub, object.prv, err = ed25519.GenerateKey(rand.Reader); nil == err {
		err = object.toFile(keyPath)
	}
	return
}

// NewFromPrv factory method
func NewFromPrv(prv ed25519.PrivateKey) *ED25519 {
	object := &ED25519{
		prv: prv,
		pub: ed25519.PublicKey(prv[32:]),
	}
	return object
}

// NewFromPub factory method
func NewFromPub(pub ed25519.PublicKey) *ED25519 {
	return &ED25519{
		pub: pub,
	}
}

// fromFile load private key from file
func (object *ED25519) fromFile(path string) (err error) {
	var input []byte
	if input, err = ioutil.ReadFile(path); nil != err {
		return
	}
	output := make([]byte, base64.StdEncoding.DecodedLen(len(input)))
	var n int
	if n, err = base64.StdEncoding.Decode(output, input); nil != err {
		return
	}
	object.prv = output[:n]
	object.pub = []byte(object.prv)[32:]
	return
}

// toFile save private key to file
func (object *ED25519) toFile(path string) (err error) {
	output := make([]byte, base64.StdEncoding.EncodedLen(len(object.prv)))
	base64.StdEncoding.Encode(output, object.prv)
	err = ioutil.WriteFile(path, output, 0666)
	return
}

// Sign sign input bytes
func (object *ED25519) Sign(input []byte) (output []byte) {
	output = ed25519.Sign(object.prv, input)
	return
}

// Verify input bytes use signed bytes
func (object *ED25519) Verify(input, sig []byte) bool {
	return ed25519.Verify(object.pub, input, sig)
}
