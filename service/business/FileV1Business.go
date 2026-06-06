package business

import (
	"amper/cache/business"
	"amper/common/argument"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/jsons"
	"amper/data/database"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetDriveDirectory(userId *int64) (result *string, err error) {
	user := business.GetUser(userId, true)
	if user == nil || *user.AmperId == 0 || *user.AmperId != *business.AmperId() {
		return result, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	instance := business.GetAmperInstance(*business.AmperId())
	rootDriveDirectory := instance.Directory
	if rootDriveDirectory == nil || len(*rootDriveDirectory) < 1 {
		return result, fmt.Errorf("empty directory found, user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	result = util.PointerString(filepath.Join(*rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive"))
	errUD := os.MkdirAll(*result, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return result, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", *result)
	}
	return result, nil
}

func getTargetDirectory(driveDirectory *string, directory *string) *string {
	var result *string
	directory = util.PointerString(strings.Trim(*directory, " "))
	directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
	directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
	directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
	result = util.PointerString(filepath.Join(*driveDirectory, *directory))

	return result
}

func Upversion(userId *int64, id *string, newId *string, chunk *[]byte, name *string, Type *string, Size *int64, directory *string) (success bool, fileMetadata *structs.FileMetadata, err error) {
	fileId, errFi := structs.ParseId(id)
	if errFi != nil {
		util.Loggify(errFi)
		return false, nil, fmt.Errorf("failed upversioning, the suplied id is not in correct format, contact the support")

	}
	if fileId.Parent != nil {
		id = fileId.Parent.Format()
	}
	err = argument.Validate(map[string]interface{}{"userId": userId, "id": id, "chunk": chunk, "size": Size, "directory": directory})
	if err != nil {
		return false, nil, err
	}
	if util.EmptyString(name) {
		name = util.UUID()
	}

	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, nil, fmt.Errorf("failed upversioning, user is not allocated to an active directory, contact the support")
	}
	upversionDirectory := getTargetDirectory(driveDirectory, directory)
	var amperDatastoreInstance *structs.Amper

	metadata := &structs.FileMetadata{}
	parentMetadata := &structs.FileMetadata{}
	parentFile := "__file__" + *id
	parentMetadataPath := filepath.Join(*upversionDirectory, parentFile, "metadata")
	if newId == nil || len(*newId) < 1 {
		var errADI error
		instance := business.GetAmperInstance(*business.AmperId())
		amperDatastoreInstance, errADI = GetAvailableAmperDatastore(userId)
		if errADI != nil || amperDatastoreInstance == nil || instance == nil {
			util.Loggify(errADI)
			return false, nil, fmt.Errorf("failed upversioning, not able to allocate amper datastore instance")
		}
		newId = GetFileId(userId, instance.Identifier, amperDatastoreInstance.Identifier, id, nil)

		data, errMet := os.ReadFile(parentMetadataPath)
		if errMet != nil || data == nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("failed upversioning, not able to access the metadata information, file %s is corrupted", *id)
		}
		errP := parentMetadata.Parse(util.PointerString(string(data)))
		if errP != nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("failed upversioning, not able to parse the metadata information, contact the support")
		}

		metadata.Id = newId
		metadata.Name = name
		metadata.Size = *Size
		metadata.LastModified = time.Now().UnixNano()
		metadata.Type = Type
		metadata.RenditionType = util.PointerString("?")

		var latestVersion *structs.FileMetadata
		if parentMetadata.Versions != nil && len(*parentMetadata.Versions) > 0 {
			latestVersion, _ = GetLatestVersion(parentMetadata.Versions)
		} else {
			latestVersion = parentMetadata
		}
		version := structs.Version{
			Major: latestVersion.Version.Major,
			Minor: latestVersion.Version.Minor,
			Patch: latestVersion.Version.Patch,
		}
		version.UpVersion()
		metadata.Version = &version

		var versions []*structs.FileMetadata
		if parentMetadata.Versions != nil {
			versions = *parentMetadata.Versions
		} else {
			versions = make([]*structs.FileMetadata, 0)
		}
		versions = append(versions, metadata)
		parentMetadata.Versions = &versions

		oldMetadataJson, errMJ := parentMetadata.Json()
		if errMJ != nil || oldMetadataJson == nil {
			util.Loggify(errR)
			return false, nil, fmt.Errorf("failed upversioning, not able to initialize the file metadata, try again later or contact support")
		}
		errMet = os.WriteFile(parentMetadataPath, []byte(*oldMetadataJson), 0644)
		if errMet != nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("failed upversioning, not able to write the metadata, try again later or contact the support")
		}
	} else {
		idStruct, errI := structs.ParseId(newId)
		if errI != nil {
			util.Loggify(errI)
			return false, nil, fmt.Errorf("failed upversioning, not able to parse the id, make sure id is correct")
		}
		var errDI error
		amperDatastoreInstance, errDI = database.GetInstance(&idStruct.InstanceId, false)
		if errDI != nil {
			return false, nil, fmt.Errorf("failed upversioning, not able to locate the datastore instance, make sure instanceId is correct")
		}

		data, errMet := os.ReadFile(parentMetadataPath)
		if errMet != nil || data == nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("failed upversioning, not able to access the metadata information, file %s is corrupted", *id)
		}
		errP := parentMetadata.Parse(util.PointerString(string(data)))
		if errP != nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("failed upversioning, not able to parse the metadata information, contact the support")
		}
		if parentMetadata.Versions != nil {
			found := false
			for i := 0; i < len(*parentMetadata.Versions); i++ {
				if (*(*parentMetadata.Versions)[i].Id) == *newId {
					metadata = (*parentMetadata.Versions)[i]
					found = true
					break
				}
			}
			if !found {
				return false, nil, fmt.Errorf("failed upversioning, not able to find the metadata information, contact the support")
			}
		} else {
			return false, nil, fmt.Errorf("failed upversioning, not able to parse the metadata information, contact the support")
		}
	}

	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)
	part, errFF := w.CreateFormFile("value", *name)
	if errFF != nil {
		util.Loggify(errFF)
		return false, nil, fmt.Errorf("not able to prepare file package")
	}
	_, errW := part.Write(*chunk)
	if errW != nil {
		util.Loggify(errW)
		return false, nil, fmt.Errorf("not able to write file package")
	}
	userIdString := strconv.FormatInt(*userId, 10)
	parameters := map[string]*string{
		"userId":    &userIdString,
		"key":       newId,
		"name":      name,
		"type":      Type,
		"size":      util.PointerString(strconv.FormatInt(*Size, 10)),
		"directory": directory,
	}
	errAFF := addFormFields(parameters, w)
	if errAFF != nil {
		util.Loggify(errAFF)
		return false, nil, fmt.Errorf("not able to write form field parameters")
	}
	errC := w.Close()
	if errC != nil {
		util.Loggify(errC)
		return false, nil, fmt.Errorf("not able to close file package")
	}

	url := fmt.Sprintf("http://%s:%s/key-value-store/put", *amperDatastoreInstance.Address, *amperDatastoreInstance.Port)
	req, errNR := http.NewRequest("POST", url, buf)
	if errNR != nil {
		util.Loggify(errNR)
		return false, nil, fmt.Errorf("not able to prepare http post request")
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("userId", strconv.FormatInt(*userId, 10))
	client := &http.Client{}
	res, errDo := client.Do(req)
	if errDo == nil {
		defer res.Body.Close()
		post := &structs.KeyValueStoreResult{}
		errD := json.NewDecoder(res.Body).Decode(post)
		if errD == nil {
			if post.Success {
				if post.Rendition != nil {
					metadata.Thumbnail = post.Rendition.Thumbnail
					metadata.Rendition = post.Rendition.Rendition
					metadata.Processing = post.Rendition.Processing
					metadata.Viewable = post.Rendition.Viewable
					metadata.RenditionType = post.Rendition.RenditionType
					if metadata.Type != nil || len(*metadata.Type) < 1 {
						metadata.Type = post.Rendition.FileType
					}
					parentMetadataJson, errMJ := parentMetadata.Json()
					if errMJ != nil || parentMetadataJson == nil {
						util.Loggify(errR)
						return false, nil, fmt.Errorf("upversion filed, not able to finalize the file metadata, try again later or contact support")
					}
					errMet := os.WriteFile(parentMetadataPath, []byte(*parentMetadataJson), 0644)
					if errMet != nil {
						util.Loggify(errMet)
						return false, nil, fmt.Errorf("upversion filed, not able to write the metadata update, try again later or contact the support")
					}
				}
			} else {
				return false, nil, fmt.Errorf("upversion filed, received failing result whlie storing the data, contact the support")
			}
		} else {
			util.Loggify(errD)
			return false, nil, fmt.Errorf("upversion filed, not able to evaluate service response, contact the support")
		}
	} else {
		util.Loggify(errDo)
		return false, nil, fmt.Errorf("upversion filed, not able to reach the datastore service, contact the support")
	}
	return true, metadata, nil
}

