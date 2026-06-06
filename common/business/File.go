package business

import (
	"fmt"
	"io/ioutil"
	"log"
)

//Read is reading a file with the provided filepath
func Read(path *string) (string, error) {
	result, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Print(err.Error(), err)
		return "", fmt.Errorf("Unable to locate a file with the provided path %s", *path)
	}
	return string(result), nil
}
