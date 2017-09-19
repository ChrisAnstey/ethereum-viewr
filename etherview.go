package main

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

func callApiWithParams(method string, params interface{}) []byte {

    url := "http://192.168.1.145:8545"
    url = "http://127.0.0.1:8545"

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

    req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(queryJson))
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

    return body

}

func callApi(method string) []byte {

    return callApiWithParams(method, [0]string{})

}

func blockNumber() string {

    body := callApi("eth_blockNumber")

    output := fmt.Sprintf("Res:  %s!", body)

    var dat map[string]interface{}

    if err := json.Unmarshal(body, &dat); err != nil {
        panic(err)
    }
    output += fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("Res2:  %s!", dat)
    output += fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("Res3:  %s!", dat["result"])
    output += fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("<br /><br />")

    return output
}

func syncing() string {

    body := callApi("eth_syncing")

    output := fmt.Sprintf("Res:  %s!", body)

    var dat map[string]interface{}

    if err := json.Unmarshal(body, &dat); err != nil {
        panic(err)
    }
    output += fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("Res2:  %s!", dat)
    output += fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("Res3:  %s!", dat["result"])


    output += fmt.Sprintf("<br /><br />")
    switch vv := dat["result"].(type) {
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

    // for k, v := range dat {
    //     switch vv := v.(type) {
    //     case string:
    //         fmt.Println(k, "is string", vv)
    //     case int:
    //         fmt.Println(k, "is int", vv)
    //     case float64:
    //         fmt.Println(k, "is float64", vv)
    //     case map[string]interface {}:
		  //   output += fmt.Sprintf("<table>")
    //         for i, u := range vv {
			 //    output += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", i, u)
    //         }
		  //   output += fmt.Sprintf("<table>")
    //     default:
    //         fmt.Println(k, "is of a type I don't know how to handle", reflect.TypeOf(vv))
    //     }
    // }

    return output

}


func getBlock() string {

    body := callApiWithParams("eth_getBlockByNumber", []interface{}{"latest", true})

    output := fmt.Sprintf("Res:  %s!", body)

    var dat map[string]interface{}

    if err := json.Unmarshal(body, &dat); err != nil {
        panic(err)
    }
    output += fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("Res2:  %s!", dat)
    output += fmt.Sprintf("<br /><br />")
    output += fmt.Sprintf("Res3:  %s!", dat["result"])


    output += fmt.Sprintf("<br /><br />")
    switch vv := dat["result"].(type) {
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

    // for k, v := range dat {
    //     switch vv := v.(type) {
    //     case string:
    //         fmt.Println(k, "is string", vv)
    //     case int:
    //         fmt.Println(k, "is int", vv)
    //     case float64:
    //         fmt.Println(k, "is float64", vv)
    //     case map[string]interface {}:
		  //   output += fmt.Sprintf("<table>")
    //         for i, u := range vv {
			 //    output += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", i, u)
    //         }
		  //   output += fmt.Sprintf("<table>")
    //     default:
    //         fmt.Println(k, "is of a type I don't know how to handle", reflect.TypeOf(vv))
    //     }
    // }

    return output

}




func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!doctype html>
		<html lang="en">
		<head>
		  <meta charset="utf-8">
		  <title>Etherview</title>
		</head>
		<body>
	`)

    body := syncing()
    body += blockNumber()
    body += getBlock()

    fmt.Fprintf(w, "%s", body)

	fmt.Fprintf(w, `</body></html>`)

}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8088", nil)
}