package main

import (
	"github.com/ChrisAnstey/etherview/lib"
	"html/template"
	"log"
	"net/http"
)

type PageVariables struct {
	Body template.HTML
}

var gethClient = lib.Client{
	Url: "http://192.168.1.145:8545",
}

func status(w http.ResponseWriter, r *http.Request) {
	syncData, err := gethClient.IsSyncing()
	if err != nil {
		log.Print("API error: ", err)
		http.Error(w, "Error", 500)
		return
	}

	blockNumber, err := gethClient.BlockNumber()
	if err != nil {
		log.Print("API error: ", err)
		http.Error(w, "Error", 500)
		return
	}

	var PageVars = struct {
		PageTitle   string
		LatestBlock interface{}
		SyncData    lib.EthSyncingResponse
	}{"Status", blockNumber, syncData}

	outputPage(w, "html/page/status.html", PageVars)
}

func viewBlock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get by number
	var block string
	var blockData lib.Block
	var err error

	// check if we got a number
	if block = r.Form.Get("block"); block != "" {
		blockData, err = gethClient.GetBlockDataByNumber(block)
	} else {
		// otherwise, try hash
		blockData, err = gethClient.GetBlockDataByHash(r.Form.Get("blockHash"))
	}
	if err != nil {
		log.Print("API error: ", err)
		http.Error(w, "Error finding block", 500)
		return
	}

	var PageVars = struct {
		PageTitle string
		BlockData lib.Block
	}{"View Block", blockData}

	outputPage(w, "html/page/block.html", PageVars)
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
		http.Error(w, "Error finding Transaction", 500)
		return
	}

	var PageVars = struct {
		PageTitle string
		Txn       lib.Transaction
		TxReceipt lib.TransactionReceipt
	}{"View Transaction", txData, txReceipt}

	outputPage(w, "html/page/transaction.html", PageVars)
}

func viewAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	acc := r.Form.Get("acc")
	accData, err := gethClient.GetAccountBalance(acc)
	if err != nil {
		log.Print("API error: ", err)
		http.Error(w, "Error", 500)
		return
	}

	var PageVars = struct {
		PageTitle string
		Acc       lib.Account
	}{"View Account", accData}

	outputPage(w, "html/page/account.html", PageVars)
}

func outputPage(w http.ResponseWriter, pageTemplate string, pageVars interface{}) {
	t, err := template.ParseFiles(pageTemplate, "html/layout/template.html") //parse the html files
	if err != nil {
		log.Print("template parsing error: ", err)
		http.Error(w, "Error", 500)
		return
	}

	//execute the template, pass it the PageVars struct to fill in the gaps, and the ResponseWriter to output the result
	err = t.ExecuteTemplate(w, "layout", pageVars)
	if err != nil {
		log.Print("template executing error: ", err)
		http.Error(w, "Error", 500)
		return
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
