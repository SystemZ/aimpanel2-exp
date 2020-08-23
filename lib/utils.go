package lib

import (
	"encoding/json"
	"fmt"
	"github.com/cavaliercoder/grab"
	log "github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"io/ioutil"
	"math/rand"
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
