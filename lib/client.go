package lib

import (
    "net/http"
     // "net/url"
    "encoding/json"
    "time"
    "log"
    "io/ioutil"
    "bytes"
)

type Request1 struct {
    Jsonrpc   string      `json:"jsonrpc"`
    Method   string `json:"method"`
    Id   int `json:"id"`
    Params interface{} `json:"params"`
}

type Client struct{
    Url string
}

func (c *Client) callApiWithParams(method string, params interface{}) interface{} {

    queryData := &Request1{
        Jsonrpc:   "2.0",
        Method: method,
        Id: 1,
        Params: params,
    }
    queryJson, _ := json.Marshal(queryData)
    // fmt.Println(string(queryJson))

    ethClient := http.Client{
        Timeout: time.Second * 2, // Maximum of 2 secs
    }

    req, err := http.NewRequest(http.MethodGet, c.Url, bytes.NewBuffer(queryJson))
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("User-Agent", "Etherview")

    // make the request
    res, getErr := ethClient.Do(req)
    if getErr != nil {
        log.Fatal(getErr)
    }

    // read in the result
    body, readErr := ioutil.ReadAll(res.Body)
    if readErr != nil {
        log.Fatal(readErr)
    }

    var response map[string]interface{}

    if err := json.Unmarshal(body, &response); err != nil {
        panic(err)
    }

    return response["result"]

}

func (c *Client) callApi(method string) interface{} {

    return c.callApiWithParams(method, [0]string{})

}

func (c *Client) BlockNumber() interface{} {

    result := c.callApi("eth_blockNumber")

    return result
}


func (c *Client) IsSyncing() (bool, map[string]interface {}) {

    result := c.callApi("eth_syncing")

    var data map[string]interface {};
    syncing := false

    switch  resultValue := result.(type) {
        case map[string]interface {}:
            syncing = true
            data = resultValue
    }

    return syncing, data

}

func (c *Client) GetBlockData(blockNum string)  interface {} {
    result := c.callApiWithParams("eth_getBlockByNumber", []interface{}{blockNum, true})

    return result
}

func (c *Client) GetBlockDataByHash(blockHash string)  interface {} {
    result := c.callApiWithParams("eth_getBlockByHash", []interface{}{blockHash, true})

    return result
}


func (c *Client) GetTxn(txNum string) interface{} {
    result := c.callApiWithParams("eth_getTransactionByHash", []interface{}{txNum})

    return result
}

