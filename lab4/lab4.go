package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// TODO: Create a struct to hold the data sent to the template
type Data struct {
	Result  string
	Expression string
}

func CalculateGCD(a, b int) int {
	if b == 0 {
		return a
	}
	return CalculateGCD(b, a%b)
}

func CalculateLCM(a, b int) int {
	return (a * b) / CalculateGCD(a, b)
}

func Calculator(w http.ResponseWriter, r *http.Request) {
	// TODO: Finish this function
	op := r.URL.Query().Get("op")
	num1Str := r.URL.Query().Get("num1")
	num2Str := r.URL.Query().Get("num2")

	num1, err1 := strconv.Atoi(num1Str)
	num2, err2 := strconv.Atoi(num2Str)

	var expression, result string

	if err1 != nil || err2 != nil {
		tmpl := template.Must(template.ParseFiles("error.html"))
		tmpl.Execute(w, nil)
		return
	} else {
		switch op {
		case "add":
			result = strconv.Itoa(num1 + num2)
			expression = fmt.Sprintf("%d + %d", num1, num2)
		case "sub":
			result = strconv.Itoa(num1 - num2)
			expression = fmt.Sprintf("%d - %d", num1, num2)
		case "mul":
			result = strconv.Itoa(num1 * num2)
			expression = fmt.Sprintf("%d * %d", num1, num2)
		case "div":
			if num2 == 0 {
				tmpl := template.Must(template.ParseFiles("error.html"))
				tmpl.Execute(w, nil)
				return
			} else {
				result = fmt.Sprintf("%d", num1 / num2)
			}
			expression = fmt.Sprintf("%d / %d", num1, num2)
		case "gcd":
			result = strconv.Itoa(CalculateGCD(num1, num2))
			expression = fmt.Sprintf("GCD(%d, %d)", num1, num2)
		case "lcm":
			result = strconv.Itoa(CalculateLCM(num1, num2))
			expression = fmt.Sprintf("LCM(%d, %d)", num1, num2)
		default:
			tmpl := template.Must(template.ParseFiles("error.html"))
			tmpl.Execute(w, nil)
			return
		}
	}
	
	data := Data{
		Result:     result,
		Expression: expression,
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", Calculator)
	log.Fatal(http.ListenAndServe(":8084", nil))
}