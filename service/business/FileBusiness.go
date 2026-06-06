package business

import (
	"amper/cache/business"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/ampstrings"
	"amper/common/util/files"
	"amper/common/util/jsons"
	"amper/data/database"
	"amper/service/processor/rendition"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func UploadChunk(userId *int64, name *string, lastModified *int64, size *int64, Type *string, chunk *string, chunkByte *[]byte, start *int64, directory *string) (success bool, startRes int64, fileMetadata *structs.FileMetadata, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, -1, nil, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, -1, nil, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, -1, nil, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	uploadPath := drivePath
	if directory != nil {
		directory = util.PointerString(strings.Trim(*directory, " "))
		directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
		directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
		directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
		uploadPath = filepath.Join(drivePath, *directory)
	}
	if lastModified == nil {
		currentTime := time.Now().UnixNano() / int64(time.Millisecond)
		lastModified = &currentTime
	}

	extension := filepath.Ext(*name)
	tempFile := files.Name(util.PointerString("__progress__" + base64.StdEncoding.EncodeToString([]byte(extension+"|"+strconv.FormatInt(int64(*size), 10)+"|"+strconv.FormatInt(*lastModified, 10)+"|"+*Type))))
	progressPath := filepath.Join(uploadPath, *tempFile)
	metadata := structs.FileMetadata{}
	if _, errS := os.Stat(progressPath); errors.Is(errS, os.ErrNotExist) {
		errMA := os.MkdirAll(progressPath, os.ModePerm)
		if errMA != nil {
			util.Loggify(errMA)
			return false, -1, nil, fmt.Errorf("not able to locate the user's active directory in drive for file '%s', please contect the support", progressPath)
		}
		metadata.Name = name
		metadata.Size = *size
		metadata.LastModified = *lastModified
		metadata.Type = Type
		metadata.RenditionType = util.PointerString("?")
		version := structs.Version{
			Major: 0,
			Minor: 0,
			Patch: 1,
		}
		metadata.Version = &version
		versions := make([]structs.Version, 0)
		versions = append(versions, version)
		//metadata.Versions = &versions

		metadataPath := filepath.Join(progressPath, "metadata")
		metadataJson, errMJ := metadata.Json()
		if errMJ != nil || metadataJson == nil {
			util.Loggify(errMJ)
			return false, -1, nil, fmt.Errorf("not able to convert the metadata to json, try again or contact the support")
		}
		errMet := os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
		if errMet != nil {
			util.Loggify(errMet)
			return false, -1, nil, fmt.Errorf("not able to initialize the metadata, try again or contact the support")
		}
	} else {
		metadataPath := filepath.Join(progressPath, "metadata")
		data, errMet := os.ReadFile(metadataPath)
		if errMet != nil || data == nil {
			util.Loggify(errMet)
			return false, -1, nil, fmt.Errorf("not able to access the metadata information, contact the support")
		}
		errP := metadata.Parse(util.PointerString(string(data)))
		if errP != nil {
			util.Loggify(errMet)
			return false, -1, nil, fmt.Errorf("not able to parse the metadata information, contact the support")
		}
	}

	//Use base64 decoder if chunk exists, otherwise use already decoded byte array
	var decodedChunk []byte
	if chunk != nil {
		var errCh error
		decodedChunk, errCh = base64.StdEncoding.DecodeString(*chunk)
		if errCh != nil || decodedChunk == nil {
			util.Loggify(errCh)
			return false, -1, nil, fmt.Errorf("not able to decode the chang using base 64, try a base 64 decoded input instead")
		}
	} else {
		decodedChunk = *chunkByte
	}

	filePath := filepath.Join(progressPath, "file")
	file, errF := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		util.Loggify(errF)
		return false, -1, nil, fmt.Errorf("not able to append the data to a file, contact the support")
	}

	if start != nil {
		fileInfo, errStat := file.Stat()
		if errStat != nil {
			defer file.Close()
			util.Loggify(errStat)
			return false, -1, nil, fmt.Errorf("not able to gather stats for the file, contact the support")
		}
		if fileInfo.Size() >= int64(*size) {
			file.Close()
			errProg := os.RemoveAll(progressPath)
			if errProg != nil {
				util.Loggify(errProg)
				return false, -1, nil, fmt.Errorf("not able to finalize the file, contact the support")
			}
			return false, -1, nil, fmt.Errorf("try uploading file '%s' again, this time it failed", *name)
		}
		if fileInfo.Size() > *start {
			return true, fileInfo.Size(), nil, nil
		}
	}
	number, errWrF := file.Write(decodedChunk)
	if errWrF != nil || number != len(decodedChunk) {
		defer file.Close()
		util.Loggify(errF)
		return false, -1, nil, fmt.Errorf("not able to append the data to a file, contact the support")
	}

	fileInfo, errStat := file.Stat()
	if errStat != nil {
		defer file.Close()
		util.Loggify(errStat)
		return false, -1, nil, fmt.Errorf("not able to gather stats for the file, contact the support")
	}
	if metadata.Size == fileInfo.Size() {
		defer file.Close()
		uploaded := time.Now().UnixNano() / int64(time.Millisecond)
		tempFile := files.Name(util.PointerString("__file__" + base64.StdEncoding.EncodeToString([]byte(extension+"|"+strconv.FormatInt(int64(*size), 10)+"|"+strconv.FormatInt(*lastModified, 10)+"|"+*Type+"|"+strconv.FormatInt(uploaded, 10)))))
		finalPath := filepath.Join(uploadPath, *tempFile)
		errRen := os.Rename(progressPath, finalPath)
		if errRen != nil {
			util.Loggify(errRen)
			return false, fileInfo.Size(), nil, fmt.Errorf("not able to finalize the file, contact the support")
		}
		exifMetadata, errEx := structs.Exif(filepath.Join(finalPath, "file"))
		if errEx != nil {
			//log.Println(errEx.Error(), errEx)
		} else {
			metadata.ExifMetadata = &exifMetadata
		}
		thumbnail, rendition, pricessing, viewable, fileType, renditionType := rendition.Process(util.PointerString(filepath.Join(finalPath, "file")), &finalPath, false)
		metadata.Thumbnail = thumbnail
		metadata.Rendition = rendition
		metadata.Processing = pricessing
		metadata.RenditionType = renditionType
		metadata.Viewable = viewable

		if Type == nil && fileType != nil {
			metadata.Type = fileType
		}

		if metadata.Thumbnail {
			bytes, errT := os.ReadFile(filepath.Join(finalPath, "thumbnail"))
			if errT == nil {
				metadata.ThumbnailImage = util.PointerString(base64.StdEncoding.EncodeToString(bytes))
			} else {
				util.Loggify(errT)
				metadata.Thumbnail = false
			}
		}
		metadata.Id = tempFile
		fileMetadata = &metadata

		metadataPath := filepath.Join(finalPath, "metadata")
		metadataJson, errMJ := metadata.Json()
		if errMJ != nil || metadataJson == nil {
			util.Loggify(errMJ)
		}
		errMet := os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
		if errMet != nil {
			util.Loggify(errMet)
		}
	} else if metadata.Size < fileInfo.Size() {
		file.Close()
		errProg := os.RemoveAll(progressPath)
		if err != nil {
			util.Loggify(errProg)
			return false, -1, nil, fmt.Errorf("not able to finalize the file, contact the support")
		}
		return false, -1, nil, fmt.Errorf("try uploading file '%s' again, this time it failed", *name)
	}
	return true, -2, fileMetadata, nil
}

