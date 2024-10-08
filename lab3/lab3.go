package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Calculator(w http.ResponseWriter, r *http.Request) {
	// TODO: implement a calculator
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 3 {
			fmt.Fprint(w, "Error!")
			return
	}

	operation := parts[0]
	num1, err1 := strconv.Atoi(parts[1])
	num2, err2 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil {
			fmt.Fprint(w, "Error!")
			return
	}


	switch operation {
	case "add":
			fmt.Fprintf(w, "%d + %d = %d", num1, num2, num1+num2)
	case "sub":
			fmt.Fprintf(w, "%d - %d = %d", num1, num2, num1-num2)
	case "mul":
			fmt.Fprintf(w, "%d * %d = %d", num1, num2, num1*num2)
	case "div":
			if num2 == 0 {
					fmt.Fprint(w, "Error!")
			} else {
					fmt.Fprintf(w, "%d / %d = %d, remainder %d", num1, num2, num1/num2, num1%num2)
			}
	default:
			fmt.Fprint(w, "Error!")
	}

}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8083", nil))
}