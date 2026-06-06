package application

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/magiconair/properties"
)

var configuration *properties.Properties
var once sync.Once

// Get reads and loads the application configuration if not loaded yet
func Get() (result *properties.Properties, err error) {

	once.Do(func() {
		path, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Printf("unable to get root directory: %v", err)
			err = fmt.Errorf("unable to get the root directory")
		}
		var errP error
		configuration, errP = properties.LoadFile(path+"/application.properties", properties.UTF8)
		if errP != nil {
			log.Printf("unable to load application properties: %v", errP)
			err = fmt.Errorf("unable to load application properties")
		}
	})
	return configuration, err
}

func SetValue(values map[string]interface{}) (success bool, err error) {
	if values == nil || len(values) < 1 {
		return true, nil
	}
	appConfiduration, errAC := Get()
	if errAC != nil {
		log.Printf("nota able to load the app configuration with error: %v", err)
		return false, fmt.Errorf("nota able to load the app configuration")
	}
	for key, value := range values {
		appConfiduration.SetValue(key, value)
	}

	path, err := filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	if err != nil {
		log.Printf("unable to get root directory: %v", err)
		return false, fmt.Errorf("not anle to store the properties the root directory can't be found")
	}
	propertiesFile, errF := os.OpenFile(path+"/properties/application/application.properties", os.O_WRONLY, os.ModeAppend)
	if errF != nil {
		log.Printf(errF.Error(), errF)
		return false, fmt.Errorf("not anle to store data, the properties file can't be found")
	}
	n, errW := appConfiduration.Write(propertiesFile, properties.UTF8)
	if errW != nil || n < 1 {
		log.Printf(errW.Error(), errW)
		return false, fmt.Errorf("not anle to store the properties there was a data write issue")
	}
	return true, nil
}
