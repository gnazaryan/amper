package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
	"fmt"
	"strconv"
)

func UploadChunk(userId *int64, name *string, lastModified *int64, size *int64, Type *string, chunk *string, start *int64, directory *string) (success bool, startRes int64, fileMetadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "name": name, "size": size, "chunk": chunk})
	if err != nil {
		return false, -1, nil, err
	}
	return business.UploadChunk(userId, name, lastModified, size, Type, chunk, nil, start, directory)
}

func GetMetadataOld(userId *int64, directory *string, id *string, major *int64, minor *int64, patch *int64) (metadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "directory": directory, "id": id})
	if err != nil {
		return nil, err
	}
	var version *structs.Version
	if major != nil && minor != nil && patch != nil {
		version = &structs.Version{
			Major: *major,
			Minor: *minor,
			Patch: *patch,
		}
	}
	return business.GetMetadata(userId, directory, id, version)
}

func NewDirectoryOld(userId *int64, directory *string, name *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "name": name})
	if err != nil {
		return false, err
	}
	return business.NewDirectory(userId, directory, name)
}

func RemoveFileOld(userId *int64, directory *string, id *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "directory": directory, "id": id})
	if err != nil {
		return false, err
	}
	return business.RemoveFile(userId, directory, id)
}

func RemoveFilesOld(userId *int64, root *string, ids *[]string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "root": root, "ids": ids})
	if err != nil {
		return false, err
	}
	return business.RemoveFiles(userId, root, ids)
}

func MoveFilesOld(userId *int64, root *string, ids *[]string, directory *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "root": root, "ids": ids, "directory": directory})
	if err != nil {
		return false, err
	}
	return business.MoveFiles(userId, root, ids, directory)
}

func RemoveDirectory(userId *int64, directory *string, id *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "directory": directory, "id": id})
	if err != nil {
		return false, err
	}
	return business.RemoveFile(userId, directory, id)
}

func PasteFilesOld(userId *int64, root *string, copy *string, cut *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "root": root})
	if err != nil {
		return false, err
	}
	return business.PasteFiles(userId, root, copy, cut)
}

func DiscoverRootOld(userId *int64) (result *structs.Folder, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId})
	if err != nil {
		return nil, err
	}
	return business.DiscoverRoot(userId)
}
func UpversionFile(userId *int64, root *string, id *string, major *string, minor *string, patch *string, file *[]byte) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "root": root, "id": id, "major": major, "minor": minor, "patch": patch, "file": file})
	if err != nil {
		return false, err
	}
	majorV, errM := strconv.ParseInt(*major, 10, 64)
	minorV, errMi := strconv.ParseInt(*minor, 10, 64)
	patchV, errP := strconv.ParseInt(*patch, 10, 64)
	if errM != nil || errMi != nil || errP != nil {
		return false, fmt.Errorf("we prefer numbers for major, minor and patch values, try again")
	}
	version := structs.Version{
		Major: majorV,
		Minor: minorV,
		Patch: patchV,
	}
	return business.UpversionFile(userId, root, id, &version, file)
}

func UpversionFileChunk(userId *int64, directory *string, id *string, major *int64, minor *int64, patch *int64, name *string, lastModified *int64, size *int64, dataType *string, start *int64, chunk *string) (success bool, startRes int64, fileMetadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "name": name, "size": size, "root": directory, "id": id, "major": major, "minor": minor, "patch": patch, "chunk": chunk})
	if err != nil {
		return false, *start, nil, err
	}
	version := structs.Version{
		Major: *major,
		Minor: *minor,
		Patch: *patch,
	}
	return business.UpversionFileChunk(userId, directory, id, version, name, lastModified, size, dataType, start, chunk, nil)
}