func Upload(userId *int64, id *string, chunk *[]byte, name *string, Type *string, Size *int64, directory *string, optionalValue *string) (success bool, fileMetadata *structs.FileMetadata, err error) {

	if name == nil {
		name = util.UUID()
	}
	if directory == nil {
		directory = util.PointerString(string(os.PathSeparator))
	}
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, nil, fmt.Errorf("user is not allocated to an active directory, contact the support")
	}
	uploadDirectory := getTargetDirectory(driveDirectory, directory)
	errUD := os.MkdirAll(*uploadDirectory, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, nil, fmt.Errorf("not able to initiate the target directory")
	}
	var amperDatastoreInstance *structs.Amper

	metadata := structs.FileMetadata{}
	var progressFile string
	var progressPath string
	var metadataPath string
	if id == nil || len(*id) < 1 {
		var errADI error
		instance := business.GetAmperInstance(*business.AmperId())
		amperDatastoreInstance, errADI = GetAvailableAmperDatastore(userId)
		if errADI != nil || amperDatastoreInstance == nil || instance == nil {
			util.Loggify(errADI)
			return false, nil, fmt.Errorf("not able to allocate amper datastore instance")
		}
		id = GetFileId(userId, instance.Identifier, amperDatastoreInstance.Identifier, nil, optionalValue)

		progressFile = "__progress__" + *id
		progressPath = filepath.Join(*uploadDirectory, progressFile)
		metadataPath = filepath.Join(progressPath, "metadata")

		metadata.Id = id
		metadata.Name = name
		metadata.Size = *Size
		metadata.LastModified = time.Now().UnixNano()
		metadata.Type = Type
		metadata.RenditionType = util.PointerString("?")
		version := structs.Version{
			Major: 0,
			Minor: 0,
			Patch: 1,
		}
		metadata.Version = &version
		metadata.Versions = nil
		metadataJson, errMJ := metadata.Json()
		if errMJ != nil || metadataJson == nil {
			util.Loggify(errR)
			return false, nil, fmt.Errorf("not able to initialize the file metadata, try again later or contact support")
		}
		errMA := os.MkdirAll(progressPath, os.ModePerm)
		if errMA != nil {
			util.Loggify(errMA)
			return false, nil, fmt.Errorf("not able to locate the user's active directory in drive for file '%s', please contect the support", progressPath)
		}
		errMet := os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
		if errMet != nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("not able to write the metadata, try again later or contact the support")
		}
	} else {
		progressFile = "__progress__" + *id
		progressPath = filepath.Join(*uploadDirectory, progressFile)
		metadataPath = filepath.Join(progressPath, "metadata")

		idStruct, errI := structs.ParseId(id)
		if errI != nil {
			util.Loggify(errI)
			return false, nil, fmt.Errorf("not able to parse the id, make sure id is correct")
		}
		var errDI error
		amperDatastoreInstance, errDI = database.GetInstance(&idStruct.InstanceId, false)
		if errDI != nil {
			return false, nil, fmt.Errorf("not able to locate the datastore instance, make sure instanceId is correct")
		}

		data, errMet := os.ReadFile(metadataPath)
		if errMet != nil || data == nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("not able to access the metadata information, file %s is corrupted", *id)
		}
		errP := metadata.Parse(util.PointerString(string(data)))
		if errP != nil {
			util.Loggify(errMet)
			return false, nil, fmt.Errorf("not able to parse the metadata information, contact the support")
		}
	}

	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)
	part, errFF := w.CreateFormFile("value", *name)
	if errFF != nil {
		util.Loggify(errFF)
		return false, nil, fmt.Errorf("not able to prepare file package")
	}
	_, errW := part.Write(*chunk)
	if errW != nil {
		util.Loggify(errW)
		return false, nil, fmt.Errorf("not able to write file package")
	}
	userIdString := strconv.FormatInt(*userId, 10)
	parameters := map[string]*string{
		"userId":    &userIdString,
		"key":       id,
		"name":      name,
		"type":      Type,
		"size":      util.PointerString(strconv.FormatInt(*Size, 10)),
		"directory": directory,
	}
	errAFF := addFormFields(parameters, w)
	if errAFF != nil {
		util.Loggify(errAFF)
		return false, nil, fmt.Errorf("not able to write form field parameters")
	}
	errC := w.Close()
	if errC != nil {
		util.Loggify(errC)
		return false, nil, fmt.Errorf("not able to close file package")
	}

	url := fmt.Sprintf("http://%s:%s/key-value-store/put", *amperDatastoreInstance.Address, *amperDatastoreInstance.Port)
	req, errNR := http.NewRequest("POST", url, buf)
	if errNR != nil {
		util.Loggify(errNR)
		return false, nil, fmt.Errorf("not able to prepare http post request")
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("userId", strconv.FormatInt(*userId, 10))
	client := &http.Client{}
	res, errDo := client.Do(req)
	if errDo == nil {
		defer res.Body.Close()
		post := &structs.KeyValueStoreResult{}
		errD := json.NewDecoder(res.Body).Decode(post)
		if errD == nil {
			if post.Success {
				if post.Rendition != nil {
					metadata.Thumbnail = post.Rendition.Thumbnail
					metadata.Rendition = post.Rendition.Rendition
					metadata.Processing = post.Rendition.Processing
					metadata.Viewable = post.Rendition.Viewable
					metadata.RenditionType = post.Rendition.RenditionType
					if metadata.Type != nil || len(*metadata.Type) < 1 {
						metadata.Type = post.Rendition.FileType
					}
					metadataJson, errMJ := metadata.Json()
					if errMJ != nil || metadataJson == nil {
						util.Loggify(errR)
						return false, nil, fmt.Errorf("not able to finalize the file metadata, try again later or contact support")
					}
					errMet := os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
					if errMet != nil {
						util.Loggify(errMet)
						return false, nil, fmt.Errorf("not able to write the metadata update, try again later or contact the support")
					}

					finalFile := "__file__" + *id
					finalPath := filepath.Join(*uploadDirectory, finalFile)
					errRen := os.Rename(progressPath, finalPath)
					if errRen != nil {
						util.Loggify(errRen)
						return false, nil, fmt.Errorf("not able to finalize the file, contact the support")
					}
				}
			} else {
				return false, nil, fmt.Errorf("received failing result whlie storing the data, contact the support")
			}
		} else {
			util.Loggify(errD)
			return false, nil, fmt.Errorf("not able to evaluate service response, contact the support")
		}
	} else {
		util.Loggify(errDo)
		return false, nil, fmt.Errorf("not able to reach the datastore service, contact the support")
	}
	return true, &metadata, nil
}

