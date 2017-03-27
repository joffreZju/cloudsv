package main

import (
	"flag"
	"fmt"
	"s4s/common/lib/keycrypt"
)

var (
	key        = flag.String("k", "", "crypt key")
	plaintext  = flag.String("e", "", "plaintext")
	ciptertext = flag.String("d", "", "ciptertext")

	helper = `
s4skeycipter help

encode text
	s4skeycipter -k "this is key" -e "this is plaintext"

encode text
	s4skeycipter -k "this is key" -d "this is ciptertext"
	`
)

func main() {
	flag.Parse()
	if *key == "" {
		fmt.Println(helper)
		return
	}
	if len(*plaintext) > 0 {
		fmt.Println(keycrypt.Encode(*key, *plaintext))
		return
	}
	if len(*ciptertext) > 0 {
		fmt.Println(keycrypt.Decode(*key, *ciptertext))
		return
	}
	fmt.Println(helper)
}
