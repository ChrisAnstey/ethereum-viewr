package main

import (
    "fmt"
    "net/http"
    "html/template"
    "log"
    "./lib"
)

type PageVariables struct {
	Body         template.HTML
}


func handler(w http.ResponseWriter, r *http.Request) {
    c := lib.Client{}
    body := c.Syncing()
    body += fmt.Sprintf("Res3:  %s!", c.BlockNumber())
    body += c.GetBlock("latest")

    PageVars := PageVariables{ //store the data in a struct
      Body: template.HTML(body),
    }

    t, err := template.ParseFiles("html/page/generic.html", "html/layout/template.html") //parse the html file
    if err != nil {
  	  log.Print("template parsing error: ", err)
  	}

  	//execute the template, pass it the PageVars struct to fill in the gaps, and the ResponseWriter to output the result
    err =t.ExecuteTemplate(w, "layout", PageVars)
    if err != nil {
  	  log.Print("template executing error: ", err)
	}
}


func status(w http.ResponseWriter, r *http.Request) {
    c := lib.Client{}
    body := c.Syncing()

    var PageVars = struct{Body, LatestBlock interface{}}{template.HTML(body), c.BlockNumber()}

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
    block := r.Form.Get("block")

    body := "Viewing Block: " + block + " "
    c := lib.Client{}
    body += c.GetBlock(block)

    PageVars := PageVariables{ //store the data in a struct
      Body: template.HTML(body),
    }

    t, err := template.ParseFiles("html/page/generic.html", "html/layout/template.html") //parse the html file
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
    http.HandleFunc("/", handler)
    http.HandleFunc("/status", status)
    http.HandleFunc("/block", viewBlock)
    log.Fatal(http.ListenAndServe(":8088", nil))
}