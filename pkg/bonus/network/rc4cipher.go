package network

import "crypto/rc4"

type Rc4Cipher struct {
	cipher *rc4.Cipher
}

func (slf *Rc4Cipher) Encrypt(data []byte) []byte {
	dst := make([]byte, len(data))
	//dst := data
	slf.cipher.XORKeyStream(dst, data)
	return dst
}

func (slf *Rc4Cipher) Decrypt(data []byte) []byte {
	dst := make([]byte, len(data))
	//dst := data
	slf.cipher.XORKeyStream(dst, data)
	return dst
}

func NewRc4Cipher(key []byte) *Rc4Cipher {
	c, _ := rc4.NewCipher(key)
	rc4 := &Rc4Cipher{
		cipher: c,
	}
	return rc4
}
