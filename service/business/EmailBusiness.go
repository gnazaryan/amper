package business

import (
	"amper/api/email/imap"
	"amper/api/email/imapclient"
	"amper/api/email/imapserver/imapmemserver"
	"amper/cache/business"
	"amper/common/crypto"
	"amper/common/email"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/ampstrings"
	"amper/common/util/datetime"
	"amper/common/util/files"
	"amper/data/database"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var EMAIL_PATH = filepath.Join("__system__", "Email")
var ALL_EMAIL = "ALL_EMAIL"
var EMAIL_PARTITION = 100

var emailLock EmailLock

func init() {
	emailLock = EmailLock{emails: make(map[string]bool)}
}

type EmailLock struct {
	mu     sync.Mutex
	emails map[string]bool
}

// Lock the email for checking the emails, until the check is complete
func (c *EmailLock) Lock(email string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.emails[email] {
		return false
	}
	// Lock so only one goroutine at a time can access the map email locks
	c.emails[email] = true
	return true
}

// Check to see if the given email is locked
func (c *EmailLock) Locked(key string) bool {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map
	defer c.mu.Unlock()
	return c.emails[key]
}

// Lock the email for checking the emails, until the check is complete
func (c *EmailLock) Unlock(email string) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map email locks
	c.emails[email] = false
	c.mu.Unlock()
}

var EMAIL_FLAGS = map[string]imap.Flag{
	string(imap.FlagSeen):      imap.FlagSeen,
	string(imap.FlagAnswered):  imap.FlagAnswered,
	string(imap.FlagFlagged):   imap.FlagFlagged,
	string(imap.FlagDeleted):   imap.FlagDeleted,
	string(imap.FlagDraft):     imap.FlagDraft,
	string(imap.FlagForwarded): imap.FlagForwarded,
	string(imap.FlagMDNSent):   imap.FlagMDNSent,
	string(imap.FlagJunk):      imap.FlagJunk,
	string(imap.FlagNotJunk):   imap.FlagNotJunk,
	string(imap.FlagPhishing):  imap.FlagPhishing,
	string(imap.FlagImportant): imap.FlagImportant,
	string(imap.FlagWildcard):  imap.FlagWildcard,
}

func InitializeEmail(userID *int64, emailConfig map[string]interface{}) error {
	driveDir, errD := GetDriveDirectory(userID)
	if errD != nil {
		util.Loggify(errD)
		return fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	email := emailConfig["email"].(string)
	emailDirectory := util.PointerString(filepath.Join(*directory, email))
	errA = files.CreateIfNotExists(*emailDirectory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return fmt.Errorf("not able to reserve an email directory for: %s", email)
	}

	statusMetadataPath := filepath.Join(*emailDirectory, "metadata")
	statusMetadata, errSM := getStatusMetadata(&statusMetadataPath, &email)
	if errSM != nil || statusMetadata == nil {
		Mailboxes := make(map[string]map[string]interface{})
		mailboxes := emailConfig["mailboxes"].([]interface{})
		for _, mailbox := range mailboxes {
			label, okL := mailbox.(map[string]interface{})["label"]
			numMessagesSource, okNM := mailbox.(map[string]interface{})["numMessages"]
			numMessagesSourceInt64, eNMSN := util.I2Num(numMessagesSource)
			syncNumber, okSN := mailbox.(map[string]interface{})["syncNumber"]
			syncNumberInt64, eSNI64 := util.I2Num(syncNumber)
			all, okA := mailbox.(map[string]interface{})["all"]
			if !okL || !okNM || !okSN || !okA || eNMSN != nil || eSNI64 != nil {
				continue
			}
			NumMessages := 1
			if !all.(bool) {

				NumMessages = int(numMessagesSourceInt64 - syncNumberInt64)
				if NumMessages < 1 {
					NumMessages = 1
				}
			}
			Mailboxes[label.(string)] = map[string]interface{}{
				"Count":       0,
				"NumMessages": NumMessages,
			}
		}
		statusMetadata = &structs.EmailStatusMetadata{
			Email:     &email,
			Mailboxes: &Mailboxes,
		}

		statusMetadataJson, sMJErr := statusMetadata.Json()
		if sMJErr != nil {
			util.Loggify(sMJErr)
		}
		errMet := os.WriteFile(statusMetadataPath, []byte(*statusMetadataJson), 0644)
		if errMet != nil {
			util.Loggify(errMet)
			return fmt.Errorf("not able to update the metadata for email: %s", email)
		}
	}

	return nil
}

func Send(userID *int64, Id *string, From *string, TO *string, CC *string, BCC *string, Subject *string, Content *string) (bool, error) {
	tOAr := strings.Split(*TO, ";")
	cCAr := strings.Split(*CC, ";")
	bCCAr := strings.Split(*BCC, ";")
	message := email.Message{
		From:    From,
		To:      &tOAr,
		CC:      &cCAr,
		BCC:     &bCCAr,
		Subject: Subject,
		Body:    Content,
	}
	user, errU := database.GetUser(userID, nil, nil, true, util.PointerBoolean(false), true, true)
	if errU != nil {
		util.Loggify(errU)
		return false, fmt.Errorf("not able to locate the user, contact the support")
	}
	if *user.AmperId != *business.AmperId() {
		errB := fmt.Errorf("the amper instance you are requesting doesn't belong to the user")
		util.Loggify(errB)
		return false, errB
	}
	errI := user.Initialize(true)
	if errI != nil {
		util.Loggify(errI)
		return false, fmt.Errorf("not able to initialize the user configuration, contact the support")
	}
	var emailConfig map[string]interface{}
	if user.Emails != nil {
		for _, emailItem := range *user.Emails {
			if emailItem["email"] == *From {
				emailConfig = emailItem
				break
			}
		}
	}
	if emailConfig == nil {
		return false, fmt.Errorf("not able to locate the email configuration for the email: %s", *From)
	}

	password, okP := emailConfig["password"].(string)
	if !okP && len(password) < 1 {
		//skip checking the email not a valid email address
		return false, fmt.Errorf("the password supplied is too short for email: %s", *From)
	}
	passwordEncrypted, _ := hex.DecodeString(password)
	passwordDecrypted := string(crypto.Decrypt(passwordEncrypted, *user.Password))
	serverName, port, auth, errSN := getSmtpServerNameParts(user.ID, *From)
	if errSN != nil {
		util.Loggify(errSN)
		return false, fmt.Errorf("failed to lookup SMTP server for email: %s", *From)
	}

	success, errSG := message.SendEmail(*From, passwordDecrypted, serverName, port, auth)
	if errSG != nil || !success {
		util.Loggify(errSG)
		return false, fmt.Errorf("not able to send the email from: %s", *From)
	}

	if Id != nil {
		idsSplit := strings.Split(*Id, "_")
		if len(idsSplit) == 4 {
			timeInt, errTimeInt := strconv.ParseInt(idsSplit[0], 10, 64)
			if errTimeInt == nil {
				messageTime := time.Unix(timeInt, 0)
				year := messageTime.Year()
				month := messageTime.Month()

				driveDir, errD := GetDriveDirectory(userID)
				if errD == nil {
					directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH, *From, "Drafts", "success", strconv.Itoa(year), strconv.Itoa(int(month)), *Id))
					if files.Exists(*directory) {
						files.RemoveAll(*directory)
					}
				}
			}
		}
	}
	return true, nil
}

