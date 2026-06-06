package files

import (
	"amper/common/util"
	"fmt"
	"os"
)

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func RemoveAll(filePath string) bool {
	if err := os.RemoveAll(filePath); err != nil {
		util.Loggify(err)
		return false
	}
	return true
}

func RecreateIfExist(dir string, perm os.FileMode) error {
	if Exists(dir) && !RemoveAll(dir) {
		return fmt.Errorf("not able to remove the directory and it's content for %s", dir)
	}
	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func Name(value *string) *string {
	if value == nil || len(*value) == 0 {
		return util.UUID()
	} else if len(*value) > 255 {
		return util.PointerString((*value)[0:255])
	} else {
		return value
	}
}
