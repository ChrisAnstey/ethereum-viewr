package lib

import (
    "net/http"
     // "net/url"
    "encoding/json"
    "time"
    "log"
    "io/ioutil"
    "bytes"
    "fmt"
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

type Transaction struct {
    Hash string
    Data map[string]string
}

type Block struct {
    Hash string
    Data map[string]string
    Transactions map[string]Transaction
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

func (c *Client) GetBlockDataByNumber(blockNum string)  Block {
    result := c.callApiWithParams("eth_getBlockByNumber", []interface{}{blockNum, true})

    return extractBlockData(result)
}

func (c *Client) GetBlockDataByHash(blockHash string) Block {
    result := c.callApiWithParams("eth_getBlockByHash", []interface{}{blockHash, true})

    return extractBlockData(result)
}

func extractBlockData(input interface{}) Block {

    var response Block
    data := make(map[string]string)

    for i, u := range input.(map[string]interface {}) {
        switch v := u.(type) {
            case string:
                data [i] = v
            case []interface {}:
                if i == "transactions" {
                    transactions := make(map[string]Transaction)
                    for _, tu := range u.([]interface {}) {
                        tdata := make(map[string]string)
                        for tti, ttu := range tu.(map[string]interface {}) {
                            if  ttus, ok := ttu.(string); ok {
                                tdata[tti] = ttus
                            }
                        }

                        transactions[tdata["hash"]] = Transaction{tdata["hash"], tdata}
                    }
                    response.Transactions = transactions
                } else {

                    fmt.Printf(i, "unexpected type %T", v)
                }
            default:
                fmt.Printf(i, "unexpected type %T", v)
        }
    }
    response.Data = data
    response.Hash = data["hash"]

    return response
}




func (c *Client) GetTxn(txNum string) interface{} {
    result := c.callApiWithParams("eth_getTransactionByHash", []interface{}{txNum})

    return result
}