func SaveEmailDraft(userID *int64, Id *string, From *string, TO *string, CC *string, BCC *string, Subject *string, Content *string) (*string, error) {
	driveDir, errD := GetDriveDirectory(userID)
	if errD != nil {
		util.Loggify(errD)
		return nil, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return nil, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	emailDirectory := util.PointerString(filepath.Join(*directory, *From))
	errA = files.CreateIfNotExists(*emailDirectory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return nil, fmt.Errorf("not able to reserve an email directory for: %s", *From)
	}
	draftDirectory := util.PointerString(filepath.Join(*emailDirectory, "Drafts"))
	statusMetadataPath := filepath.Join(*emailDirectory, "metadata")
	statusMetadata, errSM := getStatusMetadata(&statusMetadataPath, From)
	if errSM != nil {
		return nil, fmt.Errorf("not able to locate the status metadata for email: %s", *From)
	}
	_, okD := (*statusMetadata.Mailboxes)["Drafts"]
	if !okD {
		(*statusMetadata.Mailboxes)["Drafts"] = map[string]interface{}{
			"Local": true,
			"Count": 0,
		}
		statusMetadataJson, sMJErr := statusMetadata.Json()
		if sMJErr != nil {
			util.Loggify(sMJErr)
			return nil, fmt.Errorf("not able to format the metadata for email: %s", *From)
		}
		errMet := os.WriteFile(statusMetadataPath, []byte(*statusMetadataJson), 0644)
		if errMet != nil {
			util.Loggify(errMet)
			return nil, fmt.Errorf("not able to update the metadata for email: %s", *From)
		}

		errA = files.CreateIfNotExists(*draftDirectory, 0755)
		if errA != nil {
			util.Loggify(errA)
			return nil, fmt.Errorf("not able to reserve a draft directory for: %s", *From)
		}
	}

	var messageTime time.Time
	var flags []imap.Flag
	draftMailbox := (*statusMetadata.Mailboxes)["Drafts"]
	count, errCF := util.I2Num(draftMailbox["Count"])
	if errCF != nil {
		util.Loggify(errCF)
		return nil, fmt.Errorf("not able to format the count field of the draft metadata for email: %s", *From)
	}
	fromAddress := imapmemserver.ParseAddressList(*From)
	toAddress := imapmemserver.ParseAddressList(*TO)
	newDraft := true
	if Id == nil || len(*Id) < 1 {
		messageTime = time.Now()
		Id = util.PointerString(strconv.FormatInt(messageTime.Unix(), 10) + "_" + strconv.FormatInt(count, 10) + "_" + strconv.FormatInt(count, 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	} else {
		newDraft = false
		idsSplit := strings.Split(*Id, "_")
		if len(idsSplit) != 4 {
			return nil, fmt.Errorf("not able to format the id of draft %s", *Id)
		}
		timeInt, errTimeInt := strconv.ParseInt(idsSplit[0], 10, 64)
		if errTimeInt != nil {
			util.Loggify(errTimeInt)
			return nil, fmt.Errorf("not able to parse the id of draft %s", *Id)
		}
		messageTime = time.Unix(timeInt, 0)
	}
	envelop := imap.Envelope{
		Date:      messageTime,
		Subject:   *Subject,
		From:      fromAddress,
		Sender:    fromAddress,
		ReplyTo:   fromAddress,
		To:        toAddress,
		Cc:        imapmemserver.ParseAddressList(*CC),
		Bcc:       imapmemserver.ParseAddressList(*BCC),
		MessageID: *Id,
	}
	emailMessage := structs.Email{
		ID:           Id,
		SeqNum:       uint32(count),
		Flags:        flags,
		Envelope:     &envelop,
		InternalDate: messageTime,
		RFC822Size:   0,
		UID:          0,
		BodyHTML:     Content,
	}
	emailRelativeDirectory := filepath.Join(EMAIL_PATH, *From, "Drafts")
	if newDraft {
		errSS := saveSuccess(userID, &emailRelativeDirectory, emailMessage, make([]structs.Attachment, 0))
		if errSS != nil {
			util.Loggify(errSS)
			return nil, fmt.Errorf("not able to save the draft for email: %s", *From)
		}
	} else {
		errUM := updateMetadata(userID, &emailRelativeDirectory, emailMessage)
		if errUM != nil {
			util.Loggify(errUM)
			return nil, fmt.Errorf("not able to save the draft for email: %s", *Id)
		}
	}

	draftMailbox["Count"] = count + 1
	(*statusMetadata.Mailboxes)["Drafts"] = draftMailbox
	statusMetadataJson, sMJErr := statusMetadata.Json()
	if sMJErr != nil {
		util.Loggify(sMJErr)
		return nil, fmt.Errorf("not able to format the metadata for email: %s", *From)
	}
	errMet := os.WriteFile(statusMetadataPath, []byte(*statusMetadataJson), 0644)
	if errMet != nil {
		util.Loggify(errMet)
		return nil, fmt.Errorf("not able to update the metadata for email: %s", *From)
	}
	return Id, nil
}

func updateMetadata(userId *int64, path *string, email structs.Email) (err error) {
	time := email.Envelope.Date

	year := time.Year()
	month := time.Month()
	messageId := email.Envelope.MessageID
	emailDirectory := filepath.Join(*path, "success", strconv.Itoa(year), strconv.Itoa(int(month)), *files.Name(email.ID))

	amperFiles, errF := FetchFiles(userId, &emailDirectory)
	if errF != nil {
		util.Loggify(errF)
		return fmt.Errorf("not able to fetch files for the email metadata: %s", *email.ID)
	}
	metadataRemoved := false
	for _, amperFile := range *amperFiles {
		if *amperFile.Name == "Metadata" && *amperFile.Type == "application/email" {
			success, errRF := RemoveFile(userId, &emailDirectory, amperFile.Id)
			if !success || errRF != nil {
				util.Loggify(errRF)
				return fmt.Errorf("not able to remove the previos email state for %s", emailDirectory)
			} else {
				metadataRemoved = true
			}
		}
	}
	if !metadataRemoved {
		return fmt.Errorf("not able to remove the email state metadata for %s", emailDirectory)
	}
	emailMetadata, errEM := email.Json()
	if errEM != nil {
		util.Loggify(errEM)
		return fmt.Errorf("not able to format the email metadata, skipping for message: %s", messageId)
	}
	emailMetadataByte := []byte(*emailMetadata)
	start := int64(0)
	size := int64(len(emailMetadataByte))

	success, _, _, errUP := UploadChunk(userId, util.PointerString(EMAIL_METADATA), nil, &size, util.PointerString("application/email"), nil, &emailMetadataByte, &start, &emailDirectory)
	if errUP != nil || !success {
		util.Loggify(errUP)
		return fmt.Errorf("not able to store the email metadata, skipping for message: %s", messageId)
	}
	return err
}

func Mailboxes(userID *int64, Email *string) (*[]structs.Mailbox, error) {
	result := make([]structs.Mailbox, 0)
	driveDir, errD := GetDriveDirectory(userID)
	if errD != nil {
		util.Loggify(errD)
		return nil, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return nil, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	emailDirectory := util.PointerString(filepath.Join(*directory, *Email))
	errA = files.CreateIfNotExists(*emailDirectory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return nil, fmt.Errorf("not able to reserve an email directory for: %s", *Email)
	}

	statusMetadataPath := filepath.Join(*emailDirectory, "metadata")
	statusMetadata, errSM := getStatusMetadata(&statusMetadataPath, Email)
	if errSM != nil || statusMetadata == nil {
		util.Loggify(errSM)
		return nil, fmt.Errorf("not able to retrieve the status metadata, please contact the support")
	}
	for label, mailBox := range *statusMetadata.Mailboxes {
		Count, errC := util.I2Num(mailBox["Count"])
		NumMessages, _ := util.I2Num(mailBox["NumMessages"])
		if errC == nil {
			result = append(result, structs.Mailbox{
				Label:       label,
				Count:       int(Count),
				NumMessages: int(NumMessages),
			})
		}
	}
	return &result, nil
}

func ConfigureEmail(userID *int64, Email *string, Password *string) (*[]structs.Mailbox, error) {
	serverName, errSN := getImapServerName(userID, *Email)
	if errSN != nil {
		util.Loggify(errSN)
		return nil, fmt.Errorf("failed to lookup IMAP server for email: %s", *Email)
	}
	connection, errDTLS := imapclient.DialTLS(serverName, nil)
	if errDTLS != nil {
		util.Loggify(errDTLS)
		return nil, fmt.Errorf("failed to dial IMAP server for email: %s", *Email)
	}

	if errL := connection.Login(*Email, *Password).Wait(); errL != nil {
		util.Loggify(errL)
		return nil, fmt.Errorf("failed to log in into the email server for: %s", *Email)
	}
	mailboxes, errL := connection.List("", "%", nil).Collect()
	if errL != nil {
		util.Loggify(errL)
		return nil, fmt.Errorf("failed to dial IMAP server for email: %s", *Email)
	}
	result := make([]structs.Mailbox, 0)
	for _, mailbox := range mailboxes {
		mailboxSelected, errS := connection.Select(mailbox.Mailbox, nil).Wait()
		if errS != nil && mailboxSelected.NumMessages >= 0 {
			util.Loggify(errS)
			continue
		}
		result = append(result, structs.Mailbox{
			Label:       mailbox.Mailbox,
			NumMessages: int(mailboxSelected.NumMessages),
			SyncNumber:  100,
			All:         false,
		})
	}
	return &result, nil
}

func FlagEmails(userId *int64, Emails *[]map[string]interface{}, Box *string, Flags *[]string) (result bool, err error) {
	for _, flag := range *Flags {
		_, okF := EMAIL_FLAGS[flag]
		if !okF {
			return false, fmt.Errorf("the requested action can't be completed, since the flag %s is not recongnized by us", flag)
		}
	}
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), true, false)
	if errU != nil {
		util.Loggify(errU)
		return false, fmt.Errorf("not able to locate the user, contact the support")
	}
	if *user.AmperId != *business.AmperId() {
		errB := fmt.Errorf("the amper instance you are requesting doesn't belong to the user")
		util.Loggify(errB)
		return false, errB
	}
	errI := user.Initialize(true)
	if errI != nil {
		util.Loggify(errI)
		return false, fmt.Errorf("not able to initialize the user configuration, contact the support")
	}
	driveDir, errD := GetDriveDirectory(userId)
	if errD != nil {
		util.Loggify(errD)
		return false, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return false, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	EmailGroups := make(map[string][]map[string]interface{})
	for _, emailItem := range *Emails {
		email, okE := emailItem["email"].(string)
		if okE && len(email) > 0 {
			emailGroup := EmailGroups[email]
			if emailGroup == nil {
				emailGroup = make([]map[string]interface{}, 0)
			}
			EmailGroups[email] = append(emailGroup, emailItem)
		}
	}
	result = true
	for emailAddress, emailItems := range EmailGroups {
		success, errME := flagEmails(user, &emailAddress, Box, Flags, emailItems, directory)
		if !success || errME != nil {
			err = fmt.Errorf("not able to flag some of the email items for inbox email: %s", emailAddress)
			result = false
		}
	}
	return result, err
}

func flagEmails(user *structs.User, email *string, Box *string, Flags *[]string, emailItems []map[string]interface{}, directory *string) (result bool, err error) {
	emailFilePath := filepath.Join(*directory, *email)
	emailFileRelativePath := filepath.Join(EMAIL_PATH, *email)
	BoxPath := filepath.Join(emailFilePath, *Box)
	BoxRelativePath := filepath.Join(emailFileRelativePath, *Box)
	BoxPath = filepath.Join(BoxPath, "success")
	BoxRelativePath = filepath.Join(BoxRelativePath, "success")

	for _, emailItem := range emailItems {
		id, IOk := emailItem["id"]
		if IOk {
			idParts := strings.Split(id.(string), "_")
			if len(idParts) == 4 {
				secondTime, errNT := strconv.ParseInt(idParts[0], 10, 64)
				if errNT == nil {
					date := time.Unix(secondTime, 0)
					emailInboxItemPath := filepath.Join(BoxRelativePath, strconv.Itoa(date.Year()), strconv.Itoa(int(date.Month())), *files.Name(util.PointerString(id.(string))))
					amperFiles, errFF := FetchFiles(user.ID, util.PointerString(emailInboxItemPath))
					if errFF == nil {
						for _, amperFile := range *amperFiles {
							if *amperFile.Name == "Metadata" && *amperFile.Type == "application/email" {
								emailItemFilePath := filepath.Join(BoxPath, strconv.Itoa(date.Year()), strconv.Itoa(int(date.Month())), *files.Name(util.PointerString(id.(string))), *amperFile.Id, "file")
								data, errMet := os.ReadFile(emailItemFilePath)
								if errMet == nil {
									readEmailItem := structs.Email{}
									errUM := json.Unmarshal(data, &readEmailItem)
									if errUM != nil {
										util.Loggify(errUM)
										err = fmt.Errorf("not able to parse an email file on path %s to struct", emailItemFilePath)
									} else {
										flags := readEmailItem.Flags
										if flags == nil {
											flags = make([]imap.Flag, 0)
										}
										for _, flag := range *Flags {
											flagImap, okF := EMAIL_FLAGS[flag]
											if okF {
												flags = append(flags, flagImap)
											}
										}
										readEmailItem.Flags = flags

										//save the content by replacing it
										emailMetadata, errEM := readEmailItem.Json()
										if errEM != nil || emailMetadata == nil {
											util.Loggify(errEM)
											err = fmt.Errorf("not able to flag the email %s", *readEmailItem.ID)
											break
										}
										emailMetadataByte := []byte(*emailMetadata)
										start := int64(0)
										size := int64(len(emailMetadataByte))

										success, _, _, errUP := UploadChunk(user.ID, util.PointerString(EMAIL_METADATA), nil, &size, util.PointerString("application/email"), nil, &emailMetadataByte, &start, util.PointerString(emailInboxItemPath))
										if errUP != nil || !success {
											util.Loggify(errUP)
											err = fmt.Errorf("not able to flag the email %s", *readEmailItem.ID)
											break
										} else {
											rF, rFErr := RemoveFile(user.ID, util.PointerString(emailInboxItemPath), amperFile.Id)
											if rFErr != nil || !rF {
												util.Loggify(rFErr)
												err = fmt.Errorf("not able to flag the email %s", *readEmailItem.ID)
												break
											} else {
												result = true
											}
										}
									}
								} else {
									util.Loggify(errMet)
									err = fmt.Errorf("not able to read an email file on path %s", emailItemFilePath)
								}
							}
						}
					} else {
						util.Loggify(errFF)
						err = fmt.Errorf("not able to fetch emails on path %s", emailInboxItemPath)
					}
				} else {
					util.Loggify(errNT)
					err = fmt.Errorf("not able to read an email file on path %s, the first part of id should be a date format: %s", BoxPath, idParts[0])
				}
			} else {
				err = fmt.Errorf("not able to read an email file on path %s, the id is of a wrong format: %s", BoxPath, id.(string))
			}
		} else {
			err = fmt.Errorf("not able to read an email file on path %s", BoxPath)
		}
	}
	return result, err
}

func MoveEmails(userId *int64, Emails *[]map[string]interface{}, From *string, To *string) (result bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), true, false)
	if errU != nil {
		util.Loggify(errU)
		return false, fmt.Errorf("not able to locate the user, contact the support")
	}
	if *user.AmperId != *business.AmperId() {
		errB := fmt.Errorf("the amper instance you are requesting doesn't belong to the user")
		util.Loggify(errB)
		return false, errB
	}
	errI := user.Initialize(true)
	if errI != nil {
		util.Loggify(errI)
		return false, fmt.Errorf("not able to initialize the user configuration, contact the support")
	}
	driveDir, errD := GetDriveDirectory(userId)
	if errD != nil {
		util.Loggify(errD)
		return false, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return false, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	EmailGroups := make(map[string][]map[string]interface{})
	for _, emailItem := range *Emails {
		email, okE := emailItem["email"].(string)
		if okE && len(email) > 0 {
			emailGroup := EmailGroups[email]
			if emailGroup == nil {
				emailGroup = make([]map[string]interface{}, 0)
			}
			EmailGroups[email] = append(emailGroup, emailItem)
		}
	}
	result = true
	for emailAddress, emailItems := range EmailGroups {
		success, errME := moveEmails(user, &emailAddress, emailItems, directory, 0, From, 0, To)
		if !success || errME != nil {
			err = fmt.Errorf("not able to move some of the email items for inbox email: %s", emailAddress)
			result = false
		}
	}
	/*Emails := make([]map[string]string, 0)
	if user.Emails != nil {
		for _, emailItem := range *user.Emails {
			if *Email == ALL_EMAIL || emailItem["email"] == *Email {
				Emails = append(Emails, emailItem)
			}
		}
	}*/
	return result, err
}

func moveEmails(user *structs.User, email *string, emailItems []map[string]interface{}, directory *string, fromState int, From *string, toState int, To *string) (result bool, err error) {
	result = true
	emailFilePath := filepath.Join(*directory, *email)

	fromPath := filepath.Join(emailFilePath, *From)
	toPath := filepath.Join(emailFilePath, *To)
	if *From == "inbox" {
		fromPath = filepath.Join(fromPath, "success")
	}
	if *To == "inbox" {
		toPath = filepath.Join(toPath, "success")
	}
	for _, emailItem := range emailItems {
		envelope, eOk := emailItem["envelope"].(map[string]interface{})
		if eOk {
			date, dOk := envelope["date"]
			if dOk {
				date, errT := datetime.ParseDateTime(util.PointerString(date.(string)))
				if errT == nil {
					id, IOk := emailItem["id"]
					if IOk {
						emailInboxItemPath := filepath.Join(fromPath, strconv.Itoa(date.Year()), strconv.Itoa(int(date.Month())), *files.Name(util.PointerString(id.(string))))
						emailTrashItemPath := filepath.Join(toPath, strconv.Itoa(date.Year()), strconv.Itoa(int(date.Month())))
						errA := files.CreateIfNotExists(emailTrashItemPath, 0755)
						if errA != nil {
							util.Loggify(errA)
							err = fmt.Errorf("not able to initiate a email year and month directory for path %s", emailTrashItemPath)
							result = false
						} else {
							errRen := os.Rename(emailInboxItemPath, filepath.Join(emailTrashItemPath, *files.Name(util.PointerString(id.(string)))))
							if errRen != nil {
								util.Loggify(errRen)
								result = false
								err = fmt.Errorf("not able to move email to the trash bin")
							}
						}
					} else {
						err = fmt.Errorf("email with id: %v has no valid messageId provided", emailItem["id"])
						result = false
					}
				} else {
					util.Loggify(errT)
					err = fmt.Errorf("email with id: %v has no valid date provided", emailItem["id"])
					result = false
				}
			} else {
				err = fmt.Errorf("email with id: %v has no date provided", emailItem["id"])
				result = false
			}
		}
	}
	return result, err
}

func FetchEmails(userId *int64, Email *string, Box *string, Search *string, Start *int, Limit *int, Pointer *structs.PagePointer) (result []structs.Email, resultTotalCount uint32, err error) {
	if Box == nil {
		Box = util.PointerString("INBOX")
	}
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), true, false)
	if errU != nil {
		util.Loggify(errU)
		return nil, resultTotalCount, fmt.Errorf("not able to locate the user, contact the support")
	}
	if *user.AmperId != *business.AmperId() {
		errB := fmt.Errorf("the amper instance you are requesting doesn't belong to the user")
		util.Loggify(errB)
		return nil, resultTotalCount, errB
	}
	errI := user.Initialize(true)
	if errI != nil {
		util.Loggify(errI)
		return nil, resultTotalCount, fmt.Errorf("not able to initialize the user configuration, contact the support")
	}
	driveDir, errD := GetDriveDirectory(userId)
	if errD != nil {
		util.Loggify(errD)
		return nil, resultTotalCount, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return nil, resultTotalCount, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	Emails := make([]map[string]interface{}, 0)
	if user.Emails != nil {
		for _, emailItem := range *user.Emails {
			if *Email == ALL_EMAIL || emailItem["email"].(string) == *Email {
				Emails = append(Emails, emailItem)
			}
		}
	}
	if len(Emails) > 0 {
		for _, emailItem := range Emails {
			var emails []structs.Email
			var errE error
			var totalCount uint32
			emails, totalCount, errE = fetchEmail(user, util.PointerString(emailItem["email"].(string)), util.PointerString(filepath.Join(*directory, emailItem["email"].(string))), Box, Search, Start, Limit, Pointer)
			if totalCount > resultTotalCount {
				resultTotalCount = totalCount
			}
			if errE != nil {
				util.Loggify(errE)
				err = fmt.Errorf("not able to fetch the emails for address: %s", emailItem["email"])
			} else {
				result = append(result, emails...)
			}
		}
	}
	return result, resultTotalCount, err
}

func fetchEmail(user *structs.User, email *string, directory *string, Box *string, Search *string, Start *int, Limit *int, Pointer *structs.PagePointer) (result []structs.Email, totalCount uint32, err error) {
	result = make([]structs.Email, 0)
	if !files.Exists(*directory) {
		return result, 0, err
	}
	boxDirectory := filepath.Join(*directory, *Box, "success")
	files, errRD := ioutil.ReadDir(boxDirectory)
	if errRD != nil {
		util.Loggify(errRD)
	}
	var startId *string
	var startYear *int
	var startMonth *int
	currentPage := *Limit / (*Limit - *Start)
	if Pointer != nil && Pointer.Pages != nil && *Start != 0 {
		pages := make([]int, 0)
		for page, _ := range Pointer.Pages {
			pages = append(pages, page)
		}
		sort.Ints(pages)
		for _, page := range pages {
			if currentPage > page {
				startId = util.PointerString(Pointer.Pages[page])
			}
		}
	}
	if startId != nil {
		startIdParts := strings.Split(*startId, "_")
		if len(startIdParts) == 4 {
			secondTime, errNT := strconv.ParseInt(startIdParts[0], 10, 64)
			if errNT == nil {
				startTime := time.Unix(secondTime, 0)
				startYear = util.PointerInt(startTime.Year())
				startMonth = util.PointerInt(int(startTime.Month()))
			}
		}
	}
	yearsFolders := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			yearsFolders = append(yearsFolders, file.Name())
		}
	}
	sort.Strings(yearsFolders)
	boxRelativeDirectory := filepath.Join(EMAIL_PATH, *email, *Box, "success")

	statusMetadata, errSM := getStatusMetadata(util.PointerString(filepath.Join(*directory, "metadata")), email)
	if errSM != nil {
		util.Loggify(errSM)
		return nil, 0, fmt.Errorf("not able to retrieve the status metadata for email: %s", *email)
	}
	mailbox, okMB := (*statusMetadata.Mailboxes)[*Box]
	if !okMB || len(mailbox) < 1 {
		return nil, 0, fmt.Errorf("not able to retrieve the status metadata for email: %s and mailbox: %s", *email, *Box)
	}
	Count, errC := util.I2Num(mailbox["Count"])
	if errC != nil {
		util.Loggify(errC)
		return nil, 0, fmt.Errorf("not able to retrieve the status metadata for email: %s and mailbox: %s", *email, *Box)
	}
	CountUint := uint32(Count)
	for i := len(yearsFolders) - 1; i >= 0; i-- {
		yearFolder := yearsFolders[i]
		yearFolderInt, errYI := strconv.Atoi(yearFolder)
		if errYI != nil || (startYear != nil && yearFolderInt > *startYear) {
			continue
		}
		files, errRD = ioutil.ReadDir(filepath.Join(boxDirectory, yearFolder))
		if errRD != nil {
			util.Loggify(errRD)
		}
		monthsFolders := make([]string, 0)
		for _, file := range files {
			if file.IsDir() {
				monthsFolders = append(monthsFolders, file.Name())
			}
		}
		sort.Strings(monthsFolders)
		for l := len(monthsFolders) - 1; l >= 0; l-- {
			monthFolder := monthsFolders[l]
			monthFolderInt, errMI := strconv.Atoi(monthFolder)
			if errMI != nil || (startYear != nil && monthFolderInt > *startMonth) {
				continue
			}

			monthDirectory := filepath.Join(boxDirectory, yearFolder, monthFolder)
			files, errRD = ioutil.ReadDir(monthDirectory)
			if errRD == nil {
				emailIds := make([]string, 0)
				for _, file := range files {
					if file.IsDir() {
						emailIds = append(emailIds, file.Name())
					}
				}
				sort.Strings(emailIds)
				for i := len(emailIds) - 1; i >= 0; i-- {
					emailId := emailIds[i]
					if startId != nil && strings.Compare(emailId, *startId) >= 0 {
						continue
					}
					emailItemPath := util.PointerString(filepath.Join(boxRelativeDirectory, yearFolder, monthFolder, emailId))
					amperFiles, errFF := FetchFiles(user.ID, emailItemPath)
					if errFF == nil {
						for _, amperFile := range *amperFiles {
							if *amperFile.Name == "Metadata" && *amperFile.Type == "application/email" {
								emailItemFilePath := filepath.Join(monthDirectory, emailId, *amperFile.Id, "file")
								data, errMet := os.ReadFile(emailItemFilePath)
								if errMet == nil {
									emailItem := structs.Email{}
									errUM := json.Unmarshal(data, &emailItem)
									emailItem.Email = email
									if errUM != nil {
										util.Loggify(errUM)
										err = fmt.Errorf("not able to parse an email file on path %s to struct", emailItemFilePath)
									} else {
										result = append(result, emailItem)
										if len(result) > *Limit-*Start-1 {
											return result, CountUint, err
										}
									}
								} else {
									util.Loggify(errMet)
									err = fmt.Errorf("not able to read an email file on path %s", emailItemFilePath)
								}
							}
						}
					} else {
						util.Loggify(errFF)
						err = fmt.Errorf("not able to read an email file on path %s", *emailItemPath)
						break
					}
				}
			} else {
				util.Loggify(errRD)
			}
		}
	}
	return result, CountUint, err
}

