package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"
const SLEEP = 2500

const nicUrl string = "http://www.nic.io/cgi-bin/whois"
const failureMsg string = "DomainNotFound"
const reservedMsg string = "Reserved Auction"
const successMsg string = "Whois Search Successful"

// Return == 0: available
// Return == 1: reserved
// Return == 2: occupied
// Return == 3: error occurred
func lookup(domain string) int {
	formData := url.Values{}
	formData.Set("query", domain)

	resp, err := http.PostForm(nicUrl, formData)
	if err != nil {
		fmt.Printf("ERROR: could not query NIC for domain %s\n", domain)
		return 3
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: could not read response body\n")
		return 3
	}

	resp.Body.Close()

	if strings.Contains(string(body), failureMsg) {
		return 0
	} else if strings.Contains(string(body), reservedMsg) {
		return 1
	} else if strings.Contains(string(body), successMsg) {
		return 2
	} else {
		return 3
	}
}

func respond(result int) {
	switch result {
	case 0:
		fmt.Println("AVAILABLE")
	case 1:
		fmt.Println("RESERVED")
	case 2:
		fmt.Println("OCCUPIED")
	default:
		fmt.Println("ERROR OCCURRED")
	}
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		result := lookup(arg)
		respond(result)
		os.Exit(0)
	}

	for a := range ALPHABET {
		for b := range ALPHABET {
			for c := range ALPHABET {
				domain := fmt.Sprintf("%c%c%c.io", ALPHABET[a], ALPHABET[b], ALPHABET[c])
				fmt.Printf("%s: ", domain)
				respond(lookup(domain))
				time.Sleep(SLEEP * time.Millisecond)
			}
		}
	}
}
