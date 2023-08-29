package domain_scan

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	UserAgent          string
	IsUint32           bool
	InsecureSkipVerify bool
	ReqTimeOut         int = 10
)

func FromURLGetContent(requrl string) (content []byte, err error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(ReqTimeOut),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify}, // param
		},
	}
	req, err := http.NewRequest("GET", requrl, nil) //nolint
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := MyError{"favicon not exist"}
		return nil, err
	}
	defer resp.Body.Close() //nolint
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetFaviconMmh3Hash(url string) string {
	faviconUrl := fmt.Sprintf("%v/favicon.ico", url)
	content, err := FromURLGetContent(faviconUrl)
	if err != nil {
		return ""
	}
	return Mmh3Hash32(StandBase64(content))

}

func Mmh3Hash32(raw []byte) string {
	var h32 hash.Hash32 = murmur3.New32()
	h32.Write(raw)
	if IsUint32 {
		return fmt.Sprintf("%d", h32.Sum32())
	}
	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

// StandBase64 计算 base64 的值
func StandBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}

type MyError struct {
	message string
}

func (e MyError) Error() string {
	return e.message
}