func CheckEmails(userId *int64, Email *string) (result bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), true, true)
	if errU != nil {
		util.Loggify(errU)
		return false, fmt.Errorf("not able to locate the user, contact the support")
	}
	if *user.AmperId != *business.AmperId() {
		errB := fmt.Errorf("the amper instance you are requesting doesn't belong to the user")
		util.Loggify(errB)
		return false, errB
	}
	errI := user.Initialize(true)
	if errI != nil {
		util.Loggify(errI)
		return false, fmt.Errorf("not able to initialize the user configuration, contact the support")
	}
	driveDir, errD := GetDriveDirectory(userId)
	if errD != nil {
		util.Loggify(errD)
		return false, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, EMAIL_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return false, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}

	Emails := make([]map[string]interface{}, 0)
	if user.Emails != nil {
		for _, emailItem := range *user.Emails {
			if *Email == ALL_EMAIL || emailItem["email"] == *Email {
				Emails = append(Emails, emailItem)
			}
		}
	}
	failedEmailsMap := make(map[string][]string)
	result = true
	if len(Emails) > 0 {
		for _, emailItem := range Emails {
			success, errChE := checkEmail(user, emailItem, directory, failedEmailsMap)
			if errChE != nil || !success {
				if err != nil {
					err = errors.New(err.Error() + ", " + emailItem["email"].(string))
				} else {
					err = fmt.Errorf("not able to check the incoming emails for %s", emailItem["email"])
				}

				util.Loggify(errChE)
				result = false
			}
		}
	}
	//TODO notify the user of failed emails once the notification framework is complete
	if len(failedEmailsMap) > 0 {

	}

	/*mailboxes, err := c.List("", "%", nil).Collect()
	if err != nil {
		log.Fatalf("failed to list mailboxes: %v", err)
	}
	log.Printf("Found %v mailboxes", len(m--ailboxes))
	for _, mbox := range mailboxes {
		log.Printf(" - %v", mbox.Mailbox)
	}*/
	return result, err
}