func addFormFields(parameters map[string]*string, w *multipart.Writer) error {
	for key, value := range parameters {
		formField, errFF := w.CreateFormField(key)
		if errFF != nil {
			util.Loggify(errFF)
			return fmt.Errorf("not able to add form fields to the multipart request writer")
		}
		_, errW := formField.Write([]byte(*value))
		if errW != nil {
			util.Loggify(errW)
			return fmt.Errorf("not able to write form fields to the multipart request writer")
		}
	}
	return nil
}

func Update(userId *int64, id *string, data *[]byte, directory *string) (success bool, err error) {
	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)
	part, errFF := w.CreateFormFile("value", *id)
	if errFF != nil {
		util.Loggify(errFF)
		return false, fmt.Errorf("update failed, not able to prepare file package")
	}
	_, errW := part.Write(*data)
	if errW != nil {
		util.Loggify(errW)
		return false, fmt.Errorf("update failed, not able to write file package")
	}
	userIdString := strconv.FormatInt(*userId, 10)
	parameters := map[string]*string{
		"userId":    &userIdString,
		"key":       id,
		"directory": directory,
	}
	errAFF := addFormFields(parameters, w)
	if errAFF != nil {
		util.Loggify(errAFF)
		return false, fmt.Errorf("update failed, not able to write form field parameters")
	}
	errC := w.Close()
	if errC != nil {
		util.Loggify(errC)
		return false, fmt.Errorf("update failed, not able to close file package")
	}

	idStruct, errI := structs.ParseId(id)
	if errI != nil {
		util.Loggify(errI)
		return false, fmt.Errorf("update failed, not able to parse the id, make sure id is correct")
	}
	var errDI error
	amperDatastoreInstance, errDI := database.GetInstance(&idStruct.InstanceId, false)
	if errDI != nil {
		return false, fmt.Errorf("update failed, not able to locate the datastore instance, make sure instanceId is correct")
	}

	url := fmt.Sprintf("http://%s:%s/key-value-store/update", *amperDatastoreInstance.Address, *amperDatastoreInstance.Port)
	req, errNR := http.NewRequest("POST", url, buf)
	if errNR != nil {
		util.Loggify(errNR)
		return false, fmt.Errorf("update failed, not able to prepare http post request")
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("userId", strconv.FormatInt(*userId, 10))
	client := &http.Client{}
	res, errDo := client.Do(req)
	if errDo == nil {
		defer res.Body.Close()
		post := &structs.KeyValueStoreResult{}
		errD := json.NewDecoder(res.Body).Decode(post)
		if errD == nil {
			if post.Success {
				return true, nil
			} else {
				return false, fmt.Errorf("update failed, received failing result whlie storing the data, contact the support")
			}
		} else {
			util.Loggify(errD)
			return false, fmt.Errorf("update failed, not able to evaluate service response, contact the support")
		}
	} else {
		util.Loggify(errDo)
		return false, fmt.Errorf("update failed, not able to reach the datastore service, contact the support")
	}
}

