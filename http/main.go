package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type customWriter struct{
	body string
}


func main() {



	res, err := http.Get("http://google.com")
	
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	//method1
	// body:=make([]byte,99999)
	// item,_ := res.Body.Read(body)
	// fmt.Println(string(body))
	// fmt.Println(item)


	//methot2
	// io.Copy(os.Stdout,res.Body)

	

	//method3: custom writer struct
	cus := customWriter{}
	numberElement , _ := io.Copy(&cus,res.Body)
	fmt.Printf(strconv.Itoa(int(numberElement)))
	fmt.Println(cus.body)
	



	// byteBody := resp.Body

	//return value
	// testMake:= make([]byte, 99)

	//return &[]
	// newTest:= new([]byte)
	
	// fmt.Println(testMake)

	// lw := logWriter{}

	// io.Copy(lw, resp.Body)
}


//custom writer
func (cw *customWriter) Write(p []byte) (int,  error){
	// fmt.Printf(string(p))
	cw.body=string(p)
	return 1 , nil
}


































// func (logWriter) Write(bs []byte) (int, error) {
// 	fmt.Println(string(bs))
// 	fmt.Println("Just wrote this many bytes:", len(bs))
// 	return len(bs), nil
// }