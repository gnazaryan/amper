package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"amper/common/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// FileController is responsible for dispatching requests related to
// files managment functionalities
func FileController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "uploadChunk":
			resultStruct = uploadChunk(userID, w, r)
		case "upversion":
			resultStruct = upversionFile(userID, w, r)
		case "upversionChunk":
			resultStruct = upversionFileChunk(userID, w, r)
		case "fetch":
			resultStruct = fetchFiles(userID, w, r)
		case "newDir":
			resultStruct = newDirectory(userID, w, r)
		case "remove":
			resultStruct = removeFile(userID, w, r)
		case "removeDirectory":
			resultStruct = removeDirectory(userID, w, r)
		case "removeFiles":
			resultStruct = removeFiles(userID, w, r)
		case "moveFiles":
			resultStruct = moveFiles(userID, w, r)
		case "pasteFiles":
			resultStruct = pasteFiles(userID, w, r)
		case "discover":
			resultStruct = discoverRoot(userID, w, r)
		case "viewFile":
			viewFile(userID, w, r)
		case "metadata":
			resultStruct = metadata(userID, w, r)
		case "download":
			downloadFile(userID, w, r)
			return
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func upversionFileChunk(userId *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChunkResult) {
	var parameters struct {
		Name         *string
		LastModified *int64
		Size         *int64
		Type         *string
		Chunk        *string
		Start        *int64
		Directory    *string
		Id           *string
		Major        *int64
		Minor        *int64
		Patch        *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, startRes, metadata, err := authorization.UpversionFileChunk(userId, parameters.Directory, parameters.Id, parameters.Major, parameters.Minor, parameters.Patch,
		parameters.Name, parameters.LastModified, parameters.Size, parameters.Type, parameters.Start, parameters.Chunk)

	result.Start = startRes
	if err == nil {
		result.Metadata = metadata
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return result
}

func upversionFile(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	root := r.URL.Query().Get("root")
	id := r.URL.Query().Get("id")
	major := r.URL.Query().Get("major")
	minor := r.URL.Query().Get("minor")
	patch := r.URL.Query().Get("patch")
	file, errBR := ioutil.ReadAll(r.Body)
	if errBR != nil {
		result.Error = "we prefer to receive file content with the upversion request on the request body"
		result.Success = false
		return
	}
	success, err := authorization.UpversionFile(userID, &root, &id, &major, &minor, &patch, &file)
	if err == nil && success {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func pasteFilesOld(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Root *string
		Copy *string
		Cut  *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.PasteFiles(userID, parameters.Root, parameters.Copy, parameters.Cut)

	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func moveFilesOld(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Root      *string
		Ids       *[]string
		Directory *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.MoveFiles(userID, parameters.Root, parameters.Ids, parameters.Directory)

	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func removeDirectoryOld(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
	var parameters struct {
		Root *string
		Id   *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.RemoveDirectory(userID, parameters.Root, parameters.Id)

	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func removeFilesOld(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
	var parameters struct {
		Root *string
		Ids  *[]string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.RemoveFiles(userID, parameters.Root, parameters.Ids)

	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func removeFileOld(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
	var parameters struct {
		Root *string
		Id   *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.RemoveFile(userID, parameters.Root, parameters.Id)

	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func downloadFileOld(userID *int64, w *http.ResponseWriter, r *http.Request) {
	Directory := r.URL.Query().Get("root")
	Id := r.URL.Query().Get("id")
	RenditionIndicator := r.URL.Query().Get("rendition")
	rendition, _ := strconv.ParseBool(RenditionIndicator)
	Major := r.URL.Query().Get("major")
	Minor := r.URL.Query().Get("minor")
	Patch := r.URL.Query().Get("patch")
	file, metadata, err := authorization.GetFile(userID, &Directory, &Id, &Major, &Minor, &Patch, &rendition)

	if err == nil && file != nil {
		fileName := *metadata.Name
		if rendition {
			if metadata.Rendition && metadata.RenditionType != nil && *metadata.RenditionType != "?" {
				switch *metadata.RenditionType {
				case "application/pdf":
					fileName = metadata.FileNameWithoutExt() + ".pdf"
				case "image/jpeg":
					fileName = metadata.FileNameWithoutExt() + ".jpeg"
				case "image/png":
					fileName = metadata.FileNameWithoutExt() + ".png"
				}
			}
		}
		(*w).Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		(*w).Header().Set("Content-Type", "application/octet-stream")
		io.Copy((*w), *file)
	}
}

func viewFileOld(userID *int64, w *http.ResponseWriter, r *http.Request) {
	Directory := r.URL.Query().Get("root")
	Id := r.URL.Query().Get("id")
	Major := r.URL.Query().Get("major")
	Minor := r.URL.Query().Get("minor")
	Patch := r.URL.Query().Get("patch")
	file, metadata, err := authorization.GetFile(userID, &Directory, &Id, &Major, &Minor, &Patch, nil)

	if err == nil && file != nil {
		(*w).Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", *metadata.Name))
		(*w).Header().Set("Content-Type", util.IfElse(metadata.Rendition, *metadata.RenditionType, *metadata.Type).(string))
		io.Copy((*w), *file)
	}
}

func metadataOLD(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.MetadataResult) {
	var parameters struct {
		Directory *string
		Id        *string
		Major     *int64
		Minor     *int64
		Patch     *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	metadata, err := authorization.GetMetadata(userID, parameters.Directory, parameters.Id, parameters.Major, parameters.Minor, parameters.Patch)

	if err == nil {
		result.Success = true
		result.Data = metadata
	} else {
		result.Error = err.Error()
	}
	return result
}

func newDirectoryOld(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
	var parameters struct {
		Root *string
		Name *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.NewDirectory(userID, parameters.Root, parameters.Name)

	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func discoverRootOld(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.FolderResult) {
	data, err := authorization.DiscoverRoot(userID)
	if err == nil {
		result.Success = true
		result.Data = data
	} else {
		result.Error = err.Error()
	}
	return result
}

func uploadChunk(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChunkResult) {
	var parameters struct {
		Name         *string
		LastModified *int64
		Size         *int64
		Type         *string
		Chunk        *string
		Start        *int64
		Directory    *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, startRes, metadata, err := authorization.UploadChunk(userID, parameters.Name, parameters.LastModified, parameters.Size, parameters.Type, parameters.Chunk, parameters.Start, parameters.Directory)
	result.Start = startRes
	if err == nil {
		result.Metadata = metadata
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return result
}
