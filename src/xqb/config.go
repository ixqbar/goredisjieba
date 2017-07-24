package xqb

import (
	"errors"
	"os"
	"encoding/xml"
	"fmt"
)

type JiebaXmlConfig struct {
	Address string `xml:"address"`
	DB int `xml:"db"`
	DictPath string `xml:"dict"`
}

var jiebaXmlConfig *JiebaXmlConfig

func ParseXmlConfig(path string) (*JiebaXmlConfig, error) {
	if len(path) == 0 {
		return nil, errors.New("not found configure xml file")
	}

	r,e := os.Stat(path);
	if  e !=nil || r.Size() == 0 {
		return nil, errors.New("not found configure xml file")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	jiebaXmlConfig = &JiebaXmlConfig{
		Address : ":6379",
		DB : 0,
		DictPath : "",
	}

	data := make([]byte, r.Size())

	n, err := f.Read(data)
	if err != nil {
		return nil, err
	}

	if int64(n) != r.Size() {
		return nil, errors.New(fmt.Sprintf("expect read configure xml file size %d but result is %d", r.Size(), n))
	}

	err = xml.Unmarshal(data, &jiebaXmlConfig)
	if err != nil {
		return nil, err
	}

	return jiebaXmlConfig, nil
}