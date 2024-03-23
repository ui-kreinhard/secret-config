package urltag

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"syscall"

	"github.com/ui-kreinhard/secret-config/crypt"
	"golang.org/x/term"
)

func ScanForUrlAndOpen[T any](object T, envName string) T {
	urlSecretMap := make(map[string]string)
	if _, isSet := os.LookupEnv(envName); !isSet {
		log.Println("env not set")
		return object
	}
	rt := reflect.TypeOf(object)
	rf := reflect.Indirect(reflect.ValueOf(&object))
	n := rt.NumField()
	for i := 0; i < n; i++ {
		urlOfTag := rt.Field(i).Tag.Get("secret_url")
		cipher := rf.Field(i).String()
		if urlOfTag != "" {
			pwAsString := urlSecretMap[urlOfTag]
			if urlSecretMap[urlOfTag] == "" {
				fmt.Println("Opening url", urlOfTag)
				fmt.Print("Enter secret:")
				openbrowser(urlOfTag)
				bytepw, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					fmt.Println("err reading pw", err)
					os.Exit(1)
				}
				fmt.Println()
				pwAsString = string(bytepw)
				urlSecretMap[urlOfTag] = pwAsString
			}
			decrypted, err := crypt.Decrypt(cipher, pwAsString)
			if err != nil {
				fmt.Println("err during decryption", err)
				os.Exit(1)
			}
			rf.Field(i).SetString(decrypted)
		}
	}
	return object
}

func ScanForUrlAndEncrypt[T any](object T, key string) T {
	rt := reflect.TypeOf(object)
	rf := reflect.Indirect(reflect.ValueOf(&object))
	n := rt.NumField()
	for i := 0; i < n; i++ {
		urlOfTag := rt.Field(i).Tag.Get("secret_url")
		cipher := rf.Field(i).String()
		if urlOfTag != "" {
			encrypted, err := crypt.Encrypt(cipher, key)
			if err != nil {
				fmt.Println("err during encryption", err)
				os.Exit(1)
			}
			rf.Field(i).SetString(encrypted)
		}
	}
	return object
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