func GetFileId(userId *int64, sourceInstanceId *int64, instanceId *int64, parentId *string, optionalValue *string) *string {
	var parent *structs.FileId
	var errP error
	if parentId != nil {
		parent, errP = structs.ParseId(parentId)
		if errP != nil {
			util.Loggify(errP)
		}
	}
	if optionalValue == nil {
		optionalValue = util.PointerString("?")
	}
	now := time.Now()
	result := structs.FileId{
		SourceInstanceId: *sourceInstanceId,
		InstanceId:       *instanceId,
		UserId:           *userId,
		Year:             now.Year(),
		Month:            int(now.Month()),
		Day:              now.Day(),
		UUID:             util.UUID(),
		OptionalValue:    optionalValue,
		Parent:           parent,
	}

	return result.Format()
}

func UpdateMetadata(id *string, directory *string, Thumbnail bool, Rendition bool, RenditionType *string, Viewable bool, Processing bool) (success bool, err error) {
	idStruct, errI := structs.ParseId(id)
	if errI != nil {
		util.Loggify(errI)
		return false, fmt.Errorf("failed updating metadata, not able to parse the id, make sure id is correct")
	}
	targetId := id
	if idStruct.Parent != nil {
		targetId = idStruct.Parent.Format()
	}
	driveDirectory, errR := GetDriveDirectory(&idStruct.UserId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, fmt.Errorf("failed updating metadata, user is not allocated to an active directory, contact the support")
	}
	uploadDirectory := getTargetDirectory(driveDirectory, directory)
	file := "__file__" + *targetId
	metadataPath := filepath.Join(*uploadDirectory, file, "metadata")

	metadata := structs.FileMetadata{}
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return false, fmt.Errorf("failed updating metadata, not able to access the metadata information, file %s is corrupted", *id)
	}
	errP := metadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return false, fmt.Errorf("failed updating metadata, not able to parse the metadata information, contact the support")
	}

	if *metadata.Id == *id {
		metadata.Thumbnail = Thumbnail
		metadata.Rendition = Rendition
		metadata.RenditionType = RenditionType
		metadata.Viewable = Viewable
		metadata.Processing = Processing
	} else {
		if metadata.Versions != nil {
			for i := 0; i < len(*metadata.Versions); i++ {
				metadataVersion := (*metadata.Versions)[i]
				if *metadataVersion.Id == *id {
					metadataVersion.Thumbnail = Thumbnail
					metadataVersion.Rendition = Rendition
					metadataVersion.RenditionType = RenditionType
					metadataVersion.Viewable = Viewable
					metadataVersion.Processing = Processing
				}
			}
		}
	}
	metadataJson, errMJ := metadata.Json()
	if errMJ != nil || metadataJson == nil {
		util.Loggify(errR)
		return false, fmt.Errorf("failed updating metadata, not able to finalize the file metadata, try again later or contact support")
	}
	errMet = os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
	if errMet != nil {
		util.Loggify(errMet)
		return false, fmt.Errorf("failed updating metadata, not able to write the metadata update, try again later or contact the support")
	}
	return true, nil
}

func GetLatestVersion(versions *[]*structs.FileMetadata) (latestMetadataVersion *structs.FileMetadata, availableVersions *[]structs.Version) {
	var tempAvailableVersions []structs.Version
	for i := 0; i < len(*versions); i++ {
		tempAvailableVersions = append(tempAvailableVersions, *(*versions)[i].Version)
		if latestMetadataVersion != nil {
			if (*versions)[i].Version.CompareTo(latestMetadataVersion.Version) > 0 {
				latestMetadataVersion = (*versions)[i]
			}
		} else {
			latestMetadataVersion = (*versions)[i]
		}
	}
	availableVersions = &tempAvailableVersions
	return latestMetadataVersion, availableVersions
}

