package lib

import (
    "fmt"
    "net/http"
     // "net/url"
    "encoding/json"
    "time"
    "log"
    "io/ioutil"
    "bytes"
    "reflect"
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

    dat := c.callApi("eth_blockNumber")

    return dat
}

func (c *Client) Syncing() string {

    result := c.callApi("eth_syncing")

    output := fmt.Sprintf("<br />")

    switch vv := result.(type) {
	    case bool:
		    output += fmt.Sprintf("Not syncing")
	    case map[string]interface {}:
	        output += fmt.Sprintf("Syncing")
		    output += fmt.Sprintf("<table>")
            for i, u := range vv {
			    output += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", i, u)
            }
		    output += fmt.Sprintf("<table>")
	    default:
	        output += fmt.Sprintf("is of a type I don't know how to handle", reflect.TypeOf(vv))
    }
    output += fmt.Sprintf("<br /><br />")

    return output
}

func (c *Client) IsSyncing() bool {

    result := c.callApi("eth_syncing")

    syncing := false

    switch result.(type) {
        case map[string]interface {}:
            syncing = true
    }

    return syncing

}


func (c *Client) GetBlock(blockNum string) string {

    result := c.callApiWithParams("eth_getBlockByNumber", []interface{}{blockNum, true})

    output := fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("Res3:  %s!", result)


    output += fmt.Sprintf("<br /><br />")
    switch vv := result.(type) {
	    case map[string]interface {}:
		    output += fmt.Sprintf("<table>")
            for i, u := range vv {
			    output += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", i, u)
            }
		    output += fmt.Sprintf("<table>")
	    default:
	        output += fmt.Sprintf("is of a type I don't know how to handle", reflect.TypeOf(vv))
    }
    output += fmt.Sprintf("<br /><br />")

    return output
}


func (c *Client) GetTxn(txNum string) interface{} {
    result := c.callApiWithParams("eth_getTransactionByHash", []interface{}{txNum})

    return result
}