func getStatusMetadata(path *string, email *string) (*structs.EmailStatusMetadata, error) {
	if _, err := os.Stat(*path); !errors.Is(err, os.ErrNotExist) {
		data, errStatMet := os.ReadFile(*path)
		if errStatMet != nil {
			util.Loggify(errStatMet)
			return nil, fmt.Errorf("not able to read the status metadata from the directory: %s", *path)
		} else if len(data) > 0 {
			Mailboxes := make(map[string]map[string]interface{})
			statusMetadata := structs.EmailStatusMetadata{
				Email:     email,
				Mailboxes: &Mailboxes,
			}
			errP := statusMetadata.Parse(&data)
			if errP != nil {
				util.Loggify(errP)
				return nil, fmt.Errorf("not able to parse the status metadata from the directory: %s", *path)
			} else {
				return &statusMetadata, nil
			}
		}
	}
	return nil, fmt.Errorf("the status metadata was not found for email %s", *email)
}

func getSmtpServerNameParts(userId *int64, email string) (string, int, string, error) {
	emailSplit := strings.Split(email, "@")
	if len(emailSplit) != 2 {
		return "", -1, "", fmt.Errorf("wrongly formatted email address supplied: %s", email)
	}
	smtpS := GetSetting(userId, util.PointerString("amper.smtp"), nil)
	smtp := structs.Smtp{}
	json.Unmarshal([]byte(smtpS), &smtp)
	for _, domain := range smtp.Domains {
		if domain.Domain == emailSplit[1] {
			return domain.ServerName, domain.Port, domain.Auth, nil
		}
	}
	return "", -1, "", fmt.Errorf("not able to locat the domain server name in admin imap settings for email: %s", email)
}