func FetchFiles(userId *int64, directory *string) (result *[]structs.FileMetadata, err error) {
	files := []structs.FileMetadata{}
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return nil, fmt.Errorf("failed fetching files, user is not allocated to an active directory, contact the support")
	}

	targetDirectory := getTargetDirectory(driveDirectory, directory)
	filesDirectories, errFD := os.ReadDir(*targetDirectory)
	if errFD != nil {
		util.Loggify(errFD)
		return nil, fmt.Errorf("not able to retrieve directories and files at this moment, try again later or contact the support")
	}
	metadataMap := make(map[string]*structs.FileMetadata, 0)
	var thumbnailsToFetch []*string
	for _, file := range filesDirectories {
		if file.IsDir() && strings.HasPrefix(file.Name(), "__file__") {
			var availableVersions []structs.Version
			metadataPath := filepath.Join(*targetDirectory, file.Name(), "metadata")
			metadataJson, errMet := os.ReadFile(metadataPath)
			if errMet != nil || metadataJson == nil {
				util.Loggify(errMet)
				//TODO consider notifying about broken file so users can manage it
				continue
			}
			parentMetadata := structs.FileMetadata{}
			errP := parentMetadata.Parse(util.PointerString(string(metadataJson)))
			if errP != nil {
				util.Loggify(errP)
				//TODO consider notifying about broken file so users can manage it
				continue
			}
			//collect available versions
			availableVersions = append(availableVersions, *parentMetadata.Version)
			var metadata *structs.FileMetadata
			if parentMetadata.Versions != nil && len(*parentMetadata.Versions) > 0 {
				var tempAvailableVersions *[]structs.Version
				metadata, tempAvailableVersions = GetLatestVersion(parentMetadata.Versions)
				availableVersions = append(availableVersions, *tempAvailableVersions...)
			} else {
				metadata = &parentMetadata
			}
			metadata.AvailableVersions = &availableVersions
			if metadata.Thumbnail {
				thumbnailsToFetch = append(thumbnailsToFetch, metadata.Id)
			}
			metadataMap[*metadata.Id] = metadata
		} else if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			directory := structs.FileMetadata{
				Name:  util.PointerString(file.Name()),
				IsDir: true,
			}
			files = append(files, directory)
		}
	}
	//First group the thumbnails by instance id wher the thumbnail file is stored
	thumbnailInstanceGroup := map[int64][]string{}
	for _, id := range thumbnailsToFetch {
		idStruct, errIS := structs.ParseId(id)
		if errIS != nil {
			util.Loggify(errIS)
			continue
		}
		groupIds := thumbnailInstanceGroup[idStruct.InstanceId]
		if groupIds == nil {
			groupIds = []string{}
		}
		groupIds = append(groupIds, *id)
		thumbnailInstanceGroup[idStruct.InstanceId] = groupIds
	}

	//Next reach the datastore instances and collect the thumbnail files
	for instanceId, fileId := range thumbnailInstanceGroup {
		thumbnailMap, errTM := retrieveThumbnailsFromDatastore(userId, instanceId, &fileId)
		if errTM != nil {
			util.Loggify(errTM)
			continue
		}
		for id, thumbnail := range thumbnailMap {
			metadataMap[id].ThumbnailImage = thumbnail
		}
	}
	for _, metadata := range metadataMap {
		files = append(files, *metadata)
	}
	result = &files
	return result, nil
}

func retrieveThumbnailsFromDatastore(userID *int64, instanceId int64, ids *[]string) (result map[string]*string, err error) {
	parametersString, _ := json.Marshal(map[string]interface{}{
		"userId": *userID,
		"ids":    ids,
	})
	body := []byte(parametersString)
	instance := business.GetAmperInstance(instanceId)
	url := fmt.Sprintf("http://%s:%s/%s", *instance.Address, *instance.Port, "key-value-store/tumbnails")
	r, errNR := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if errNR == nil {
		r.Header.Add("userId", strconv.FormatInt(*userID, 10))

		client := &http.Client{}
		res, errDo := client.Do(r)
		if errDo == nil {
			defer res.Body.Close()
			post := &structs.ThumbnailsResult{}
			errD := json.NewDecoder(res.Body).Decode(post)
			if errD == nil {
				if post.Success {
					return post.Data, nil
				} else {
					err = fmt.Errorf("retrieve thumbnails failes, node `%s` received failed internal responce with error: %s", url, post.Error)
				}
			} else {
				util.Loggify(errD)
				err = fmt.Errorf("retrieve thumbnails failes, node `%s` failed: not able to decode the response from the amper instance, try a different address or port", url)
			}
		} else {
			util.Loggify(errDo)
			err = fmt.Errorf("retrieve thumbnails failes, node `%s` failed: not able to reach the amper instance, try a different address or port", url)
		}
	} else {
		util.Loggify(errNR)
		err = fmt.Errorf("retrieve thumbnails failes, node `%s` failed: not able to reach the amper instance, please try again with different address and port or reach the support", url)
	}

	return nil, err
}

func NewDirectory(userId *int64, directory *string, name *string) (result bool, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, fmt.Errorf("failed creating directory, user is not allocated to an active directory, contact the support")
	}

	targetDirectory := getTargetDirectory(driveDirectory, directory)

	if name == nil || strings.HasPrefix(*name, ".") || strings.HasPrefix(*name, "__progress__") || strings.HasPrefix(*name, "__file__") || strings.HasPrefix(*name, "__system__") {
		return false, fmt.Errorf("a folder name is required and can't start with '.', '__progress__', or '__file__'")
	}
	pattern := regexp.MustCompile(`^[^\s^\x00-\x1f\\?*:"";<>|\/.][^\x00-\x1f\\?*:"";<>|\/]*[^\s^\x00-\x1f\\?*:"";<>|\/.]+$`)
	found := pattern.FindString(*name)
	if found == "" {
		return false, fmt.Errorf("a folder name is required and can't start with '.', ' ', '__progress__', or '__file__' and contain characters '*, \\, :, \", /, >, <, ?, |'")
	}

	newDirectory := filepath.Join(*targetDirectory, *name)
	if _, err := os.Stat(newDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(newDirectory, os.ModePerm)
		if err != nil {
			return false, fmt.Errorf("not able to make up a new directory with supplied naming '%s'", *name)
		}
	} else {
		return false, fmt.Errorf("a directory with a name '%s' already exists, try something different", *name)
	}

	return true, nil
}

func RemoveFile(userId *int64, directory *string, id *string) (result bool, err error) {
	return RemoveFiles(userId, directory, &[]string{*id})
}

