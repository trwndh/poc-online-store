package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func main() {

	client := &http.Client{}
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go DoRequest(&wg, client, i)
	}

	wg.Wait()
	fmt.Println("Test Done")
}

func DoRequest(wg *sync.WaitGroup, client *http.Client, sequence int) {
	defer wg.Done()
	url := "http://localhost:9999/api/v1/order/checkout?no_locking=true"
	method := "POST"
	str, _ := ioutil.ReadFile("request.json")
	payload := strings.NewReader(fmt.Sprintf(string(str), sequence, sequence))
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Sequence %d, res: %s", sequence, string(body))
}
