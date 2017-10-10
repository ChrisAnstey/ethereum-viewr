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
    syncData, _ := gethClient.IsSyncing()
    blockNumber, _ := gethClient.BlockNumber()

    var PageVars = struct{PageTitle string; LatestBlock interface{}; SyncData lib.EthSyncingResponse}{"Status", blockNumber, syncData}

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
      blockData, _ = gethClient.GetBlockDataByNumber(block)
    } else {
	    // otherwise, try hash
      blockData, _ = gethClient.GetBlockDataByHash(r.Form.Get("blockHash"))
    }

    var PageVars = struct{PageTitle string; BlockData lib.Block}{"View Block", blockData}

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
    txData, err := gethClient.GetTxn(tx)
    if err != nil {
      log.Print("API error: ", err)
      http.Error(w, "Error", 500)
      return
    }

    txReceipt, err := gethClient.GetTxnReceipt(tx)
    if err != nil {
      log.Print("API error: ", err)
      http.Error(w, "Error", 500)
      return
    }

    var PageVars = struct{PageTitle string; Txn lib.Transaction; TxReceipt lib.TransactionReceipt}{"View Transaction", txData, txReceipt}

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

func viewAccount(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    acc := r.Form.Get("acc")
    accData, _ := gethClient.GetAccountBalance(acc)

    var PageVars = struct{PageTitle string; Acc lib.Account}{"View Account", accData}

    t, err := template.ParseFiles("html/page/account.html", "html/layout/template.html") //parse the html files
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
  log.SetFlags(log.LstdFlags | log.Lshortfile)
    http.HandleFunc("/", status)
    http.HandleFunc("/status", status)
    http.HandleFunc("/block", viewBlock)
    http.HandleFunc("/tx", viewTransaction)
    http.HandleFunc("/account", viewAccount)
    log.Fatal(http.ListenAndServe(":8088", nil))
}