func RemoveDirectory(userId *int64, root *string, directory *string) (success bool, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, fmt.Errorf("failed removing directory, user is not allocated to an active directory, contact the support")
	}
	success = true
	targetDirectory := getTargetDirectory(driveDirectory, root)

	fileIds := CollectFileIdsRecursively(filepath.Join(*targetDirectory, *directory))

	successRF, errRF := RemoveFiles(userId, root, &fileIds)
	if errRF == nil && successRF {
		errR := os.RemoveAll(filepath.Join(*targetDirectory, *directory))
		if errR != nil {
			success = false
			util.Loggify(errR)
			err = fmt.Errorf("remove directory failed, not able to remove all of the file or folders at this time, contact the support")
		}
	} else {
		util.Loggify(errRF)
		success = false
		err = fmt.Errorf("remove directory files failed, not able to remove all of the file or folders at this time, contact the support")
	}
	return success, err
}

func CollectFileIdsRecursively(directory string) (result []string) {
	directories, errFD := os.ReadDir(directory)
	if errFD != nil || directories == nil || len(directories) < 1 {
		util.Loggify(errFD)
		return
	}
	for _, file := range directories {
		if file.IsDir() && (strings.HasPrefix(file.Name(), "__file__") || strings.HasPrefix(file.Name(), "__progress__")) {
			metadataPath := filepath.Join(directory, file.Name(), "metadata")
			metadataJson, errMet := os.ReadFile(metadataPath)
			if errMet != nil || metadataJson == nil {
				util.Loggify(errMet)
				//TODO consider notifying about broken file so users can manage it
				continue
			}
			parentMetadata := structs.FileMetadata{}
			errP := parentMetadata.Parse(util.PointerString(string(metadataJson)))
			if errP != nil {
				util.Loggify(errP)
				//TODO consider notifying about broken file so users can manage it
				continue
			}
			result = append(result, *parentMetadata.Id)
		} else if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			result = append(result, CollectFileIdsRecursively(filepath.Join(directory, file.Name()))...)
		}
	}
	return result
}

func RemoveFiles(userId *int64, root *string, ids *[]string) (result bool, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, fmt.Errorf("failed removing files, user is not allocated to an active directory, contact the support")
	}

	targetDirectory := getTargetDirectory(driveDirectory, root)
	result = true

	for _, id := range *ids {
		idStruct, errIS := structs.ParseId(&id)
		if errIS != nil {
			//If parsing fails, this implies the id is a directory
			//Recursively delete all files and irectories containing the directory if exists
			partialsuccess, errRD := RemoveDirectory(userId, root, &id)
			if errRD != nil {
				util.Loggify(errRD)
			}
			if !partialsuccess {
				result = false
			}
			continue
		}
		var versionIdsToRemove []string
		var parentId *string
		if idStruct.Parent != nil {
			parentId = idStruct.Parent.Format()
			versionIdsToRemove = append(versionIdsToRemove, *parentId)

			metadataPath := filepath.Join(*targetDirectory, "__file__"+*parentId, "metadata")
			metadataJson, errMet := os.ReadFile(metadataPath)
			if errMet != nil || metadataJson == nil {
				util.Loggify(errMet)
				continue
			}
			parentMetadata := structs.FileMetadata{}
			errP := parentMetadata.Parse(util.PointerString(string(metadataJson)))
			if errP != nil {
				util.Loggify(errP)
				//TODO consider notifying about broken file so users can manage it
				continue
			}
			if parentMetadata.Versions != nil && len(*parentMetadata.Versions) > 0 {
				for _, metadataVersion := range *parentMetadata.Versions {
					versionIdsToRemove = append(versionIdsToRemove, *metadataVersion.Id)
				}
			}
		} else {
			versionIdsToRemove = append(versionIdsToRemove, id)
			parentId = &id
		}

		for _, metadataVersionId := range versionIdsToRemove {
			versionIdStruct, errVIS := structs.ParseId(&metadataVersionId)
			if errVIS != nil {
				util.Loggify(errVIS)
				continue
			}

			instance := business.GetAmperInstance(versionIdStruct.InstanceId)
			partialSuccess, _, errPS := DedicatedCallWithRetry(userId, nil, map[string]string{
				"amperDatastoreInstance": "key-value-store/delete",
			}, map[string]interface{}{
				"key": metadataVersionId,
			}, instance)
			if errPS != nil || !partialSuccess {
				util.Loggify(errPS)
				result = false
				continue
			}
		}
		removeFilePath := filepath.Join(*targetDirectory, "__file__"+*parentId)
		errR := os.RemoveAll(removeFilePath)
		if errR != nil {
			result = false
			util.Loggify(errR)
			err = fmt.Errorf("not able to remove all the file or folders at this time, contact the support")
		}
	}

	return result, err
}

func MoveFiles(userId *int64, root *string, ids *[]string, directory *string) (result bool, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, fmt.Errorf("failed moving directory, user is not allocated to an active directory, contact the support")
	}

	targetDirectory := getTargetDirectory(driveDirectory, root)

	result = true
	for _, id := range *ids {
		var fileName *string
		if *directory == id {
			continue
		}
		_, errFI := structs.ParseId(&id)
		//if parse failes, implies id is a directory
		if errFI != nil {
			fileName = &id
		} else {
			fileName = util.PointerString("__file__" + id)
		}

		originalFilePath := filepath.Join(*targetDirectory, *fileName)
		moveFilePath := filepath.Join(*targetDirectory, *directory, *fileName)

		errRen := os.Rename(originalFilePath, moveFilePath)
		if errRen != nil {
			util.Loggify(errRen)
			result = false
			err = fmt.Errorf("not able to move all files at this time, contact the support")
		}
	}

	return result, err
}