func FetchFilesOld(userId *int64, directory *string) (files *[]structs.FileMetadata, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return nil, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return nil, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return nil, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}
	if directory != nil {
		directory = util.PointerString(strings.Trim(*directory, " "))
		directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
		directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
		directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
		drivePath = filepath.Join(drivePath, *directory)
	}
	filesDirectories, errFD := os.ReadDir(drivePath)
	if errFD != nil {
		util.Loggify(errFD)
		return nil, fmt.Errorf("not able to retrieve directories and files at this moment, try again later or contact the support")
	}
	data := make([]structs.FileMetadata, 0)
	for _, file := range filesDirectories {
		if file.IsDir() && strings.HasPrefix(file.Name(), "__file__") {
			metadataPath := filepath.Join(drivePath, file.Name(), "metadata")
			metadataJson, errMet := os.ReadFile(metadataPath)
			if errMet != nil || metadataJson == nil {
				util.Loggify(errMet)
				//TODO consider notifying about broken file so users can manage it
				continue
			}
			metadata := structs.FileMetadata{}
			metadata.Id = util.PointerString(file.Name())
			errP := metadata.Parse(util.PointerString(string(metadataJson)))
			if errP != nil {
				util.Loggify(errP)
				//TODO consider notifying about broken file so users can manage it
				continue
			}
			if metadata.Thumbnail {
				bytes, errT := os.ReadFile(filepath.Join(drivePath, file.Name(), "thumbnail"))
				if errT == nil {
					metadata.ThumbnailImage = util.PointerString(base64.StdEncoding.EncodeToString(bytes))
				} else {
					util.Loggify(errT)
					metadata.Thumbnail = false
				}
			}
			if metadata.Processing {
				rendition.AssignRenditionWork(util.PointerString(filepath.Join(drivePath, file.Name())), metadata.Type, false)
			}
			data = append(data, metadata)
		} else if file.IsDir() && strings.HasPrefix(file.Name(), "__progress__") {
			//TODO consider removing progres directories if not active for certain time
		} else if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			directory := structs.FileMetadata{
				Name:  util.PointerString(file.Name()),
				IsDir: true,
			}
			data = append(data, directory)
		}
	}
	files = &data
	return files, nil
}

