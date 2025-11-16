package domain

import (
	"errors"
	"strconv"
)

type Files struct {
	Name     string  `json:"name"`
	Size     int32   `json:"size"`
	Type     string  `json:"type"`
	CustomID *string `json:"custom_id,omitempty"`
}

type UploadthingUploadRequest struct {
	Files              []Files `json:"files"`
	ACL                string  `json:"acl"`
	Metadata           any     `json:"metadata"`
	ContentDisposition string  `json:"contentDisposition"`
}

func (i UploadthingUploadRequest) Validate() error {
	// Files cannot be empty
	if len(i.Files) == 0 {
		return errors.New("files missing")
	}

	for idx, f := range i.Files {
		if f.Name == "" {
			return errors.New("file[" + strconv.Itoa(idx) + "]: name missing")
		}
		if f.Size <= 0 {
			return errors.New("file[" + strconv.Itoa(idx) + "]: size invalid")
		}
		if f.Type == "" {
			return errors.New("file[" + strconv.Itoa(idx) + "]: type missing")
		}
	}

	// ACL is required (UploadThing accepts: private, public-read)
	if i.ACL == "" {
		return errors.New("acl missing")
	}

	// contentDisposition must not be empty (inline / attachment)
	if i.ContentDisposition == "" {
		return errors.New("contentDisposition missing")
	}

	return nil
}

type UploadthingFile struct {
	Key                string         `json:"key"`
	FileName           string         `json:"fileName"`
	FileType           string         `json:"fileType"`
	FileUrl            string         `json:"fileUrl"`
	ContentDisposition string         `json:"contentDisposition"`
	PollingJwt         string         `json:"pollingJwt"`
	PollingUrl         string         `json:"pollingUrl"`
	CustomId           *string        `json:"customId"`
	Url                string         `json:"url"`
	Fields             map[string]any `json:"fields"`
}

type UploadthingUploadResponse struct {
	Data []UploadthingFile `json:"data"`
}

func (r UploadthingUploadResponse) Validate() error {
	if len(r.Data) == 0 {
		return errors.New("uploadthing: response returned no files")
	}

	for idx, f := range r.Data {
		prefix := "uploadthing: data[" + strconv.Itoa(idx) + "]."

		if f.Key == "" {
			return errors.New(prefix + "key missing")
		}
		if f.FileName == "" {
			return errors.New(prefix + "fileName missing")
		}
		if f.FileType == "" {
			return errors.New(prefix + "fileType missing")
		}
		if f.FileUrl == "" {
			return errors.New(prefix + "fileUrl missing")
		}
		if f.Url == "" {
			return errors.New(prefix + "url missing")
		}
	}

	return nil
}