func PasteFiles(userId *int64, root *string, copyJson *string, cutJson *string) (result bool, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return false, fmt.Errorf("failed pasting directory, user is not allocated to an active directory, contact the support")
	}

	targetDirectory := getTargetDirectory(driveDirectory, root)

	/*copy, errC := jsons.GetJsonObject(copyJson)
	if errC != nil {
		util.Loggify(errC)
		return false, fmt.Errorf("the supplied copy parameter: %s is of wrong format", *copyJson)
	}*/
	cut, errCut := jsons.GetJsonObject(cutJson)
	if errCut != nil {
		util.Loggify(errCut)
		return false, fmt.Errorf("the supplied cut parameter: %s is of wrong format", *cutJson)
	}
	result = true
	if len(cut) < 1 {
		return result, err
	}
	resultError := ""
	/*for file, directoryInt := range copy {
		if len(file) > 0 && reflect.TypeOf(directoryInt).String() == "string" {
			directory := directoryInt.(string)
			if len(directory) > 0 {
				fullFilePath := getFullPath(*driveDirectory, directory, file)
				if _, errS := os.Stat(fullFilePath); errors.Is(errS, os.ErrNotExist) {
					result = false
					resultError = resultError + ", " + fmt.Sprintf("Not able to find the file '%s' in the specified directory '%s'", file, directory)
				} else {
					fullTargetFilePath := filepath.Join(*targetDirectory, file)
					_, errS1 := os.Stat(fullTargetFilePath)
					index := int64(0)
					for !errors.Is(errS1, os.ErrNotExist) {
						index++
						_, errS1 = os.Stat(fullTargetFilePath + "(" + strconv.FormatInt(index, 10) + ")")
					}
					if index > 0 {
						fullTargetFilePath = fullTargetFilePath + "(" + strconv.FormatInt(index, 10) + ")"
					}
					errDest := os.MkdirAll(fullTargetFilePath, os.ModePerm)
					if errDest == nil {
						errCopy := CopyDirectory(fullFilePath, fullTargetFilePath)
						if errCopy != nil {
							util.Loggify(errCopy)
							resultError = resultError + ", " + fmt.Sprintf("Not able to copy the file '%s' to the specified directory '%s'", file, *root)
						}
					} else {
						util.Loggify(errDest)
						result = false
						resultError = resultError + ", " + fmt.Sprintf("Not able to copy the file '%s' to the specified directory '%s' due to missing directory", file, *root)
					}
				}
			}
		}
	}*/

	for file, directoryInt := range cut {
		if len(file) > 0 && reflect.TypeOf(directoryInt).String() == "string" {
			var fileName *string
			_, errFI := structs.ParseId(&file)
			//if parse failes, implies id is a directory
			if errFI != nil {
				fileName = &file
			} else {
				fileName = util.PointerString("__file__" + file)
			}
			directory := directoryInt.(string)
			if len(directory) > 0 {
				fullFilePath := getFullPath(*driveDirectory, directory, *fileName)
				if _, errS := os.Stat(fullFilePath); errors.Is(errS, os.ErrNotExist) {
					result = false
					resultError = resultError + ", " + fmt.Sprintf("Not able to find the file '%s' in the specified directory '%s'", file, directory)
				} else {
					fullTargetFilePath := filepath.Join(*targetDirectory, *fileName)
					if fullFilePath != fullTargetFilePath {
						errMove := os.Rename(fullFilePath, fullTargetFilePath)
						if errMove != nil {
							util.Loggify(errMove)
							result = false
							resultError = resultError + ", " + fmt.Sprintf("Not able to move the file '%s' to the specified directory '%s'", *fileName, *root)
						}
					}
				}
			}
		}
	}
	if len(resultError) > 0 {
		result = false
		err = fmt.Errorf(resultError)
	}
	return result, err
}

func DiscoverRoot(userId *int64) (result *structs.Folder, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return nil, fmt.Errorf("failed creating directory, user is not allocated to an active directory, contact the support")
	}

	result = &structs.Folder{
		Name:    util.PointerString("Root"),
		Path:    util.PointerString(string(os.PathSeparator)),
		Folders: RootWalk(*driveDirectory, string(os.PathSeparator)),
	}
	return result, err
}

func RootWalk(root string, relativePath string) *[]structs.Folder {
	directories, errFD := os.ReadDir(root)
	if errFD != nil || directories == nil || len(directories) < 1 {
		util.Loggify(errFD)
		return nil
	}
	result := make([]structs.Folder, 0)
	for _, file := range directories {
		if file.IsDir() && !strings.HasPrefix(file.Name(), "__file__") && !strings.HasPrefix(file.Name(), "__progress__") {
			folder := structs.Folder{}
			folder.Name = util.PointerString(file.Name())
			folder.Path = util.PointerString(filepath.Join(relativePath, file.Name()))
			folder.Folders = RootWalkOld(filepath.Join(root, file.Name()), *folder.Path)
			result = append(result, folder)
		}
	}
	return &result
}

func GetFileBody(userId *int64, id *string) (result *io.ReadCloser, err error) {
	idStruct, errIS := structs.ParseId(id)
	if errIS != nil {
		util.Loggify(errIS)
		return nil, fmt.Errorf("failed fetching file body, id is not of a valid format")
	}

	instanceId := idStruct.InstanceId
	amperDatastoreInstance := business.GetAmperInstance(instanceId)

	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)
	userIdString := strconv.FormatInt(*userId, 10)

	parameters := map[string]interface{}{
		"userId":    &userIdString,
		"key":       id,
		"rendition": false,
	}
	parametersString, _ := json.Marshal(parameters)
	body := []byte(parametersString)

	url := fmt.Sprintf("http://%s:%s/key-value-store/get", *amperDatastoreInstance.Address, *amperDatastoreInstance.Port)
	req, errNR := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if errNR != nil {
		util.Loggify(errNR)
		return nil, fmt.Errorf("failed accessing the file in key value store, not able to prepare http post request")
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("userId", strconv.FormatInt(*userId, 10))
	client := &http.Client{}
	res, errDo := client.Do(req)
	if errDo == nil {
		return &res.Body, nil
	} else {
		util.Loggify(errDo)
		return nil, fmt.Errorf("failed accessing the file on key value store, not able to reach the datastore service, contact the support")
	}
}

