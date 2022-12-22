package connection_strings

import (
	"bytes"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"strings"
	"unicode"

	"go-Sample3FrlConnect/logger"
)

var ErrUnPadding = errors.New("UnPadding error")

//func main() {
//	cryptoKeyEnc, err := decryptCryptoKey("TestPub")
//	if err != nil {
//		logger.Error(err.Error())
//	}
//
//	log.Println(cryptoKeyEnc)
//
//	key, err := tripleDESECBEncrypt("devpwd", cryptoKeyEnc)
//	if err != nil {
//		logger.Error(err.Error())
//	}
//
//	log.Println(key)
//}

func tripleDESECBEncrypt(data, key string) (string, error) {
	keyByte, err := getKey(key)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	value := []byte(data)

	k1, k2, k3 := get3DESKey(keyByte)

	enc, err := ecbEncrypter(value, k1, k2, k3)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc), nil
}

image.pngfunc tripleDESECBDecrypt(data, key string) (string, error) {
	value, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	//keyByte, err := getKey()
	//if err != nil {
	//	logger.Error(err.Error())
	//	return "", err
	//}

	keyByte, err := getKey(key)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	k1, k2, k3 := get3DESKey(keyByte)

	dec, err := ecbDecrypter(value, k1, k2, k3)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	decClean := strings.TrimFunc(string(dec), func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	return decClean, nil
}

func ecbEncrypter(origData, k1, k2, k3 []byte) ([]byte, error) {
	block, err := des.NewCipher(k1)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()

	origData = pKCS7Padding(origData, bs)

	buf1, err := encrypt(origData, k1)
	if err != nil {
		return nil, err
	}

	buf2, err := decrypt(buf1, k2)
	if err != nil {
		return nil, err
	}

	out, err := encrypt(buf2, k3)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ecbDecrypter(crypted, k1, k2, k3 []byte) ([]byte, error) {
	buf1, err := decrypt(crypted, k3)
	if err != nil {
		return nil, err
	}
	buf2, err := encrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := decrypt(buf2, k1)
	if err != nil {
		return nil, err
	}

	out, err = pKCS7UnPadding(out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func encryptCryptoKey(pubKey string) (string, error) {
	cryptoKey := "TestKey" // value of key to database
	privKey := config["PRIV_KEY"]

	key1, err := tripleDESECBEncrypt(cryptoKey, privKey)
	if err != nil {
		return "", err
	}

	key2, err := base64.StdEncoding.DecodeString(key1)
	if err != nil {
		return "", err
	}

	key3, err := tripleDESECBEncrypt(string(key2), pubKey)
	if err != nil {
		return "", err
	}

	//key3, err := base64.StdEncoding.DecodeString(key2)
	//if err != nil {
	//	return "", err
	//}

	return string(key3), nil
}

func decryptCryptoKey(pubKey string) (string, error) {
	cryptoKey := config["CRYPTO_KEY"]
	privKey := config["PRIV_KEY"]

	//pubKey = base64.StdEncoding.EncodeToString([]byte(pubKey))

	key1, err := tripleDESECBDecrypt(cryptoKey, pubKey)
	if err != nil {
		return "", err
	}

	key2 := base64.StdEncoding.EncodeToString([]byte(key1))

	key3, err := tripleDESECBDecrypt(key2, privKey)
	if err != nil {
		return "", err
	}

	return key3, nil
}

func getKey(key string) ([]byte, error) {
	//keyStr, err := base64.StdEncoding.DecodeString(config["CRYPTO_KEY"])
	//if err != nil {
	//	logger.Error(err.Error())
	//	return []byte(""), err
	//}

	//keyStr, err := base64.StdEncoding.DecodeString(key)
	//if err != nil {
	//	logger.Error(err.Error())
	//	return []byte(""), err
	//}

	keyHex := md5.Sum([]byte(key))
	keyByte := keyHex[:]

	return keyByte, nil
}

func get3DESKey(key []byte) ([]byte, []byte, []byte) {
	var k1, k2, k3 []byte
	if len(key) < 16 {
		k1 = key
		k2 = key
		k3 = key
	} else if len(key) < 24 {
		k1 = key[:8]
		k2 = key[8:16]
		k3 = k1
	} else {
		k1 = key[:8]
		k2 = key[8:16]
		k3 = key[16:]
	}

	return k1, k2, k3
}

func encrypt(origData, key []byte) ([]byte, error) {
	if len(origData) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()

	if len(origData)%bs != 0 {
		return nil, errors.New("wrong padding")
	}

	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

func pKCS7Padding(src []byte, blockSize int) []byte {
	if len(src) == blockSize {
		return src
	}

	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pKCS7UnPadding(src []byte) ([]byte, error) {

	length := len(src)

	if length == 0 {
		return src, ErrUnPadding
	}
	unpadding := int(src[length-1])

	if length > 2 {
		if unpadding != int(src[length-2]) {
			return src, nil
		}
	}

	if length < unpadding {
		return src, ErrUnPadding
	}
	return src[:(length - unpadding)], nil
}
