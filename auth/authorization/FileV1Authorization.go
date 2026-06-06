package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
	"io"
	"strconv"
)

func Upload(userId *int64, id *string, chunk *[]byte, name *string, Type *string, Size *int64, directory *string) (success bool, fileMetadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "chunk": chunk, "size": Size})
	if err != nil {
		return false, nil, err
	}
	return business.Upload(userId, id, chunk, name, Type, Size, directory, nil)
}

func Upversion(userId *int64, id *string, newId *string, major *string, minor *string, patch *string, chunk *[]byte, name *string, Type *string, Size *int64, directory *string) (success bool, fileMetadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "id": id, "chunk": chunk, "size": Size, "directory": directory})
	if err != nil {
		return false, nil, err
	}
	return business.Upversion(userId, id, newId, chunk, name, Type, Size, directory)
}

func UpdateMetadata(id *string, directory *string, Thumbnail bool, Rendition bool, RenditionType *string, Viewable bool, Processing bool) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"id": id, "directory": directory})
	if err != nil {
		return false, err
	}
	return business.UpdateMetadata(id, directory, Thumbnail, Rendition, RenditionType, Viewable, Processing)
}

func FetchFiles(userId *int64, directory *string) (files *[]structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId})
	if err != nil {
		return nil, err
	}
	return business.FetchFiles(userId, directory)
}

func NewDirectory(userId *int64, directory *string, name *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "name": name})
	if err != nil {
		return false, err
	}
	return business.NewDirectory(userId, directory, name)
}

func RemoveFile(userId *int64, directory *string, id *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "directory": directory, "id": id})
	if err != nil {
		return false, err
	}
	return business.RemoveFile(userId, directory, id)
}

func RemoveFiles(userId *int64, root *string, ids *[]string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "root": root, "ids": ids})
	if err != nil {
		return false, err
	}
	return business.RemoveFiles(userId, root, ids)
}

func MoveFiles(userId *int64, root *string, ids *[]string, directory *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "root": root, "ids": ids, "directory": directory})
	if err != nil {
		return false, err
	}
	return business.MoveFiles(userId, root, ids, directory)
}

func PasteFiles(userId *int64, root *string, copy *string, cut *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "root": root})
	if err != nil {
		return false, err
	}
	return business.PasteFiles(userId, root, copy, cut)
}

func DiscoverRoot(userId *int64) (result *structs.Folder, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId})
	if err != nil {
		return nil, err
	}
	return business.DiscoverRoot(userId)
}

func GetFile(userId *int64, directory *string, id *string, major *string, minor *string, patch *string, rendition *bool) (result *io.ReadCloser, metadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "directory": directory, "id": id})
	if err != nil {
		return nil, nil, err
	}
	var version *structs.Version
	majorV, errM := strconv.ParseInt(*major, 10, 64)
	minorV, errMi := strconv.ParseInt(*minor, 10, 64)
	patchV, errP := strconv.ParseInt(*patch, 10, 64)
	if errM == nil && errMi == nil && errP == nil {
		version = &structs.Version{
			Major: majorV,
			Minor: minorV,
			Patch: patchV,
		}
	}
	return business.GetFile(userId, directory, id, version, rendition)
}

func GetMetadata(userId *int64, directory *string, id *string, major *int64, minor *int64, patch *int64) (metadata *structs.FileMetadata, err error) {
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
