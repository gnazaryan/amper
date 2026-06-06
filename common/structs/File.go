package structs

import (
	"amper/common/util"
	"amper/common/util/ampstrings"
	"amper/common/util/jsons"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

type ThumbnailsResult struct {
	Result
	Data map[string]*string
}

type Rendition struct {
	Thumbnail     bool    `json:"thumbnail"`
	Rendition     bool    `json:"rendition"`
	Processing    bool    `json:"processing"`
	Viewable      bool    `json:"viewable"`
	FileType      *string `json:"fileType"`
	RenditionType *string `json:"renditionType"`
}

type KeyValueStoreResult struct {
	Result
	Rendition *Rendition `json:"rendition"`
}

type ChunkResult struct {
	Result
	Start    int64         `json:"start"`
	Metadata *FileMetadata `json:"metadata"`
}

type Version struct {
	Major int64 `json:"major"`
	Minor int64 `json:"minor"`
	Patch int64 `json:"patch"`
}

func (v *Version) CompareTo(v1 *Version) int {
	if v == nil && v1 == nil {
		return 0
	} else if v == nil && v1 != nil {
		return -1
	} else if v != nil && v1 == nil {
		return 1
	}
	if v.Major > v1.Major {
		return 1
	} else if v.Major == v1.Major {
		if v.Minor > v1.Minor {
			return 1
		} else if v.Minor == v1.Minor {
			if v.Patch > v1.Patch {
				return 1
			} else if v.Patch == v1.Patch {
				return 0
			} else if v.Patch < v1.Patch {
				return -1
			}
		} else if v.Minor < v1.Minor {
			return -1
		}
	} else if v.Major < v1.Major {
		return -1
	}
	return 0
}

type FileMetadata struct {
	Id                *string          `json:"id"`
	Name              *string          `json:"name"`
	IsDir             bool             `json:"isDir"`
	Size              int64            `json:"size"`
	LastModified      int64            `json:"lastModified"`
	Type              *string          `json:"type"`
	Version           *Version         `json:"version"`
	Versions          *[]*FileMetadata `json:"versions"`
	AvailableVersions *[]Version       `json:"availableVersions"`
	Thumbnail         bool             `json:"thumbnail"`
	Rendition         bool             `json:"rendition"`
	Viewable          bool             `json:"viewable"`
	Processing        bool             `json:"processing"`
	ThumbnailImage    *string          `json:"thumbnailImage"`
	RenditionType     *string          `json:"renditionType"`
	ExifMetadata      *[]ExifMetadata  `json:"exifMetadata"`
	Metadata          *string
}

type MetadataResult struct {
	Result
	Data *FileMetadata `json:"data"`
}

type Files struct {
	Result
	Data *[]FileMetadata `json:"data"`
}

type Folder struct {
	Name    *string   `json:"name"`
	Path    *string   `json:"path"`
	Folders *[]Folder `json:"folders"`
}

type FolderResult struct {
	Result
	Data *Folder `json:"data"`
}

type FileId struct {
	SourceInstanceId int64   `json:"sourceInstanceId"`
	InstanceId       int64   `json:"instanceId"`
	UserId           int64   `json:"userid"`
	Year             int     `json:"year"`
	Month            int     `json:"month"`
	Day              int     `json:"day"`
	UUID             *string `json:"uuid"`
	Parent           *FileId `json:"parent"`
	OptionalValue    *string `json:"optionalValue"`
}

func (id *FileId) Format() *string {
	result := fmt.Sprintf("%d_%d_%d_%d_%d_%d_%s_%s", id.SourceInstanceId, id.InstanceId, id.UserId, id.Year, id.Month, id.Day, *id.UUID, *id.OptionalValue)
	if id.Parent != nil {
		parentId := fmt.Sprintf("%d_%d_%d_%d_%d_%d_%s_%s", id.Parent.SourceInstanceId, id.Parent.InstanceId, id.Parent.UserId, id.Parent.Year, id.Parent.Month, id.Parent.Day, *id.Parent.UUID, *id.Parent.OptionalValue)
		result = fmt.Sprintf("%s%s%s", result, ampstrings.SEPERATOR, parentId)
	}
	result = base64.StdEncoding.EncodeToString([]byte(result))
	return &result

}

func parseIdInternal(fileId *string) (result *FileId, err error) {
	result = &FileId{}
	if fileId != nil {
		idSplit := strings.Split(*fileId, "_")
		if len(idSplit) != 8 {
			return nil, fmt.Errorf("not aple to parse the file id due to inacurate format: %s", *fileId)
		}
		sourceInstanceId, errSII := strconv.ParseInt(idSplit[0], 10, 64)
		if errSII != nil {
			util.Loggify(errSII)
			return nil, fmt.Errorf("not aple to parse the file id due to inacurate instanceId value: %s", idSplit[0])
		}
		result.SourceInstanceId = sourceInstanceId
		instanceId, errII := strconv.ParseInt(idSplit[1], 10, 64)
		if errII != nil {
			util.Loggify(errII)
			return nil, fmt.Errorf("not aple to parse the file id due to inacurate instanceId value: %s", idSplit[1])
		}
		result.InstanceId = instanceId

		userId, errUI := strconv.ParseInt(idSplit[2], 10, 64)
		if errUI != nil {
			util.Loggify(errUI)
			return nil, fmt.Errorf("not aple to parse the file id due to inacurate userId value: %s", idSplit[2])
		}
		result.UserId = userId

		year, errY := strconv.Atoi(idSplit[3])
		if errY != nil {
			util.Loggify(errY)
			return nil, fmt.Errorf("not aple to parse the file id due to inacurate year value: %s", idSplit[3])
		}
		result.Year = year

		month, errM := strconv.Atoi(idSplit[4])
		if errM != nil {
			util.Loggify(errM)
			return nil, fmt.Errorf("not aple to parse the file id due to inacurate month value: %s", idSplit[4])
		}
		result.Month = month

		day, errD := strconv.Atoi(idSplit[5])
		if errD != nil {
			util.Loggify(errD)
			return nil, fmt.Errorf("not aple to parse the file id due to inacurate day value: %s", idSplit[5])
		}
		result.Day = day

		result.UUID = util.PointerString(idSplit[6])
		result.OptionalValue = util.PointerString(idSplit[7])
	}

	return result, err
}

func ParseId(id *string) (result *FileId, err error) {
	if id != nil {
		result = &FileId{}
		idDecodedByte, errId := base64.StdEncoding.DecodeString(*id)
		if errId != nil || idDecodedByte == nil {
			util.Loggify(errId)
			return nil, fmt.Errorf("not able to base64 decode the id: %s", *id)
		}
		idDecoded := string(idDecodedByte)
		idParts := strings.Split(idDecoded, ampstrings.SEPERATOR)
		if len(idParts) > 0 {
			var errFI error
			result, errFI = parseIdInternal(util.PointerString(idParts[0]))
			if errFI != nil {
				util.Loggify(errFI)
				return nil, fmt.Errorf("not able to parse the self part of the file id: %s", idParts[0])
			}
		}
		if len(idParts) > 1 {
			var errFI error
			parentId, errFI := parseIdInternal(util.PointerString(idParts[1]))
			if errFI != nil {
				util.Loggify(errFI)
				return nil, fmt.Errorf("not able to parse the parent part of the file id: %s", idParts[1])
			}
			result.Parent = parentId
		}
	} else {
		err = fmt.Errorf("not aple to parse the file id since id is null or empty")
	}
	return result, err
}

func (f *FileMetadata) FileNameWithoutExt() string {
	if f.Name == nil {
		return ""
	}
	return (*f.Name)[:len((*f.Name))-len(filepath.Ext((*f.Name)))]
}

func (fm *FileMetadata) UpVersion() {
	fm.Version.UpVersion()
}

func (v *Version) UpVersion() {
	if v.Patch < 99 {
		v.Patch++
		return
	}
	if v.Minor < 99 {
		v.Minor++
		v.Patch = 0
		return
	}
	v.Major++
	v.Minor = 0
	v.Patch = 0
}

func (v *Version) String() string {
	return fmt.Sprintf("%s_%s_%s", strconv.Itoa(int(v.Major)), strconv.Itoa(int(v.Minor)), strconv.Itoa(int(v.Patch)))
}

func (f *FileMetadata) parseInternal(metadata map[string]interface{}) error {
	if reflect.TypeOf(metadata["id"]).String() != "string" {
		return fmt.Errorf("not able to process the metadata due to parsing error in id value")
	}
	f.Id = util.PointerString(metadata["id"].(string))
	if reflect.TypeOf(metadata["name"]).String() != "string" {
		return fmt.Errorf("not able to process the metadata due to parsing error in name value")
	}
	f.Name = util.PointerString(metadata["name"].(string))
	size, errS := util.I2Num(metadata["size"])
	if errS != nil {
		return fmt.Errorf("not abple to process the metadata due to parsing error in size value")
	}
	f.Size = size
	lastModified, errL := util.I2Num(metadata["lastModified"])
	if errL != nil {
		return fmt.Errorf("not able to process the metadata due to parsing error in last modified value")
	}
	f.LastModified = lastModified

	if metadata["type"] != nil && reflect.TypeOf(metadata["type"]).String() != "string" {
		return fmt.Errorf("not able to process the metadata due to parsing error in type value")
	}
	f.Type = util.PointerString(metadata["type"].(string))

	if metadata["thumbnail"] != nil && reflect.TypeOf(metadata["thumbnail"]).String() != "bool" {
		return fmt.Errorf("not able to process the metadata due to parsing error in thumbnail value")
	} else if metadata["thumbnail"] != nil {
		f.Thumbnail = metadata["thumbnail"].(bool)
	} else {
		f.Thumbnail = false
	}

	if metadata["rendition"] != nil && reflect.TypeOf(metadata["rendition"]).String() != "bool" {
		return fmt.Errorf("not able to process the metadata due to parsing error in rendition value")
	} else if metadata["rendition"] != nil {
		f.Rendition = metadata["rendition"].(bool)
	} else {
		f.Rendition = false
	}

	if metadata["viewable"] != nil && reflect.TypeOf(metadata["viewable"]).String() != "bool" {
		return fmt.Errorf("not able to process the metadata due to parsing error in viewable value")
	} else if metadata["viewable"] != nil {
		f.Viewable = metadata["viewable"].(bool)
	} else {
		f.Viewable = false
	}

	if metadata["processing"] != nil && reflect.TypeOf(metadata["processing"]).String() != "bool" {
		return fmt.Errorf("not able to process the metadata due to parsing error in processing value")
	} else if metadata["processing"] != nil {
		f.Processing = metadata["processing"].(bool)
	} else {
		f.Processing = false
	}

	if metadata["renditionType"] != nil && reflect.TypeOf(metadata["renditionType"]).String() != "string" {
		return fmt.Errorf("not able to process the metadata due to parsing error in RenditionType value")
	} else if metadata["renditionType"] == nil {
		f.RenditionType = util.PointerString("?")
	} else {
		f.RenditionType = util.PointerString(metadata["renditionType"].(string))
	}

	if metadata["version"] != nil && reflect.TypeOf(metadata["version"]).String() != "map[string]interface {}" {
		return fmt.Errorf("not able to process the metadata due to parsing error in version value")
	}
	versionMap := metadata["version"].(map[string]interface{})
	if versionMap == nil {
		versionMap = make(map[string]interface{})
	}
	version := Version{
		Major: 0,
		Minor: 0,
		Patch: 1,
	}
	if versionMap["major"] != nil {
		version.Major, _ = util.I2Num(versionMap["major"])
	}
	if versionMap["minor"] != nil {
		version.Minor, _ = util.I2Num(versionMap["minor"])
	}
	if versionMap["patch"] != nil {
		version.Patch, _ = util.I2Num(versionMap["patch"])
	}
	f.Version = &version
	if metadata["versions"] != nil && reflect.TypeOf(metadata["versions"]).String() != "[]interface {}" {
		return fmt.Errorf("not abple to process the metadata due to parsing error in versions value")
	}
	versions := make([]*FileMetadata, 0)
	if metadata["versions"] != nil {
		versionsMap := metadata["versions"].([]interface{})
		for i := 0; i < len(versionsMap); i++ {
			versionMapInt := versionsMap[i]
			if reflect.TypeOf(versionMapInt).String() == "map[string]interface {}" {
				versionMap := versionsMap[i].(map[string]interface{})
				metadataVersion := FileMetadata{}
				errV := metadataVersion.parseInternal(versionMap)
				if errV != nil {
					util.Loggify(errV)
					continue
				}

				versions = append(versions, &metadataVersion)
			}
		}
	}

	f.Versions = &versions

	exifMetadatas := make([]ExifMetadata, 0)
	if metadata["exifMetadata"] != nil {
		exifMetadatasJson := metadata["exifMetadata"].([]interface{})
		if exifMetadatasJson != nil {
			for i := 0; i < len(exifMetadatasJson); i++ {
				exifMetadataJsonInit := exifMetadatasJson[i]
				if reflect.TypeOf(exifMetadataJsonInit).String() == "map[string]interface {}" {
					exifMetadataJson := exifMetadatasJson[i].(map[string]interface{})
					id, _ := util.I2Num(exifMetadataJson["id"])
					count, _ := util.I2Num(exifMetadataJson["count"])
					exifMetadatas = append(exifMetadatas, ExifMetadata{
						IfdPath: exifMetadataJson["ifdPath"].(string),
						Id:      uint16(id),
						Name:    exifMetadataJson["name"].(string),
						Count:   uint32(count),
						Type:    exifMetadataJson["type"].(string),
						Value:   exifMetadataJson["value"].(string),
					})
				}
			}
		}
	}
	f.ExifMetadata = &exifMetadatas
	return nil
}

func (f *FileMetadata) Parse(metadata *string) error {
	if metadata != nil {
		metadataJson, errJ := jsons.GetJsonObject(metadata)
		if errJ != nil {
			util.Loggify(errJ)
			return fmt.Errorf("not able to process the metadata due to parsing error")
		}
		return f.parseInternal(metadataJson)
	}
	return nil
}

func (f *FileMetadata) Json() (*string, error) {
	if f != nil {
		marshaled, errM := json.Marshal(*f)
		if errM != nil {
			util.Loggify(errM)
			return nil, fmt.Errorf("metadata is invalid and can not be converted to json")
		}
		return util.PointerString(string(marshaled)), nil
	}
	return nil, fmt.Errorf("the metadata is nil and caan't be converted to json")
}

func Exif(filepath string) (result []ExifMetadata, err error) {
	rawExif, errEI := exif.SearchFileAndExtractExif(filepath)
	if errEI != nil {
		return nil, fmt.Errorf("not able to exif the file with a file path %s", filepath)
	}
	im, errNIMS := exifcommon.NewIfdMappingWithStandard()
	if errNIMS != nil {
		return nil, fmt.Errorf("not able to initiate a Ifd mapping standard for file %s", filepath)
	}
	ti := exif.NewTagIndex()
	visitor := func(ite *exif.IfdTagEntry) (err error) {
		tagId := ite.TagId()
		tagType := ite.TagType()
		ii := ite.IfdIdentity()

		it, errG := ti.Get(ii, tagId)
		if errG != nil {
			log.Println(errG.Error(), errG)
		}

		valueString, errFF := ite.Format()
		if errFF != nil {
			log.Println(errFF.Error(), errFF)
		}
		result = append(result, ExifMetadata{
			IfdPath: ii.String(),
			Id:      tagId,
			Name:    it.Name,
			Count:   ite.UnitCount(),
			Type:    tagType.String(),
			Value:   valueString,
		})

		return nil
	}
	_, _, errV := exif.Visit(exifcommon.IfdStandardIfdIdentity, im, ti, rawExif, visitor, nil)
	if errV != nil {
		log.Println(errV.Error(), errV)
	}
	return result, err
}

type ExifMetadata struct {
	IfdPath string `json:"ifdPath"`
	Id      uint16 `json:"id"`
	Name    string `json:"name"`
	Count   uint32 `json:"count"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

type sniffSig interface {
	// match returns the MIME type of the data, or "" if unknown.
	match(data []byte, firstNonWS int) string
}
type htmlSig []byte
type maskedSig struct {
	mask, pat []byte
	skipWS    bool
	ct        string
}

type exactSig struct {
	sig []byte
	ct  string
}

type mp4Sig struct{}

var mp4ftype = []byte("ftyp")
var mp4 = []byte("mp4")

type textSig struct{}

func (textSig) match(data []byte, firstNonWS int) string {
	// c.f. section 5, step 4.
	for _, b := range data[firstNonWS:] {
		switch {
		case b <= 0x08,
			b == 0x0B,
			0x0E <= b && b <= 0x1A,
			0x1C <= b && b <= 0x1F:
			return ""
		}
	}
	return "text/plain; charset=utf-8"
}

func (mp4Sig) match(data []byte, firstNonWS int) string {
	// https://mimesniff.spec.whatwg.org/#signature-for-mp4
	// c.f. section 6.2.1
	if len(data) < 12 {
		return ""
	}
	boxSize := int(binary.BigEndian.Uint32(data[:4]))
	if len(data) < boxSize || boxSize%4 != 0 {
		return ""
	}
	if !bytes.Equal(data[4:8], mp4ftype) {
		return ""
	}
	for st := 8; st < boxSize; st += 4 {
		if st == 12 {
			// Ignores the four bytes that correspond to the version number of the "major brand".
			continue
		}
		if bytes.Equal(data[st:st+3], mp4) {
			return "video/mp4"
		}
	}
	return ""
}

func (e *exactSig) match(data []byte, firstNonWS int) string {
	if bytes.HasPrefix(data, e.sig) {
		return e.ct
	}
	return ""
}

func (m *maskedSig) match(data []byte, firstNonWS int) string {
	// pattern matching algorithm section 6
	// https://mimesniff.spec.whatwg.org/#pattern-matching-algorithm

	if m.skipWS {
		data = data[firstNonWS:]
	}
	if len(m.pat) != len(m.mask) {
		return ""
	}
	if len(data) < len(m.pat) {
		return ""
	}
	for i, pb := range m.pat {
		maskedData := data[i] & m.mask[i]
		if maskedData != pb {
			return ""
		}
	}
	return m.ct
}

// isTT reports whether the provided byte is a tag-terminating byte (0xTT)
// as defined in https://mimesniff.spec.whatwg.org/#terminology.
func isTT(b byte) bool {
	switch b {
	case ' ', '>':
		return true
	}
	return false
}

func (h htmlSig) match(data []byte, firstNonWS int) string {
	data = data[firstNonWS:]
	if len(data) < len(h)+1 {
		return ""
	}
	for i, b := range h {
		db := data[i]
		if 'A' <= b && b <= 'Z' {
			db &= 0xDF
		}
		if b != db {
			return ""
		}
	}
	// Next byte must be a tag-terminating byte(0xTT).
	if !isTT(data[len(h)]) {
		return ""
	}
	return "text/html; charset=utf-8"
}

// Data matching the table in section 6.
var sniffSignatures = []sniffSig{
	htmlSig("<!DOCTYPE HTML"),
	htmlSig("<HTML"),
	htmlSig("<HEAD"),
	htmlSig("<SCRIPT"),
	htmlSig("<IFRAME"),
	htmlSig("<H1"),
	htmlSig("<DIV"),
	htmlSig("<FONT"),
	htmlSig("<TABLE"),
	htmlSig("<A"),
	htmlSig("<STYLE"),
	htmlSig("<TITLE"),
	htmlSig("<B"),
	htmlSig("<BODY"),
	htmlSig("<BR"),
	htmlSig("<P"),
	htmlSig("<!--"),
	&maskedSig{
		mask:   []byte("\xFF\xFF\xFF\xFF\xFF"),
		pat:    []byte("<?xml"),
		skipWS: true,
		ct:     "text/xml; charset=utf-8"},
	&exactSig{[]byte("%PDF-"), "application/pdf"},
	&exactSig{[]byte("%!PS-Adobe-"), "application/postscript"},

	// UTF BOMs.
	&maskedSig{
		mask: []byte("\xFF\xFF\x00\x00"),
		pat:  []byte("\xFE\xFF\x00\x00"),
		ct:   "text/plain; charset=utf-16be",
	},
	&maskedSig{
		mask: []byte("\xFF\xFF\x00\x00"),
		pat:  []byte("\xFF\xFE\x00\x00"),
		ct:   "text/plain; charset=utf-16le",
	},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\x00"),
		pat:  []byte("\xEF\xBB\xBF\x00"),
		ct:   "text/plain; charset=utf-8",
	},

	// Image types
	// For posterity, we originally returned "image/vnd.microsoft.icon" from
	// https://tools.ietf.org/html/draft-ietf-websec-mime-sniff-03#section-7
	// https://codereview.appspot.com/4746042
	// but that has since been replaced with "image/x-icon" in Section 6.2
	// of https://mimesniff.spec.whatwg.org/#matching-an-image-type-pattern
	&exactSig{[]byte("\x00\x00\x01\x00"), "image/x-icon"},
	&exactSig{[]byte("\x00\x00\x02\x00"), "image/x-icon"},
	&exactSig{[]byte("BM"), "image/bmp"},
	&exactSig{[]byte("GIF87a"), "image/gif"},
	&exactSig{[]byte("GIF89a"), "image/gif"},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF\xFF\xFF"),
		pat:  []byte("RIFF\x00\x00\x00\x00WEBPVP"),
		ct:   "image/webp",
	},
	&exactSig{[]byte("\x89PNG\x0D\x0A\x1A\x0A"), "image/png"},
	&exactSig{[]byte("\xFF\xD8\xFF"), "image/jpeg"},
	&exactSig{[]byte("\x00\x00\x00$ftypheic\x00\x00\x00\x00mif1MiPrmiafMiHBheic"), "image/heic"},
	&exactSig{[]byte("II*"), "image/tiff"},
	&exactSig{[]byte("MM*"), "image/tiff"},
	&exactSig{[]byte("PK\x03\x04"), "application/docx/xlsx/pptx"},
	&exactSig{[]byte("\xd0\xcf\x11\u0871"), "application/doc/xls/ppt"},

	// Audio and Video types
	// Enforce the pattern match ordering as prescribed in
	// https://mimesniff.spec.whatwg.org/#matching-an-audio-or-video-type-pattern
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),
		pat:  []byte("FORM\x00\x00\x00\x00AIFF"),
		ct:   "audio/aiff",
	},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF"),
		pat:  []byte("ID3"),
		ct:   "audio/mpeg",
	},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\xFF\xFF"),
		pat:  []byte("OggS\x00"),
		ct:   "application/ogg",
	},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"),
		pat:  []byte("MThd\x00\x00\x00\x06"),
		ct:   "audio/midi",
	},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),
		pat:  []byte("RIFF\x00\x00\x00\x00AVI "),
		ct:   "video/avi",
	},
	&maskedSig{
		mask: []byte("\xFF\xFF\xFF\xFF\x00\x00\x00\x00\xFF\xFF\xFF\xFF"),
		pat:  []byte("RIFF\x00\x00\x00\x00WAVE"),
		ct:   "audio/wave",
	},
	// 6.2.0.2. video/mp4
	mp4Sig{},
	// 6.2.0.3. video/webm
	&exactSig{[]byte("\x1A\x45\xDF\xA3"), "video/webm"},

	// Font types
	&maskedSig{
		// 34 NULL bytes followed by the string "LP"
		pat: []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00LP"),
		// 34 NULL bytes followed by \xF\xF
		mask: []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF\xFF"),
		ct:   "application/vnd.ms-fontobject",
	},
	&exactSig{[]byte("\x00\x01\x00\x00"), "font/ttf"},
	&exactSig{[]byte("OTTO"), "font/otf"},
	&exactSig{[]byte("ttcf"), "font/collection"},
	&exactSig{[]byte("wOFF"), "font/woff"},
	&exactSig{[]byte("wOF2"), "font/woff2"},

	// Archive types
	&exactSig{[]byte("\x1F\x8B\x08"), "application/x-gzip"},
	&exactSig{[]byte("PK\x03\x04"), "application/zip"},
	// RAR's signatures are incorrectly defined by the MIME spec as per
	//    https://github.com/whatwg/mimesniff/issues/63
	// However, RAR Labs correctly defines it at:
	//    https://www.rarlab.com/technote.htm#rarsign
	// so we use the definition from RAR Labs.
	// TODO: do whatever the spec ends up doing.
	&exactSig{[]byte("Rar!\x1A\x07\x00"), "application/x-rar-compressed"},     // RAR v1.5-v4.0
	&exactSig{[]byte("Rar!\x1A\x07\x01\x00"), "application/x-rar-compressed"}, // RAR v5+

	&exactSig{[]byte("\x00\x61\x73\x6D"), "application/wasm"},

	textSig{}, // should be last
}

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 512

// isWS reports whether the provided byte is a whitespace byte (0xWS)
// as defined in https://mimesniff.spec.whatwg.org/#terminology.
func isWS(b byte) bool {
	switch b {
	case '\t', '\n', '\x0c', '\r', ' ':
		return true
	}
	return false
}

// DetectContentType implements the algorithm described
// at https://mimesniff.spec.whatwg.org/ to determine the
// Content-Type of the given data. It considers at most the
// first 512 bytes of data. DetectContentType always returns
// a valid MIME type: if it cannot determine a more specific one, it
// returns "application/octet-stream".
func DetectContentType(data []byte) string {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	// Index of the first non-whitespace byte in data.
	firstNonWS := 0
	for ; firstNonWS < len(data) && isWS(data[firstNonWS]); firstNonWS++ {
	}

	for _, sig := range sniffSignatures {
		if ct := sig.match(data, firstNonWS); ct != "" {
			return ct
		}
	}

	return "application/octet-stream" // fallback
}