func NewDirectoryOld(userId *int64, directory *string, name *string) (result bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	if name == nil || strings.HasPrefix(*name, ".") || strings.HasPrefix(*name, "__progress__") || strings.HasPrefix(*name, "__file__") || strings.HasPrefix(*name, "__system__") {
		return false, fmt.Errorf("a folder name is required and can't start with '.', '__progress__', or '__file__'")
	}
	pattern := regexp.MustCompile(`^[^\s^\x00-\x1f\\?*:"";<>|\/.][^\x00-\x1f\\?*:"";<>|\/]*[^\s^\x00-\x1f\\?*:"";<>|\/.]+$`)
	found := pattern.FindString(*name)
	if found == "" {
		return false, fmt.Errorf("a folder name is required and can't start with '.', ' ', '__progress__', or '__file__' and contain characters '*, \\, :, \", /, >, <, ?, |'")
	}

	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}
	if directory != nil {
		directory = util.PointerString(strings.Trim(*directory, " "))
		directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
		directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
		directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
		drivePath = filepath.Join(drivePath, *directory)
	}
	newDirectory := filepath.Join(drivePath, *name)
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

func RemoveFileOld(userId *int64, directory *string, id *string) (result bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	filePath := drivePath
	if directory != nil {
		directory = util.PointerString(strings.Trim(*directory, " "))
		directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
		directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
		directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
		filePath = filepath.Join(drivePath, *directory)
	}
	filePath = strings.TrimRight(filePath, string(os.PathSeparator))
	filePath = filepath.Join(filePath, *id)
	errR := os.RemoveAll(filePath)
	if errR != nil {
		return false, fmt.Errorf("not able to remove the file or folder at this time, contact the support")
	}
	return true, nil
}

func RemoveFilesOld(userId *int64, root *string, ids *[]string) (result bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	filePath := drivePath
	if root != nil {
		root = util.PointerString(strings.Trim(*root, " "))
		root = util.PointerString(strings.ReplaceAll(*root, "/", string(os.PathSeparator)))
		root = util.PointerString(strings.ReplaceAll(*root, "\\", string(os.PathSeparator)))
		root = util.PointerString(strings.TrimLeft(*root, string(os.PathSeparator)))
		filePath = filepath.Join(drivePath, *root)
	}
	filePath = strings.TrimRight(filePath, string(os.PathSeparator))
	result = true
	for _, id := range *ids {
		removeFilePath := filepath.Join(filePath, id)
		errR := os.RemoveAll(removeFilePath)
		if errR != nil {
			result = false
			util.Loggify(errR)
			err = fmt.Errorf("not able to remove all the file or folders at this time, contact the support")
		}
	}

	return result, err
}

