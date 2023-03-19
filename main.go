package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
)

var (
	claimStatus = false
	targets     []string
)

func main() {
	fmt.Println("Checking for a target list...")

	if _, err := os.Stat("targets.txt"); err == nil {
		fmt.Println("Target list has been found!")
		file, _ := os.Open("targets.txt")
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			targets = append(targets, strings.TrimSpace(scanner.Text()))
		}
		file.Close()
		claimStatus = true
	} else {
		fmt.Println("No target list detected:")
	}

	var target string
	if !claimStatus {
		fmt.Print("Target: ")
		fmt.Scanln(&target)
	}

	var email, sessionID, csrfToken string
	var threads int

	fmt.Print("Email: ")
	fmt.Scanln(&email)
	fmt.Print("Session ID: ")
	fmt.Scanln(&sessionID)
	fmt.Print("CSRF Token: ")
	fmt.Scanln(&csrfToken)
	fmt.Print("Threads: ")
	fmt.Scanln(&threads)

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			claim(csrfToken, email, sessionID, target)
		}()
	}
	wg.Wait()
}

func claim(csrfToken, email, sessionID, target string) {
	headers := []string{
		"host: solo.to",
		"connection: keep-alive",
		`sec-ch-ua: "Not A;Brand";v="99", "Chromium";v="102", "Google Chrome";v="102"`,
		"x-csrf-token: " + csrfToken,
		"sec-ch-ua-mobile: ?0",
		"user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36",
		"content-type: application/x-www-form-urlencoded; charset=UTF-8",
		"accept: */*",
		"x-requested-with: XMLHttpRequest",
		'sec-ch-ua-platform: "macOS"',
		"origin: https://solo.to",
		"sec-fetch-site: same-origin",
		"sec-fetch-mode: cors",
		"sec-fetch-dest: empty",
		"referer: https://solo.to/account",
		"accept-encoding: utf-8",
		"accept-language: en-GB,en-US;q=0.9,en;q=0.8",
		"cookie: soloto_session=" + sessionID,
	}

	var data string
	client := &fasthttp.Client{}

	if claimStatus {
		for _, target := range targets {
			data = fmt.Sprintf("_token=%s&email=%s&username=%s&domain=", csrfToken, email, target)
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()

			defer fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(resp)

			req.SetRequestURI("https://solo.to/account/update-info/1")
			req.Header.SetMethod("POST")
			req.Header.AddBytesV("Content-Type", []byte("application/x-www-form-urlencoded"))

			for _, h := range headers {
				header := strings.SplitN(h, ": ", 2)
				req.Header.Add(header[0], header[1])
			}

			req.SetBodyString(data)

			if err := client.Do(req, resp); err != nil {
				fmt.Println("error")
			}
		}
	} else {
		data = fmt.Sprintf("_token=%s&email=%s&username=%s&domain=", csrfToken, email, target)
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()

		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.SetRequestURI("https://solo.to/account/update-info/1")
		req.Header.SetMethod("POST")
		req.Header.AddBytesV("Content-Type", []byte("application/x-www-form-urlencoded"))

		for _, h := range headers {
			header := strings.SplitN(h, ": ", 2)
			req.Header.Add(header[0], header[1])
		}

		req.SetBodyString(data)

		if err := client.Do(req, resp); err != nil {
			fmt.Println("poor")
		}
	}
}

