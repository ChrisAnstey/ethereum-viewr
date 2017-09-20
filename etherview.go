package main

import (
    "fmt"
    "net/http"
    "html/template"
     // "net/url"
    "encoding/json"
    "time"
    "log"
    "io/ioutil"
    "bytes"
    "reflect"
)

type PageVariables struct {
	Body         template.HTML
}

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
    body := syncing()
    body += blockNumber()
    body += getBlock()

    PageVars := PageVariables{ //store the data in a struct
      Body: template.HTML(body),
    }

    t, err := template.ParseFiles("html/layout/template.html") //parse the html file
    if err != nil {
  	  log.Print("template parsing error: ", err)
  	}

  	//execute the template, pass it the PageVars struct to fill in the gaps, and the ResponseWriter to output the result
    err = t.Execute(w, PageVars)
    if err != nil {
  	  log.Print("template executing error: ", err)
	}
}


func status(w http.ResponseWriter, r *http.Request) {
    body := syncing()
    body += blockNumber()

    PageVars := PageVariables{ //store the data in a struct
      Body: template.HTML(body),
    }

    t, err := template.ParseFiles("html/layout/template.html") //parse the html file
    if err != nil {
  	  log.Print("template parsing error: ", err)
  	}

  	//execute the template, pass it the PageVars struct to fill in the gaps, and the ResponseWriter to output the result
    err = t.Execute(w, PageVars)
    if err != nil {
  	  log.Print("template executing error: ", err)
	}
}


func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/status", status)
    http.ListenAndServe(":8088", nil)
}