func getImapServerNameParts(userId *int64, email string) (string, int, error) {
	emailSplit := strings.Split(email, "@")
	if len(emailSplit) != 2 {
		return "", -1, fmt.Errorf("wrongly formatted email address supplied: %s", email)
	}
	imapS := GetSetting(userId, util.PointerString("amper.imap"), nil)
	imap := structs.Imap{}
	json.Unmarshal([]byte(imapS), &imap)
	for _, domain := range imap.Domains {
		if domain.Domain == emailSplit[1] {
			return domain.ServerName, domain.Port, nil
		}
	}
	return "", -1, fmt.Errorf("not able to locat the domain server name in admin imap settings for email: %s", email)
}

func getImapServerName(userId *int64, email string) (string, error) {
	serverName, port, errSNP := getImapServerNameParts(userId, email)
	return serverName + ":" + strconv.Itoa(port), errSNP
}

func checkEmail(user *structs.User, emailItem map[string]interface{}, directory *string, failedEmailsMap map[string][]string) (bool, error) {
	email, okE := emailItem["email"].(string)
	if !okE || len(email) < 4 {
		//skip checking the email not a valid email address
		return false, fmt.Errorf("the email address '%s' supplied is too short", email)
	}

	configuredMailboxes, mE := emailItem["mailboxes"].([]interface{})
	if !mE || len(configuredMailboxes) < 1 {
		return false, fmt.Errorf("the email address '%s' supplied has no mailbox", email)
	}
	emailDirectory := util.PointerString(filepath.Join(*directory, email))
	errA := files.CreateIfNotExists(*emailDirectory, 0755)
	if errA != nil {
		return false, fmt.Errorf("not able to reserve an email directory for: %s", email)
	}
	//Make sure only a single thread is checking the same email at a time
	//the second attmpt should be canceled
	if emailLock.Locked(email) {
		util.Loggify(fmt.Errorf("another process is working for this email: %s", email))
		return true, nil
	} else {
		emailLock.Lock(email)
		defer emailLock.Unlock(email)
	}

	statusMetadataPath := filepath.Join(*emailDirectory, "metadata")
	statusMetadata, errSM := getStatusMetadata(&statusMetadataPath, &email)
	if errSM != nil || statusMetadata == nil {
		util.Loggify(errSM)
		return false, fmt.Errorf("not able to retrieve the status metadata for email: %s", email)
	}
	password, okP := emailItem["password"].(string)
	if !okP && len(password) < 1 {
		//skip checking the email not a valid email address
		return false, fmt.Errorf("the password supplied is too short for email: %s", email)
	}
	passwordEncrypted, _ := hex.DecodeString(password)
	passwordDecrypted := string(crypto.Decrypt(passwordEncrypted, *user.Password))
	serverName, errSN := getImapServerName(user.ID, email)
	if errSN != nil {
		util.Loggify(errSN)
		return false, fmt.Errorf("failed to lookup IMAP server for email: %s", email)
	}
	connection, errDTLS := imapclient.DialTLS(serverName, nil)
	if errDTLS != nil {
		util.Loggify(errDTLS)
		return false, fmt.Errorf("failed to dial IMAP server for email: %s", email)
	}

	if errL := connection.Login(email, passwordDecrypted).Wait(); errL != nil {
		util.Loggify(errL)
		return false, fmt.Errorf("failed to log in into the email server for: %s", email)
	}
	remoteMailboxes, errL := connection.List("", "%", nil).Collect()
	if errL != nil {
		util.Loggify(errL)
		return false, fmt.Errorf("failed to dial IMAP server for email: %s", email)
	}

	for _, remoteMailbox := range remoteMailboxes {
		var selectedMailbox = &remoteMailbox.Mailbox

		mailbox, errS := connection.Select(*selectedMailbox, nil).Wait()
		if errS != nil {
			util.Loggify(errS)
			continue
		}
		emailRelativeDirectory := filepath.Join(EMAIL_PATH, email, *selectedMailbox)
		var mailboxStatus map[string]interface{}
		if statusMetadata.Mailboxes != nil && (*statusMetadata.Mailboxes)[*selectedMailbox] != nil {
			mailboxStatus = (*statusMetadata.Mailboxes)[*selectedMailbox]
		} else {
			mailboxStatus = make(map[string]interface{})
			mailboxStatus["NumMessages"] = 1
			mailboxStatus["Count"] = 0
		}
		NumMessages, _ := util.I2Num(mailboxStatus["NumMessages"])
		NumMessagesUint := uint32(NumMessages)
		Count, _ := util.I2Num(mailboxStatus["Count"])

		currentCount := NumMessages
		newFetched := false
		for mailbox.NumMessages > NumMessagesUint {
			newFetched = true
			partition := uint32(EMAIL_PARTITION)
			if mailbox.NumMessages-NumMessagesUint < 100 {
				partition = mailbox.NumMessages - NumMessagesUint
			}
			seqSet := imap.SeqSetRange(NumMessagesUint+1, NumMessagesUint+partition)
			fetchItems := []imap.FetchItem{imap.FetchItemEnvelope, &imap.FetchItemBodySection{}}
			messages, errF := connection.Fetch(seqSet, fetchItems, nil).Collect()
			if errF != nil {
				util.Loggify(errF)
				util.Loggify(fmt.Errorf("failed to fetch batch of messages for email: %s", email))
				break
			}
			for _, message := range messages {
				currentCount++
				var messageTime time.Time
				if message.Envelope != nil {
					messageTime = message.Envelope.Date
				} else {
					messageTime = time.Now()
					message.Envelope.Date = messageTime
				}
				var flags []imap.Flag
				emailMessage := structs.Email{
					ID:           util.PointerString(strconv.FormatInt(messageTime.Unix(), 10) + "_" + strconv.FormatUint(uint64(currentCount), 10) + "_" + strconv.FormatUint(uint64(message.SeqNum), 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)),
					SeqNum:       message.SeqNum,
					Flags:        flags,
					Envelope:     message.Envelope,
					InternalDate: message.InternalDate,
					RFC822Size:   message.RFC822Size,
					UID:          message.UID,
				}

				attachments, errA := emailMessage.Process(user.ID, message.BodySection)
				if errA == nil {
					errSS := saveSuccess(user.ID, util.PointerString(emailRelativeDirectory), emailMessage, attachments)
					if errSS != nil {
						util.Loggify(errSS)
						emailDirectory, errSF := saveFail(user.ID, util.PointerString(emailRelativeDirectory), emailMessage, message.BodySection)
						if errSF != nil {
							util.Loggify(errSF)
						}
						if emailDirectory != nil {
							failedEmails := failedEmailsMap[email]
							if failedEmails == nil {
								failedEmails = make([]string, 0)
							}
							failedEmails = append(failedEmails, *emailDirectory)
							failedEmailsMap[email] = failedEmails
						}
					} else {
						Count++
					}
				} else {
					util.Loggify(errA)
					emailDirectory, errSF := saveFail(user.ID, util.PointerString(emailRelativeDirectory), emailMessage, message.BodySection)
					if errSF != nil {
						util.Loggify(errSF)
					}
					if emailDirectory != nil {
						failedEmails := failedEmailsMap[email]
						if failedEmails == nil {
							failedEmails = make([]string, 0)
						}
						failedEmails = append(failedEmails, *emailDirectory)
						failedEmailsMap[email] = failedEmails
					}
				}
			}
			NumMessagesUint += partition
		}
		if newFetched {
			mailboxStatus["NumMessages"] = int64(NumMessagesUint)
			mailboxStatus["Count"] = Count
			(*statusMetadata.Mailboxes)[*selectedMailbox] = mailboxStatus
			statusMetadataJson, sMJErr := statusMetadata.Json()
			if sMJErr != nil {
				util.Loggify(sMJErr)
				return false, fmt.Errorf("not able to format the metadata for email: %s", email)
			}
			errMet := os.WriteFile(statusMetadataPath, []byte(*statusMetadataJson), 0644)
			if errMet != nil {
				util.Loggify(errMet)
				return false, fmt.Errorf("not able to update the metadata for email: %s", email)
			}
		}
	}

	if connection != nil {
		LogoutAndClose(connection, email)
	}
	return true, nil
}

