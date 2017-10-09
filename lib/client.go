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
    "strconv"
    "github.com/fatih/camelcase"
    "strings"
    "math"
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

type EthSyncingResponse struct{
    Status bool
    Data   map[string]string
}

type ApiResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
}

type Transaction struct {
    Hash string
    Value float64
    Data map[string]string
}

type Block struct {
    Hash string
    Number int64
    Timestamp time.Time
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
        Timeout: time.Second * 5, // Maximum of 5 secs
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

    var response ApiResponse

    if err := json.Unmarshal(body, &response); err != nil {
        panic(err)
    }

    return response.Result

}

func (c *Client) callApi(method string) interface{} {

    return c.callApiWithParams(method, [0]string{})

}

func (c *Client) BlockNumber() interface{} {

    result := c.callApi("eth_blockNumber")

    return result
}


func (c *Client) IsSyncing() EthSyncingResponse {

    result := c.callApi("eth_syncing")

    var response EthSyncingResponse
    response.Status = false

    switch v := result.(type) {
        case map[string]interface {}:
		    response.Status = true
		    data := make(map[string]string)
		    for ti, tu := range v {
		        if  tus, ok := tu.(string); ok {
		            data[humanise(ti)] = tus
		        }
		    }
		    response.Data = data
    }

    return response

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
                data [humanise(i)] = v
            case []interface {}:
                if i == "transactions" {
                    response.Transactions = extractTransactions(u)
                } else {
                    fmt.Printf(i, "unexpected type %T", v)
                }
            default:
                fmt.Printf(i, "unexpected type %T", v)
        }
    }
    response.Data = data
    response.Hash = data["Hash"]
    response.Number, _ = strconv.ParseInt(data["Number"], 0, 64)
    timestamp, _ := strconv.ParseInt(data["Timestamp"], 0, 64)
    response.Timestamp = time.Unix(timestamp, 0)

    return response
}

func extractTransactions(input interface{}) map[string]Transaction {
    transactions := make(map[string]Transaction)
    for _, tu := range input.([]interface {}) {
        transaction := extractTransactionData(tu)
        transactions[transaction.Data["hash"]] = transaction
    }
    return transactions
}

func extractTransactionData(input interface{}) Transaction {
    tdata := make(map[string]string)
    for ti, tu := range input.(map[string]interface {}) {
        if  tus, ok := tu.(string); ok {
            tdata[ti] = tus
        }
    }
    valueDec, _ := strconv.ParseInt(tdata["value"], 0, 64)
    return Transaction{tdata["hash"], float64(valueDec) / math.Pow10(18), tdata}
}

func humanise(input string) string {
	return strings.Title(strings.Join(camelcase.Split(input), " "))
}


func (c *Client) GetTxn(txNum string) Transaction {
    result := c.callApiWithParams("eth_getTransactionByHash", []interface{}{txNum})

    return extractTransactionData(result)
}