func MoveFilesOld(userId *int64, root *string, ids *[]string, directory *string) (result bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, fmt.Errorf("user is not allocated an active drive on this amper instance %d", business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	filePath := drivePath
	if root != nil {
		root = util.PointerString(strings.Trim(*root, " "))
		root = util.PointerString(strings.ReplaceAll(*root, "/", string(os.PathSeparator)))
		root = util.PointerString(strings.ReplaceAll(*root, "\\", string(os.PathSeparator)))
		root = util.PointerString(strings.TrimLeft(*root, string(os.PathSeparator)))
		filePath = filepath.Join(drivePath, *root)
	}
	filePath = strings.TrimRight(filePath, string(os.PathSeparator))
	result = true
	for _, id := range *ids {
		if *directory == id {
			continue
		}
		originalFilePath := filepath.Join(filePath, id)
		moveFilePath := filepath.Join(filePath, *directory, id)

		errRen := os.Rename(originalFilePath, moveFilePath)
		if errRen != nil {
			util.Loggify(errRen)
			result = false
			err = fmt.Errorf("not able to move all files at this time, contact the support")
		}
	}

	return result, err
}

func PasteFilesOld(userId *int64, root *string, copyJson *string, cutJson *string) (result bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	pastePath := drivePath
	if root != nil {
		root = util.PointerString(strings.Trim(*root, " "))
		root = util.PointerString(strings.ReplaceAll(*root, "/", string(os.PathSeparator)))
		root = util.PointerString(strings.ReplaceAll(*root, "\\", string(os.PathSeparator)))
		root = util.PointerString(strings.TrimLeft(*root, string(os.PathSeparator)))
		pastePath = filepath.Join(drivePath, *root)
	}
	pastePath = strings.TrimRight(pastePath, string(os.PathSeparator))
	copy, errC := jsons.GetJsonObject(copyJson)
	if errC != nil {
		util.Loggify(errC)
		return false, fmt.Errorf("the supplied copy parameter: %s is of wrong format", *copyJson)
	}
	cut, errCut := jsons.GetJsonObject(cutJson)
	if errCut != nil {
		util.Loggify(errCut)
		return false, fmt.Errorf("the supplied cut parameter: %s is of wrong format", *cutJson)
	}
	result = true
	if len(copy) < 1 && len(cut) < 1 {
		return result, err
	}
	resultError := ""
	for file, directoryInt := range copy {
		if len(file) > 0 && reflect.TypeOf(directoryInt).String() == "string" {
			directory := directoryInt.(string)
			if len(directory) > 0 {
				fullFilePath := getFullPath(drivePath, directory, file)
				if _, errS := os.Stat(fullFilePath); errors.Is(errS, os.ErrNotExist) {
					result = false
					resultError = resultError + ", " + fmt.Sprintf("Not able to find the file '%s' in the specified directory '%s'", file, directory)
				} else {
					fullTargetFilePath := filepath.Join(pastePath, file)
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
	}

	for file, directoryInt := range cut {
		if len(file) > 0 && reflect.TypeOf(directoryInt).String() == "string" {
			directory := directoryInt.(string)
			if len(directory) > 0 {
				fullFilePath := getFullPath(drivePath, directory, file)
				if _, errS := os.Stat(fullFilePath); errors.Is(errS, os.ErrNotExist) {
					result = false
					resultError = resultError + ", " + fmt.Sprintf("Not able to find the file '%s' in the specified directory '%s'", file, directory)
				} else {
					fullTargetFilePath := filepath.Join(pastePath, file)
					_, errS1 := os.Stat(fullTargetFilePath)
					index := int64(0)
					for !errors.Is(errS1, os.ErrNotExist) {
						index++
						_, errS1 = os.Stat(fullTargetFilePath + "(" + strconv.FormatInt(index, 10) + ")")
					}
					if index > 0 {
						fullTargetFilePath = fullTargetFilePath + "(" + strconv.FormatInt(index, 10) + ")"
					}
					errMove := os.Rename(fullFilePath, fullTargetFilePath)
					if errMove != nil {
						util.Loggify(errMove)
						result = false
						resultError = resultError + ", " + fmt.Sprintf("Not able to move the file '%s' to the specified directory '%s'", file, *root)
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

func CopyDirectory(scrDir, dest string) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := files.CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func getFullPath(rootPath string, relativePath string, fileName string) string {
	relativePath = strings.Trim(relativePath, " ")
	relativePath = strings.ReplaceAll(relativePath, "/", string(os.PathSeparator))
	relativePath = strings.ReplaceAll(relativePath, "\\", string(os.PathSeparator))
	relativePath = strings.TrimLeft(relativePath, string(os.PathSeparator))
	return filepath.Join(rootPath, relativePath, fileName)
}

func GetFileOld(userId *int64, directory *string, id *string, version *structs.Version, rendition bool) (result *os.File, metadata *structs.FileMetadata, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return nil, nil, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return nil, nil, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return nil, nil, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	filePath := drivePath
	if directory != nil {
		directory = util.PointerString(strings.Trim(*directory, " "))
		directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
		directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
		directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
		filePath = filepath.Join(drivePath, *directory)
	}
	filePath = strings.TrimRight(filePath, string(os.PathSeparator))

	metadata = &structs.FileMetadata{}
	metadataPath := filepath.Join(filePath, *id, "metadata")
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return nil, nil, fmt.Errorf("not able to access the metadata information, contact the support")
	}
	errP := metadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return nil, nil, fmt.Errorf("not able to parse the metadata information, contact the support")
	}
	lastVersion := metadata.Version
	metadataName := "metadata"
	if version != nil && !(version.Major == metadata.Version.Major && version.Minor == metadata.Version.Minor && version.Patch == metadata.Version.Patch) {
		metadataName = fmt.Sprintf("%s_%d_%d_%d", metadataName, version.Major, version.Minor, version.Patch)
		metadata = &structs.FileMetadata{}
		metadataPath := filepath.Join(filePath, *id, metadataName)
		data, errMet := os.ReadFile(metadataPath)
		if errMet != nil || data == nil {
			util.Loggify(errMet)
			return nil, nil, fmt.Errorf("not able to access the metadata information, contact the support")
		}
		errP := metadata.Parse(util.PointerString(string(data)))
		if errP != nil {
			util.Loggify(errMet)
			return nil, nil, fmt.Errorf("not able to parse the metadata information, contact the support")
		}
	}

	fileName := "file"
	if rendition {
		if metadata.Rendition && metadata.RenditionType != nil && *metadata.RenditionType != "?" {
			switch *metadata.RenditionType {
			case "application/pdf":
				fileName = "rendition"
			case "image/jpeg":
				fileName = "rendition"
			case "image/png":
				fileName = "rendition"
			}
		}
	}
	if version != nil && !(version.Major == lastVersion.Major && version.Minor == lastVersion.Minor && version.Patch == lastVersion.Patch) {
		fileName = fmt.Sprintf("%s_%d_%d_%d", fileName, version.Major, version.Minor, version.Patch)
	}
	filePath = filepath.Join(filePath, *id, fileName)
	file, errF := os.Open(filePath)
	if errF != nil {
		util.Loggify(errF)
		err = fmt.Errorf("the file you are try to download does not exist")
		return nil, nil, err
	}
	return file, metadata, err
}

func GetMetadataOld(userId *int64, directory *string, id *string, version *structs.Version) (metadata *structs.FileMetadata, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return nil, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return nil, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return nil, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	filePath := drivePath
	if directory != nil {
		directory = util.PointerString(strings.Trim(*directory, " "))
		directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
		directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
		directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
		filePath = filepath.Join(drivePath, *directory)
	}
	filePath = strings.TrimRight(filePath, string(os.PathSeparator))

	metadata = &structs.FileMetadata{}
	metadataPath := filepath.Join(filePath, *id, "metadata")
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return nil, fmt.Errorf("not able to access the metadata information, contact the support")
	}
	errP := metadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return nil, fmt.Errorf("not able to parse the metadata information, contact the support")
	}
	metadata.Id = id
	metadataName := "metadata"
	if version != nil && !(version.Major == metadata.Version.Major && version.Minor == metadata.Version.Minor && version.Patch == metadata.Version.Patch) {
		metadataName = fmt.Sprintf("%s_%d_%d_%d", metadataName, version.Major, version.Minor, version.Patch)
	} else {
		return metadata, nil
	}
	metadataV := &structs.FileMetadata{}
	metadataPath = filepath.Join(filePath, *id, metadataName)
	data, errMet = os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return nil, fmt.Errorf("not able to access the metadata information, contact the support")
	}
	errP = metadataV.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return nil, fmt.Errorf("not able to parse the metadata information, contact the support")
	}
	metadataV.Id = id
	//Set the versions to the latest version, to have all versions available on metadata
	metadataV.Versions = metadata.Versions
	return metadataV, nil
}

func DiscoverRootOld(userId *int64) (result *structs.Folder, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return nil, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return nil, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return nil, fmt.Errorf("not able to locate the user's active directory in drive '%s', contect the support", drivePath)
	}
	result = &structs.Folder{
		Name:    util.PointerString("Root"),
		Path:    util.PointerString(string(os.PathSeparator)),
		Folders: RootWalkOld(drivePath, string(os.PathSeparator)),
	}
	return result, err
}

func RootWalkOld(root string, relativePath string) *[]structs.Folder {
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

func UpversionFile(userId *int64, root *string, id *string, version *structs.Version, file *[]byte) (success bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	filePath := drivePath
	if root != nil {
		root = util.PointerString(strings.Trim(*root, " "))
		root = util.PointerString(strings.ReplaceAll(*root, "/", string(os.PathSeparator)))
		root = util.PointerString(strings.ReplaceAll(*root, "\\", string(os.PathSeparator)))
		root = util.PointerString(strings.TrimLeft(*root, string(os.PathSeparator)))
		filePath = filepath.Join(drivePath, *root)
	}

	metadata := structs.FileMetadata{}
	metadataPath := filepath.Join(filePath, *id, "metadata")
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return false, fmt.Errorf("not able to access the metadata information, contact the support")
	}
	errP := metadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return false, fmt.Errorf("not able to parse the metadata information, contact the support")
	}
	if metadata.Version != nil && metadata.Version.Major != version.Major && metadata.Version.Minor != version.Minor && metadata.Version.Patch != version.Patch {
		return false, fmt.Errorf("your file version doesn't match the latest file version we know, reload your file and try again")
	}
	versionString := metadata.Version.String()
	metadata.UpVersion()

	metadata.LastModified = time.Now().UnixNano() / int64(time.Millisecond)
	if len(*file) < 1 {
		return false, fmt.Errorf("we expect the new version of file have a content, try again with contentful data")
	}
	metadata.Size = int64(len(*file))
	errRen := os.Rename(filepath.Join(filePath, *id, "file"), filepath.Join(filePath, *id, "file_"+versionString))
	if errRen != nil {
		util.Loggify(errRen)
		return false, fmt.Errorf("not able to upversion the file, contact the support")
	}

	errFile := os.WriteFile(filepath.Join(filePath, *id, "file"), *file, 0644)
	if errFile != nil {
		log.Println(errFile.Error(), errFile)
		return false, fmt.Errorf("not able to store the new version of the file, contact the support")
	}
	thumbnail, rendition, pricessing, viewable, fileType, renditionType := rendition.Process(util.PointerString(filepath.Join(filePath, *id, "file")), util.PointerString(filepath.Join(filePath, *id)), true)
	metadata.Thumbnail = thumbnail
	metadata.Rendition = rendition
	metadata.Processing = pricessing
	metadata.Type = fileType
	metadata.RenditionType = renditionType
	metadata.Viewable = viewable

	metadataJson, errMJ := metadata.Json()
	if errMJ != nil || metadataJson == nil {
		util.Loggify(errMJ)
		return false, fmt.Errorf("not able to finalize the metadata information, contact the support")
	}

	errMet = os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
	if errMet != nil {
		util.Loggify(errMet)
		return false, fmt.Errorf("not able to finalize the metadata information, contact the support")
	}
	if metadata.Thumbnail {
		bytes, errT := os.ReadFile(filepath.Join(filePath, *id, "thumbnail"))
		if errT == nil {
			metadata.ThumbnailImage = util.PointerString(base64.StdEncoding.EncodeToString(bytes))
		} else {
			util.Loggify(errT)
			metadata.Thumbnail = false
		}
	}
	return true, nil
}

func UpversionFileChunk(userId *int64, directory *string, id *string, version structs.Version, name *string, lastModified *int64, size *int64, dataType *string, start *int64, chunk *string, chunkByte *[]byte) (success bool, startRes int64, fileMetadata *structs.FileMetadata, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), false, false)
	if errU != nil || user == nil || *user.AmperId == 0 && *user.AmperId != *business.AmperId() {
		util.Loggify(errU)
		return false, -1, nil, fmt.Errorf("user is not allocated an active drive on this amper instance %d", *business.AmperId())
	}
	settings, errS := FetchSettings(userId)
	if errS != nil || !ampstrings.HasValue(settings.RootDirectory) {
		util.Loggify(errS)
		return false, -1, nil, fmt.Errorf("the root directory is not configured for the instance %d, check the Administration > Settings to configure it", *business.AmperId())
	}
	rootDriveDirectory := *settings.RootDirectory
	drivePath := filepath.Join(rootDriveDirectory, strconv.FormatInt(*userId, 10), "drive")
	errUD := os.MkdirAll(drivePath, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return false, -1, nil, fmt.Errorf("not able to locate the user's active directory in drive '%s', please contect the support", drivePath)
	}

	filePath := drivePath
	if directory != nil {
		directory = util.PointerString(strings.Trim(*directory, " "))
		directory = util.PointerString(strings.ReplaceAll(*directory, "/", string(os.PathSeparator)))
		directory = util.PointerString(strings.ReplaceAll(*directory, "\\", string(os.PathSeparator)))
		directory = util.PointerString(strings.TrimLeft(*directory, string(os.PathSeparator)))
		filePath = filepath.Join(drivePath, *directory)
	}
	if lastModified == nil {
		currentTime := time.Now().UnixNano() / int64(time.Millisecond)
		lastModified = &currentTime
	}

	if dataType == nil || len(*dataType) < 1 {
		dataType = util.PointerString("?")
	}

	metadata := structs.FileMetadata{}
	metadataPath := filepath.Join(filePath, *id, "metadata")
	data, errMet := os.ReadFile(metadataPath)
	if errMet != nil || data == nil {
		util.Loggify(errMet)
		return false, *start, nil, fmt.Errorf("not able to access the metadata information, contact the support")
	}
	errP := metadata.Parse(util.PointerString(string(data)))
	if errP != nil {
		util.Loggify(errMet)
		return false, *start, nil, fmt.Errorf("not able to parse the metadata information, contact the support")
	}
	if metadata.Version != nil && metadata.Version.Major != version.Major && metadata.Version.Minor != version.Minor && metadata.Version.Patch != version.Patch {
		return false, *start, nil, fmt.Errorf("your file version doesn't match the latest file version we know, reload your file and try again")
	}

	//Use base64 decoder if chunk exists, otherwise use already decoded byte array
	var decodedChunk []byte
	if chunk != nil {
		var errCh error
		decodedChunk, errCh = base64.StdEncoding.DecodeString(*chunk)
		if errCh != nil || decodedChunk == nil {
			util.Loggify(errCh)
			return false, -1, nil, fmt.Errorf("not able to decode the chang using base 64, try a base 64 decoded input instead")
		}
	} else {
		decodedChunk = *chunkByte
	}

	upFilePath := filepath.Join(filePath, *id, "up_progress_file")
	upProgressFile, errF := os.OpenFile(upFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errF != nil {
		util.Loggify(errF)
		return false, -1, nil, fmt.Errorf("not able to append the data to a file, contact the support")
	}

	if start != nil && *start > 0 {
		fileInfo, errStat := upProgressFile.Stat()
		if errStat != nil {
			os.RemoveAll(upFilePath)
			defer upProgressFile.Close()
			util.Loggify(errStat)
			return false, -1, nil, fmt.Errorf("not able to gather stats for the file, contact the support")
		}
		if fileInfo.Size() >= *size {
			upProgressFile.Close()
			errProg := os.RemoveAll(upFilePath)
			if errProg != nil {
				util.Loggify(errProg)
				return false, -1, nil, fmt.Errorf("not able to finalize the file, contact the support")
			}
			return false, -1, nil, fmt.Errorf("try uploading file '%s' again, this time it failed", *name)
		}
		if fileInfo.Size() > *start {
			return true, fileInfo.Size(), nil, nil
		}
	}

	number, errWrF := upProgressFile.Write(decodedChunk)
	if errWrF != nil || number != len(decodedChunk) {
		os.RemoveAll(upFilePath)
		defer upProgressFile.Close()
		util.Loggify(errF)
		return false, -1, nil, fmt.Errorf("not able to append the data to a file, contact the support")
	}

	fileInfo, errStat := upProgressFile.Stat()
	if errStat != nil {
		defer upProgressFile.Close()
		util.Loggify(errStat)
		return false, -1, nil, fmt.Errorf("not able to gather stats for the file, contact the support")
	}
	if *size == fileInfo.Size() {
		//append the version to the old file
		originalFilePath := filepath.Join(filePath, *id, "file")
		oldFilePath := filepath.Join(filePath, *id, fmt.Sprintf("file_%d_%d_%d", version.Major, version.Minor, version.Patch))
		errRen := os.Rename(originalFilePath, oldFilePath)
		if errRen != nil {
			util.Loggify(errRen)
			return false, fileInfo.Size(), nil, fmt.Errorf("not able to finalize the file, contact the support")
		}
		//append the version to the old metadata
		originalMetadaaPath := filepath.Join(filePath, *id, "metadata")
		oldMetadataPath := filepath.Join(filePath, *id, fmt.Sprintf("metadata_%d_%d_%d", version.Major, version.Minor, version.Patch))
		errReM := os.Rename(originalMetadaaPath, oldMetadataPath)
		if errReM != nil {
			util.Loggify(errReM)
			return false, fileInfo.Size(), nil, fmt.Errorf("not able to finalize the file, contact the support")
		}
		//append the version to the old metadata
		newFilePath := filepath.Join(filePath, *id, "file")
		errReNF := os.Rename(upFilePath, newFilePath)
		if errReNF != nil {
			util.Loggify(errReNF)
			return false, fileInfo.Size(), nil, fmt.Errorf("not able to finalize the file, contact the support")
		}
		//version the old rendition
		renditionFilePath := filepath.Join(filePath, *id, "rendition")
		renditionFilePathVersion := filepath.Join(filePath, *id, fmt.Sprintf("rendition_%d_%d_%d", version.Major, version.Minor, version.Patch))
		if fileExist(renditionFilePath) {
			errReNFR := os.Rename(renditionFilePath, renditionFilePathVersion)
			if errReNFR != nil {
				util.Loggify(errReNFR)
			}
		}
		//version the old thumbnail
		thumbnailFilePath := filepath.Join(filePath, *id, "thumbnail")
		thumbnailFilePathVersion := filepath.Join(filePath, *id, fmt.Sprintf("thumbnail_%d_%d_%d", version.Major, version.Minor, version.Patch))
		if fileExist(thumbnailFilePath) {
			errReNFT := os.Rename(thumbnailFilePath, thumbnailFilePathVersion)
			if errReNFT != nil {
				util.Loggify(errReNFT)
			}
		}
		metadata.Name = name
		exifMetadata, errEx := structs.Exif(newFilePath)
		if errEx != nil {
			//log.Println(errEx.Error(), errEx)
		} else {
			metadata.ExifMetadata = &exifMetadata
		}
		thumbnail, rendition, pricessing, viewable, fileType, renditionType := rendition.Process(util.PointerString(newFilePath), util.PointerString(filepath.Join(filePath, *id)), false)
		metadata.Thumbnail = thumbnail
		metadata.Rendition = rendition
		metadata.Processing = pricessing
		metadata.RenditionType = renditionType
		metadata.Viewable = viewable

		if fileType != nil {
			metadata.Type = fileType
		}
		metadata.Size = fileInfo.Size()
		metadata.UpVersion()

		metadata.LastModified = util.IfElse(lastModified == nil, time.Now().UnixNano()/int64(time.Millisecond), *lastModified).(int64)

		metadataJson, errMJ := metadata.Json()
		if errMJ != nil || metadataJson == nil {
			util.Loggify(errMJ)
		}

		metadataPath := filepath.Join(filePath, *id, "metadata")
		errMet := os.WriteFile(metadataPath, []byte(*metadataJson), 0644)
		if errMet != nil {
			util.Loggify(errMet)
		}
		if metadata.Thumbnail {
			bytes, errT := os.ReadFile(filepath.Join(filePath, "thumbnail"))
			if errT == nil {
				metadata.ThumbnailImage = util.PointerString(base64.StdEncoding.EncodeToString(bytes))
			} else {
				util.Loggify(errT)
				metadata.Thumbnail = false
			}
		}
		metadata.Id = id
		fileMetadata = &metadata
	} else if *size < fileInfo.Size() {
		upProgressFile.Close()
		errProg := os.RemoveAll(upFilePath)
		if err != nil {
			util.Loggify(errProg)
			return false, -1, nil, fmt.Errorf("not able to finalize the file, contact the support")
		}
		return false, -1, nil, fmt.Errorf("try uploading file '%s' again, this time it failed", *name)
	}
	return true, -2, fileMetadata, nil
}

// function to check if file exists
func fileExist(fileName string) bool {
	_, error := os.Stat(fileName)
	return !os.IsNotExist(error)
}