func AppendToFile(path *string, content []byte) error {
	f, errOF := os.OpenFile(*path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if errOF != nil {
		util.Loggify(errOF)
		return fmt.Errorf("not able to open a shortcuts file to append with path: %s", *path)
	}
	defer f.Close()

	if _, errW := f.Write(content); errW != nil {
		util.Loggify(errW)
		return fmt.Errorf("not able to append content to the end of the file with path: %s", *path)
	}
	return nil
}

func AddressesString(addresses []imap.Address) string {
	addressResult := make([]string, 0)
	for _, address := range addresses {
		addressResult = append(addressResult, address.String())
	}
	return strings.Join(addressResult, ampstrings.SEPERATOR)
}

func LogoutAndClose(c *imapclient.Client, email string) {
	if err := c.Logout().Wait(); err != nil {
		log.Printf("failed to logoutemail: %s, with error: %v", email, err)
	}
	if err := c.Close(); err != nil {
		log.Printf("failed to close connection for email: %s, with error: %v", email, err)
	}
}

const EMAIL_METADATA = "Metadata"

func saveSuccess(userId *int64, path *string, email structs.Email, attachments []structs.Attachment) (err error) {
	var timeR time.Time
	if email.Envelope != nil {
		timeR = email.Envelope.Date
	} else if email.Headers != nil && (*email.Headers)["Date"] != nil && len((*email.Headers)["Date"]) > 0 {
		date := (*email.Headers)["Date"]
		timeH, errT := mail.ParseDate(date[0].GetValue())
		if errT == nil {
			timeR = timeH
		} else {
			util.Loggify(errT)
			return fmt.Errorf("not able to determine the date of the email, skipping")
		}
	} else {
		return fmt.Errorf("not able to determine the date of the email, skipping")
	}

	year := timeR.Year()
	month := timeR.Month()
	messageId := email.Envelope.MessageID
	emailDirectory := filepath.Join(*path, "success", strconv.Itoa(year), strconv.Itoa(int(month)), *files.Name(email.ID))
	errA := files.RecreateIfExist(emailDirectory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return fmt.Errorf("not able to initiate a email year and month directory for path %s", emailDirectory)
	}

	emailMetadata, errEM := email.Json()
	if errEM != nil {
		util.Loggify(errEM)
		return fmt.Errorf("not able to format the email metadata, skipping for message: %s", messageId)
	}
	emailMetadataByte := []byte(*emailMetadata)
	start := int64(0)
	size := int64(len(emailMetadataByte))

	success, _, _, errUP := UploadChunk(userId, util.PointerString(EMAIL_METADATA), nil, &size, util.PointerString("application/email"), nil, &emailMetadataByte, &start, &emailDirectory)
	if errUP != nil || !success {
		util.Loggify(errUP)
		return fmt.Errorf("not able to store the email metadata, skipping for message: %s", messageId)
	}

	if len(attachments) > 0 {
		for _, attachment := range attachments {
			if len(attachment.Body) > 0 {
				start = int64(0)
				size = int64(len(attachment.Body))
				contentTypes := (*attachment.Headers)["Content-Type"]
				contentType := "application/octet-stream"
				if len(contentTypes) > 0 {
					contentType = contentTypes[0].GetValue()
				}
				success, _, _, errUP := UploadChunk(userId, attachment.Name, nil, &size, util.PointerString(contentType), nil, &attachment.Body, &start, &emailDirectory)
				if errUP != nil || !success {
					util.Loggify(errUP)
					err = fmt.Errorf("failed saving an email attachment, notify the user")
				}
			}
		}
	}
	return err
}

func saveFail(userId *int64, path *string, email structs.Email, BodySection map[*imap.FetchItemBodySection][]byte) (result *string, err error) {
	var timeR time.Time
	if email.Envelope != nil {
		timeR = email.Envelope.Date
	} else if email.Headers != nil && (*email.Headers)["Date"] != nil && len((*email.Headers)["Date"]) > 0 {
		date := (*email.Headers)["Date"]
		timeH, errT := mail.ParseDate(date[0].GetValue())
		if errT == nil {
			timeR = timeH
		} else {
			util.Loggify(errT)
			timeR = time.Now()
		}
	} else {
		timeR = time.Now()
	}

	year := timeR.Year()
	month := timeR.Month()
	emailDirectory := filepath.Join(*path, "failed", strconv.Itoa(year), strconv.Itoa(int(month)), *files.Name(email.ID))
	errA := files.RecreateIfExist(emailDirectory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return nil, fmt.Errorf("not able to initiate a email year and month directory for path %s", emailDirectory)
	}

	emailMetadata, errEM := email.Json()
	if errEM != nil {
		util.Loggify(errEM)
	}
	if emailMetadata != nil {
		emailMetadataByte := []byte(*emailMetadata)
		start := int64(0)
		size := int64(len(emailMetadataByte))

		success, _, _, errUP := UploadChunk(userId, util.PointerString(EMAIL_METADATA), nil, &size, util.PointerString("application/email"), nil, &emailMetadataByte, &start, &emailDirectory)
		if errUP != nil || !success {
			util.Loggify(errUP)
		}
	}

	index := 0
	for _, bodySection := range BodySection {
		size := int64(len(bodySection))
		if size > 0 {
			start := int64(0)
			success, _, _, errUP := UploadChunk(userId, util.PointerString("Body_Section_"+strconv.Itoa(index)), nil, &size, util.PointerString("application/octet-stream"), nil, &bodySection, &start, &emailDirectory)
			if errUP != nil || !success {
				util.Loggify(errUP)
			}
			index++
		}
	}
	return &emailDirectory, nil
}
