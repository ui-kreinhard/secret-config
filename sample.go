package main

import (
	"fmt"
	"os"

	"github.com/ui-kreinhard/secret-config/crypt"
	urltag "github.com/ui-kreinhard/secret-config/url-tag"
)

type Config struct {
	Secret1       string `secret_url:"https://bitwarden.synyx.coffee/#/vault?itemId=184e5343-442d-4aa9-ba1d-b13c007fe2b8"`
	Secret2       string `secret_url:"https://bitwarden.synyx.coffee/#/vault?itemId=184e5343-442d-4aa9-ba1d-b13c007fe2b8"`
	AnotherSecret string `secret_url:"https://bitwarden.synyx.coffee/#/vault?itemId=184e5343-442d-4aa9-ba1d-b13c007fe2b8"`
	NonSecret     string
}

func encryptConfig() {
	config := Config{
		Secret1:       "This",
		Secret2:       "is",
		AnotherSecret: "very",
		NonSecret:     "secure",
	}
	key, err := crypt.GenKey()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(key)
	config = urltag.ScanForUrlAndEncrypt(config, key)
	fmt.Println(config)
	fmt.Printf("%#v", config)
}

func decryptConfig() {
	c := Config{
		Secret1:       "3d6d08ad1dbbc9d5bfa46d7966859883105b7f1fd3f85e5d616e84535c7509ed",
		Secret2:       "ae3d691abe01002544bb42da319eb909430557b91f5076bb9a76f8a9cd6a",
		AnotherSecret: "768b6f5abdc0332b7f81b2b25d0ad6e267c8d833198f4dd5668933ffadb5efaf",
		NonSecret:     "secure",
	}
	c = urltag.ScanForUrlAndOpen(c, "DEV_MODE")
	fmt.Printf("%#v\n", c)
}

func main() {
	decryptConfig()
}