func GetFile(userId *int64, directory *string, id *string, version *structs.Version, rendition *bool) (result *io.ReadCloser, metadata *structs.FileMetadata, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return nil, nil, fmt.Errorf("failed fetching metadata, user is not allocated to an active directory, contact the support")
	}

	targetDirectory := getTargetDirectory(driveDirectory, directory)

	idStruct, errIS := structs.ParseId(id)
	if errIS != nil {
		util.Loggify(errIS)
		return nil, nil, fmt.Errorf("failed fetching metadata, id is not of a valid format")
	}
	targetId := id
	if idStruct.Parent != nil {
		targetId = idStruct.Parent.Format()
	}

	file := "__file__" + *targetId
	metadataPath := filepath.Join(*targetDirectory, file, "metadata")

	originalMetadata := structs.FileMetadata{}
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return nil, nil, fmt.Errorf("failed fetching metadata, not able to access the metadata information, file %s is corrupted", *id)
	}
	errP := originalMetadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return nil, nil, fmt.Errorf("failed fetching metadata, not able to parse the metadata information, contact the support")
	}

	var availableVersion []structs.Version
	availableVersion = append(availableVersion, *originalMetadata.Version)
	if originalMetadata.Version.CompareTo(version) == 0 {
		metadata = &originalMetadata
	} else {
		if originalMetadata.Versions != nil {
			for _, metadataVersion := range *originalMetadata.Versions {
				if metadataVersion.Version.CompareTo(version) == 0 {
					metadata = metadataVersion
					break
				}
			}
		}
	}
	if metadata == nil {
		return nil, nil, fmt.Errorf("failed fetching metadata, not able to find the metadata version, contact the support")
	}
	for _, metadataVersion := range *originalMetadata.Versions {
		availableVersion = append(availableVersion, *metadataVersion.Version)
	}
	metadata.AvailableVersions = &availableVersion
	instanceId := idStruct.InstanceId
	amperDatastoreInstance := business.GetAmperInstance(instanceId)

	var buf = new(bytes.Buffer)
	var w = multipart.NewWriter(buf)
	userIdString := strconv.FormatInt(*userId, 10)
	renditionOverride := metadata.Rendition
	if rendition != nil {
		renditionOverride = *rendition
	}
	parameters := map[string]interface{}{
		"userId":    &userIdString,
		"key":       id,
		"rendition": renditionOverride,
	}
	parametersString, _ := json.Marshal(parameters)
	body := []byte(parametersString)

	url := fmt.Sprintf("http://%s:%s/key-value-store/get", *amperDatastoreInstance.Address, *amperDatastoreInstance.Port)
	req, errNR := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if errNR != nil {
		util.Loggify(errNR)
		return nil, nil, fmt.Errorf("failed viewing the file, not able to prepare http post request")
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	req.Header.Add("userId", strconv.FormatInt(*userId, 10))
	client := &http.Client{}
	res, errDo := client.Do(req)
	if errDo == nil {
		return &res.Body, metadata, nil
		/*post := &structs.KeyValueStoreResult{}
		errD := json.NewDecoder(res.Body).Decode(post)
		if errD == nil {
			if post.Success {
			} else {
				return false, nil, fmt.Errorf("upversion filed, received failing result whlie storing the data, contact the support")
			}
		} else {
			util.Loggify(errD)
			return false, nil, fmt.Errorf("upversion filed, not able to evaluate service response, contact the support")
		}*/
	} else {
		util.Loggify(errDo)
		return nil, nil, fmt.Errorf("failed viewing the file, not able to reach the datastore service, contact the support")
	}
}

func GetMetadata(userId *int64, directory *string, id *string, version *structs.Version) (metadata *structs.FileMetadata, err error) {
	driveDirectory, errR := GetDriveDirectory(userId)
	if driveDirectory == nil || errR != nil {
		util.Loggify(errR)
		return nil, fmt.Errorf("failed fetching the metadata, user is not allocated to an active directory, contact the support")
	}

	targetDirectory := getTargetDirectory(driveDirectory, directory)

	idStruct, errIS := structs.ParseId(id)
	if errIS != nil {
		util.Loggify(errIS)
		return nil, fmt.Errorf("failed fetching the metadata, id is not of a valid format")
	}
	targetId := id
	if idStruct.Parent != nil {
		targetId = idStruct.Parent.Format()
	}

	file := "__file__" + *targetId
	metadataPath := filepath.Join(*targetDirectory, file, "metadata")
	originalMetadata := structs.FileMetadata{}
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return nil, fmt.Errorf("failed fetching the metadata, not able to access the metadata information, file %s is corrupted", *id)
	}
	errP := originalMetadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return nil, fmt.Errorf("failed fetching the metadata, not able to parse the metadata information, contact the support")
	}

	var availableVersion []structs.Version
	availableVersion = append(availableVersion, *originalMetadata.Version)
	if originalMetadata.Version.CompareTo(version) == 0 {
		metadata = &originalMetadata
	} else {
		if originalMetadata.Versions != nil {
			for _, metadataVersion := range *originalMetadata.Versions {
				if metadataVersion.Version.CompareTo(version) == 0 {
					metadata = metadataVersion
					break
				}
			}
		}
	}
	for _, metadataVersion := range *originalMetadata.Versions {
		availableVersion = append(availableVersion, *metadataVersion.Version)
	}
	metadata.AvailableVersions = &availableVersion
	return metadata, nil
}
