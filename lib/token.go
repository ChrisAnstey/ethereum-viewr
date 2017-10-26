package lib

import (
	"encoding/hex"
	"strconv"
	"strings"
	"errors"
)

type Token struct {
	Address     string
	Client      Client
	Name        string
	TotalSupply int64
}

type Call struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

func (t *Token) GetName() (string, error) {

	result, err := t.Client.callApiWithParams("eth_call", []interface{}{&Call{To: t.Address, Data: "0x06fdde03"}, "latest"})

	if err != nil || len(result.(string)) == 2 {
		return "", errors.New("Not found")
	}

	//ignore the first 66 bytes ("0x", plus a 64 byte value for the length of the rest)
	decoded, err := hex.DecodeString(result.(string)[66:])
	if err != nil {
		return "", err
	}

	// trim off the null chars
	t.Name = strings.Trim(string(decoded), "\x00")
	return t.Name, nil

}

func (t *Token) GetTotalSupply() (int64, error) {

	result, err := t.Client.callApiWithParams("eth_call", []interface{}{&Call{To: t.Address, Data: "0x18160ddd"}, "latest"})
	if err != nil {
		return 0, err

	}
	supply, err := strconv.ParseInt(result.(string), 0, 64)
	if err != nil {
		return 0, err
	}
	t.TotalSupply = supply
	return supply, err

}
