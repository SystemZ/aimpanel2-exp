package ahttp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
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
	//P
	"F5:2B:15:6B:06:E3:9A:74:45:4D:A9:C3:F8:A2:43:CF:25:0C:0F:F8:C7:75:A9:18:3B:DC:8E:A8:BE:DB:9F:E4",

	//local
	"C0:5D:55:E1:A7:60:5D:EE:48:7A:01:B2:4F:E6:B7:EF:AC:E3:B4:FC:C2:0B:B2:EE:F8:28:60:89:3A:8D:8C:C2",

	//*.lvlup.pro CF
	"68:42:B1:CF:1F:69:78:AD:41:1A:79:D4:0A:FB:77:BC:9C:89:5E:93:2E:4C:D6:5E:B1:CC:06:9F:DB:69:B3:67",
}

var CurrentHost = 0
var Hosts = config.MASTER_URLS

var SlaveVersion = ""

func InitHttpClient() *http.Client {
	client := &http.Client{}
	client.Transport = &http.Transport{
		DialTLSContext: VerifyPinTLSContext,
		DialContext:    DialContext,
	}

	return client
}

// DialContext is used for http requests
func DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return nil, errors.New("http is not allowed here")
}

// VerifyPinTLSContext is used for ssl requests to check if pin is valid
func VerifyPinTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return conn, err
	}

	keyPinValid := false
	connState := conn.ConnectionState()

	for _, peerCert := range connState.PeerCertificates {
		for _, f := range Fingerprints {
			if f == GenerateCertFingerprint(peerCert.Raw) {
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
		url := Hosts[CurrentHost] + path

		logrus.Info("Request to " + url)

		resp, err := HttpClient.Get(url)
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return nil, err
			}

			if strings.HasSuffix(err.Error(), "pin is not valid") {
				return nil, err
			}
		}
		serverUnavailable := isServerUnavailable(resp)
		if !serverUnavailable {
			defer resp.Body.Close()
			if output == nil {
				return resp, nil
			} else {
				return resp, json.NewDecoder(resp.Body).Decode(output)
			}
		}

		if serverUnavailable {
			logrus.Infof("Host %v unavailable. Switching to next one", Hosts[CurrentHost])
			nextHost()
		}

		time.Sleep(b.Duration())
	}
}

func Post(path, token, jsonStr string) (*http.Response, error) {
	url := Hosts[CurrentHost] + path

	logrus.Info("Request to " + url)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	if SlaveVersion != "" {
		req.Header.Set("X-Slave-Version", SlaveVersion)
	}

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

		serverUnavailable := isServerUnavailable(resp)
		if !serverUnavailable {
			defer resp.Body.Close()
			return resp, nil
		}

		if serverUnavailable {
			logrus.Infof("Host %v unavailable. Switching to next one", Hosts[CurrentHost])
			nextHost()
		}

		time.Sleep(b.Duration())
	}
}

func nextHost() {
	CurrentHost = CurrentHost + 1
	if CurrentHost > (len(Hosts) - 1) {
		CurrentHost = 0
	}

	logrus.Infof("Switching host to %v", Hosts[CurrentHost])
}

func isServerUnavailable(resp *http.Response) bool {
	if resp == nil {
		return true
	}
	switch resp.StatusCode {
	case http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusRequestTimeout:
		return true
	default:
		return false
	}
}

func SendTaskData(url string, token string, taskMsg task.Message) (int, error) {
	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		return 0, err
	}

	resp, err := Post(url, token, jsonStr)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func SendTaskBatchData(url string, token string, taskMsg task.Messages) (int, error) {
	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		return 0, err
	}

	resp, err := Post(url, token, jsonStr)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func GenerateCertFingerprint(peerCertRaw []byte) (hash string) {
	// https://stackoverflow.com/a/38065844/1351857
	hashRaw := sha256.Sum256(peerCertRaw)
	hash = strings.ToUpper(hex.EncodeToString(hashRaw[:]))
	// https://stackoverflow.com/questions/33633168/how-to-insert-a-character-every-x-characters-in-a-string-in-golang
	for i := 2; i < len(hash); i += 3 {
		// make sure that output is 1:1 with openssl CLI tool
		hash = hash[:i] + ":" + hash[i:]
	}
	return
}
