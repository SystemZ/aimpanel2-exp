package lib

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/cavaliercoder/grab"
	log "github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"io/ioutil"
	"math/rand"
	"net"
	goHttp "net/http"
	"time"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func RandomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}

	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func MustEncode(enc *json.Encoder, v interface{}) {
	err := enc.Encode(v)
	if err != nil {
		log.Printf("ecode: %v", ecode.JsonEncode)
	}
}

type Error struct {
	ErrorCode int
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error Code: %d", e.ErrorCode)
}

func DownloadFile(url string, filepath string) error {
	_, err := grab.Get(filepath, url)
	if err != nil {
		return err
	}

	return nil
}

func CopyFile(source string, destination string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destination, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReturnError(w goHttp.ResponseWriter, httpCode int, errorCode int, err error) {
	if err != nil {
		log.Warn(err)
	}
	w.WriteHeader(httpCode)
	MustEncode(json.NewEncoder(w), response.JsonError{ErrorCode: errorCode})
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	err = listener.Close()
	if err != nil {
		return 0, err
	}
	return listener.Addr().(*net.TCPAddr).Port, nil
}

//https://gist.github.com/ukautz/cd118e298bbd8f0a88fc
func ParseCertificate(certificate string, privateKey string) (*tls.Certificate, error) {
	var cert tls.Certificate

	block, _ := pem.Decode([]byte(certificate))
	cert.Certificate = append(cert.Certificate, block.Bytes)

	block, _ = pem.Decode([]byte(privateKey))
	pk, err := ParsePrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	cert.PrivateKey = pk

	return &cert, nil
}

func ParsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("failed to parse private key")
}
