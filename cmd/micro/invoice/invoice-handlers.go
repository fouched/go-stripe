package main

import (
	"fmt"
	"net/http"
)

func (app *application) CreateAndSendInvoice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello, world")
}
