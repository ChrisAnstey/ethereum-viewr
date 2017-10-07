package main

import (
    "net/http"
    "html/template"
    "log"
    "github.com/ChrisAnstey/etherview/lib"
)

type PageVariables struct {
	Body         template.HTML
}

var gethClient = lib.Client{
      Url: "http://192.168.1.145:8545",
    }

func status(w http.ResponseWriter, r *http.Request) {
    status, syncData := gethClient.IsSyncing()

    var PageVars = struct{Syncing, LatestBlock interface{}; SyncData map[string]interface{}}{status, gethClient.BlockNumber(), syncData}

    t, err := template.ParseFiles("html/page/status.html", "html/layout/template.html") //parse the html file
    if err != nil {
  	  log.Print("template parsing error: ", err)
  	}

  	//execute the template, pass it the PageVars struct to fill in the gaps, and the ResponseWriter to output the result
    err = t.ExecuteTemplate(w, "layout", PageVars)
    if err != nil {
  	  log.Print("template executing error: ", err)
	}
}

func viewBlock(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    // get by number
    var block string
    var blockData lib.Block


    // check if we got a number
    if block = r.Form.Get("block"); block != "" {
      blockData = gethClient.GetBlockDataByNumber(block)
    } else {
	    // otherwise, try hash
      blockData = gethClient.GetBlockDataByHash(r.Form.Get("blockHash"))
    }

    var PageVars = struct{BlockData lib.Block}{blockData}

    t, err := template.ParseFiles("html/page/block.html", "html/layout/template.html") //parse the html file
    if err != nil {
  	  log.Print("template parsing error: ", err)
  	}

  	//execute the template, pass it the PageVars struct to fill in the gaps, and the ResponseWriter to output the result
    err = t.ExecuteTemplate(w, "layout", PageVars)
    if err != nil {
    	  log.Print("template executing error: ", err)
  	}
}

func viewTransaction(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    tx := r.Form.Get("tx")
    txData := gethClient.GetTxn(tx)

    var PageVars = struct{Txn lib.Transaction}{txData}

    t, err := template.ParseFiles("html/page/transaction.html", "html/layout/template.html") //parse the html files
    if err != nil {
      log.Print("template parsing error: ", err)
    }

    //execute the template, pass it the PageVars struct to fill in the gaps, and the ResponseWriter to output the result
    err = t.ExecuteTemplate(w, "layout", PageVars)
    if err != nil {
      log.Print("template executing error: ", err)
    }
}


func main() {
    http.HandleFunc("/", status)
    http.HandleFunc("/status", status)
    http.HandleFunc("/block", viewBlock)
    http.HandleFunc("/tx", viewTransaction)
    log.Fatal(http.ListenAndServe(":8088", nil))
}