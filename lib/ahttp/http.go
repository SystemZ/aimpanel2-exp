package ahttp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var b = &backoff.Backoff{
	Min:    2 * time.Second,
	Max:    1 * time.Minute,
	Factor: 2,
	Jitter: true,
}

var HttpClient *http.Client
var Fingerprints = []string{
	//local
	"74439d64c7d7c6d30fc1fbad056ded5c19674fec425d181a9207b6cb1891ccbd",

	//my-lab.aimpanel.pro
	"67d08017c05e2bca29f404947491b9055aad84da17b45951e3b9c1fa7f5126a4",
}

func InitHttpClient() *http.Client {
	client := &http.Client{}
	client.Transport = &http.Transport{
		DialTLSContext: VerifyPinTLSContext,
	}

	return client
}

func VerifyPinTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return conn, err
	}

	keyPinValid := false
	connState := conn.ConnectionState()

	for _, peerCert := range connState.PeerCertificates {
		der, err := x509.MarshalPKIXPublicKey(peerCert.PublicKey)
		hash := sha256.Sum256(der)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range Fingerprints {
			if f == hex.EncodeToString(hash[:]) {
				keyPinValid = true
			}
		}
	}

	if !keyPinValid {
		return nil, errors.New("pin is not valid")
	}

	return conn, nil
}

func Get(path string, output interface{}) (*http.Response, error) {
	for {
		logrus.Info("Request to " + path)

		resp, err := HttpClient.Get(path)
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return nil, err
			}

			if strings.HasSuffix(err.Error(), "pin is not valid") {
				return nil, err
			}
		}

		if resp != nil && !isServerUnavailable(resp.StatusCode) {
			defer resp.Body.Close()
			if output == nil {
				return resp, nil
			} else {
				return resp, json.NewDecoder(resp.Body).Decode(output)
			}
		}

		time.Sleep(b.Duration())
	}
}

func Post(path, token, jsonStr string) (*http.Response, error) {
	req, err := http.NewRequest("POST", path, bytes.NewBufferString(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	for {
		resp, err := HttpClient.Do(req)
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return nil, err
			}

			if strings.HasSuffix(err.Error(), "pin is not valid") {
				return nil, err
			}
		}

		if resp != nil && !isServerUnavailable(resp.StatusCode) {
			defer resp.Body.Close()
			return resp, nil
		}

		time.Sleep(b.Duration())
	}
}

func isServerUnavailable(code int) bool {
	switch code {
	case http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusRequestTimeout:
		return true
	default:
		return false
	}
}

func SendTaskData(url string, token string, jsonStr string) (int, error) {
	resp, err := Post(url, token, jsonStr)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
