package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	host := os.Getenv("HOST")

	p, _ := os.Getwd()
	path := fmt.Sprintf("%s/swagger/ui", p)
	if host != "" {
		err := addHost(fmt.Sprintf("%v/swagger.json", path), host)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fs := http.FileServer(http.Dir(path))
	http.Handle("/ui/", http.StripPrefix("/ui/", fs))

	fmt.Printf("serving swagger at :%v", 8080)

	err := http.ListenAndServe(":8080", fs)
	if err != nil {
		fmt.Println(err)
	}
}

func addHost(path, host string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	swagger := make(map[string]json.RawMessage)
	_ = json.Unmarshal(bytes, &swagger)
	b, _ := json.Marshal(host)

	swagger["host"] = b
	swaggerJson, _ := json.Marshal(swagger)

	_ = os.Remove(path)
	_ = os.WriteFile(path, swaggerJson, 0644)

	return nil
}
