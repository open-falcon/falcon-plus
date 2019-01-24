package encrypt

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "encoding/gob"
)


// struct -> byte
func Encode(src interface{}) []byte  {
    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    enc.Encode(src)
    return buf.Bytes()
}

// byte -> struct
func Decode(from []byte, to interface{}) {
    dec := gob.NewDecoder(bytes.NewBuffer(from))
    dec.Decode(to)
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
    padding := blockSize - len(cipherText)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(cipherText, padText...)
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    unPadding := int(origData[length-1])
    return origData[:(length - unPadding)]
}

func Encrypt(origData, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    origData = PKCS5Padding(origData, blockSize)
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    encrypted := make([]byte, len(origData))

    blockMode.CryptBlocks(encrypted, origData)
    return encrypted, nil
}

func Decrypt(encrypted, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    origData := make([]byte, len(encrypted))

    blockMode.CryptBlocks(origData, encrypted)
    origData = PKCS5UnPadding(origData)
    return origData, nil
}
