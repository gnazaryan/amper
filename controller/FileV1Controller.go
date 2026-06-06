package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"amper/common/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// FileV1Controller is responsible for dispatching requests related to
// files managment functionalities
func FileV1Controller(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "upload":
			resultStruct = upload(userID, w, r)
		case "upversion":
			resultStruct = upversion(userID, w, r)
		case "upversionFull":
			resultStruct = upversionFull(userID, w, r)
		case "updateMetadata":
			resultStruct = updateMetadata(userID, w, r)
		case "fetch":
			resultStruct = fetchFiles(userID, w, r)
		case "newDir":
			resultStruct = newDirectory(userID, w, r)
		case "remove":
			resultStruct = removeFile(userID, w, r)
		case "removeFiles":
			resultStruct = removeFiles(userID, w, r)
		case "removeDirectory":
			resultStruct = removeDirectory(userID, w, r)
		case "moveFiles":
			resultStruct = moveFiles(userID, w, r)
		case "pasteFiles":
			resultStruct = pasteFiles(userID, w, r)
		case "discover":
			resultStruct = discoverRoot(userID, w, r)
		case "viewFile":
			viewFile(userID, w, r)
		case "download":
			downloadFile(userID, w, r)
		case "metadata":
			resultStruct = metadata(userID, w, r)
			return
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func upload(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChunkResult) {
	file, _, errF := r.FormFile("chunk")
	if errF != nil || file == nil {
		util.Loggify(errF)
		result.Success = false
		result.Error = "no blob file recevied through http post method 'chunk' parameter"
		return
	}
	defer file.Close()
	var buf bytes.Buffer
	io.Copy(&buf, file)
	chunk := buf.Bytes()
	id := r.FormValue("id")
	name := r.FormValue("name")
	fType := r.FormValue("type")
	sizeString := r.FormValue("size")
	size, errS := strconv.ParseInt(sizeString, 10, 64)
	if errS != nil {
		util.Loggify(errS)
		result.Success = false
		result.Error = "size is a requiered parameter, make sure it is sent along form data"
		return
	}
	directory := r.FormValue("directory")
	success, metadata, errU := authorization.Upload(userID, &id, &chunk, &name, &fType, &size, &directory)
	if errU == nil {
		result.Success = success
		result.Metadata = metadata
	} else {
		result.Success = false
		result.Error = errU.Error()
	}
	return result
}

func upversion(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChunkResult) {
	file, _, errF := r.FormFile("chunk")
	if errF != nil || file == nil {
		util.Loggify(errF)
		result.Success = false
		result.Error = "no blob file recevied through http post method 'chunk' parameter"
		return
	}
	defer file.Close()
	var buf bytes.Buffer
	io.Copy(&buf, file)
	chunk := buf.Bytes()
	id := r.FormValue("id")
	newId := r.FormValue("newId")
	major := r.FormValue("major")
	minor := r.FormValue("minor")
	patch := r.FormValue("patch")
	name := r.FormValue("name")
	fType := r.FormValue("type")
	sizeString := r.FormValue("size")
	size, errS := strconv.ParseInt(sizeString, 10, 64)
	if errS != nil {
		util.Loggify(errS)
		size = int64(len(chunk))
	}
	directory := r.FormValue("directory")
	success, metadata, errU := authorization.Upversion(userID, &id, &newId, &major, &minor, &patch, &chunk, &name, &fType, &size, &directory)
	if errU == nil {
		result.Success = success
		result.Metadata = metadata
	} else {
		result.Success = false
		result.Error = errU.Error()
	}
	return result
}

func upversionFull(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
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
	size := int64(len(file))
	success, _, errU := authorization.Upversion(userID, &id, nil, &major, &minor, &patch, &file, nil, nil, &size, &root)
	if errU == nil {
		result.Success = success
	} else {
		result.Success = false
		result.Error = errU.Error()
	}
	return result
}

func updateMetadata(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Id            *string
		Directory     *string
		Thumbnail     bool
		Rendition     bool
		Processing    bool
		RenditionType *string
		Viewable      bool
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errUM := authorization.UpdateMetadata(parameters.Id, parameters.Directory, parameters.Thumbnail, parameters.Rendition, parameters.RenditionType, parameters.Viewable, parameters.Processing)
	if errUM == nil {
		result.Success = success
	} else {
		result.Success = false
		result.Error = errUM.Error()
	}
	return result
}

func fetchFiles(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
	var parameters struct {
		Root *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	files, err := authorization.FetchFiles(userID, parameters.Root)

	if err == nil {
		result.Success = true
		result.Data = files
	} else {
		result.Error = err.Error()
	}
	return result
}

func newDirectory(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
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

func removeFile(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
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

func removeFiles(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
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

func removeDirectory(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Files) {
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

func moveFiles(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
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

func pasteFiles(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
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

func discoverRoot(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.FolderResult) {
	data, err := authorization.DiscoverRoot(userID)
	if err == nil {
		result.Success = true
		result.Data = data
	} else {
		result.Error = err.Error()
	}
	return result
}

func viewFile(userID *int64, w *http.ResponseWriter, r *http.Request) {
	Directory := r.URL.Query().Get("root")
	Id := r.URL.Query().Get("id")
	Major := r.URL.Query().Get("major")
	Minor := r.URL.Query().Get("minor")
	Patch := r.URL.Query().Get("patch")
	reader, metadata, err := authorization.GetFile(userID, &Directory, &Id, &Major, &Minor, &Patch, nil)

	if err == nil && reader != nil {
		(*w).Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", *metadata.Name))
		(*w).Header().Set("Content-Type", util.IfElse(metadata.Rendition, *metadata.RenditionType, *metadata.Type).(string))
		io.Copy((*w), *reader)
		defer (*reader).Close()
	}
}

func downloadFile(userID *int64, w *http.ResponseWriter, r *http.Request) {
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

func metadata(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.MetadataResult) {
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
