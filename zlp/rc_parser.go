package zlp

import (
	"bufio"
	"os"
	"reflect"
	"strings"
)

type ZulipRC struct {
	Email  string `rc:"email"`
	APIKey string `rc:"key"`
	APIUrl string `rc:"site"`
}

const tagName = "rc"

// https://stackoverflow.com/a/55775573
func LoadRC(fileName string) (*ZulipRC, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(splitter)

	rc := ZulipRC{}
	rt := reflect.TypeOf(rc)

	fields := make(map[string]string)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		rcTag := field.Tag.Get(tagName)
		fields[rcTag] = field.Name
	}

	for scanner.Scan() {
		text := strings.Split(scanner.Text(), "=")
		if len(text) < 2 {
			continue
		}

		key, value := text[0], text[1]
		fieldName, ok := fields[key]
		if !ok {
			continue
		}

		reflect.ValueOf(&rc).Elem().FieldByName(fieldName).SetString(value)
	}

	return &rc, nil
}

func splitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if atEOF {
		return len(data), data, nil
	}

	if i := strings.Index(string(data), "\n"); i >= 0 {
		return i + 1, data[0:i], nil
	}
	return
}

func RCFromEnv() (*ZulipRC, error) {
	rc := ZulipRC{
		Email:  os.Getenv("ZULIP_EMAIL"),
		APIKey: os.Getenv("ZULIP_APIKEY"),
		APIUrl: os.Getenv("ZULIP_APIURL"),
	}

	return &rc, nil
}
