package urltag

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"

	"github.com/ui-kreinhard/secret-config/crypt"
	"golang.org/x/term"
)

func isDevModeSet(envName string) bool {
	_, isSet := os.LookupEnv(envName)
	return isSet
}

func getCacheFile() string {
	name, _ := os.Executable()
	name = filepath.Base(name)
	return filepath.Join("/tmp", name)
}

func checkForCacheFile() bool {
	tmpFile := getCacheFile()
	if _, err := os.Stat(tmpFile); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func writeCacheFile(urlSecretMap map[string]string) error {
	data, err := json.Marshal(urlSecretMap)
	if err != nil {
		return err
	}
	cacheFile := getCacheFile()
	return os.WriteFile(cacheFile, data, 0644)
}

func readCacheFile() (map[string]string, error) {
	ret := map[string]string{}
	rawData, err := os.ReadFile(getCacheFile())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawData, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ScanForUrlAndOpen[T any](object T, envName string) T {
	urlSecretMap := make(map[string]string)
	if !isDevModeSet(envName) {
		log.Println("env not set")
		return object
	}
	if checkForCacheFile() {
		cachedSecretMap, err := readCacheFile()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Found cached file")
			urlSecretMap = cachedSecretMap
		}
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
	err := writeCacheFile(urlSecretMap)
	if err != nil {
		log.Println("Cannot write cache file", err)
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
