package domain

import (
	"errors"
	"fmt"
	"strconv"
)

type File struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Type     string  `json:"type"`
	CustomID *string `json:"customId,omitempty"`
}

type ImageUploadRequest struct {
	Files              []File  `json:"files"`
	ACL                *string `json:"acl,omitempty"`
	Metadata           any     `json:"metadata,omitempty"`
	ContentDisposition *string `json:"contentDisposition,omitempty"`
}

func (i ImageUploadRequest) Validate() error {
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

	// UploadThing ACL accepts: private, public-read
	if i.ACL != nil && *i.ACL != "" {
		if *i.ACL != "public-read" && *i.ACL != "private" {
			return errors.New("acl must be 'public-read' or 'private'")
		}
	}

	// UploadThing ContentDisposition accepts: inline, attachment
	if i.ContentDisposition != nil && *i.ContentDisposition != "" {
		if *i.ContentDisposition != "inline" && *i.ContentDisposition != "attachment" {
			return errors.New("contentDisposition must be 'inline' or 'attachment'")
		}
	}

	return nil
}

type ImageUploadFile struct {
	Key                string         `json:"key"`
	FileName           string         `json:"fileName"`
	FileType           string         `json:"fileType"`
	FileUrl            string         `json:"fileUrl"`
	ContentDisposition string         `json:"contentDisposition"`
	PollingJwt         string         `json:"pollingJwt"`
	PollingUrl         string         `json:"pollingUrl"`
	CustomId           *string        `json:"customId,omitempty"`
	URL                string         `json:"url"` // signed URL to upload
	Fields             map[string]any `json:"fields,omitempty"`
}

type ImageUploadResponse struct {
	Data []ImageUploadFile `json:"data"`
}

func (r ImageUploadResponse) Validate() error {
	if len(r.Data) == 0 {
		return errors.New("uploadthing: response returned no files")
	}

	for idx, f := range r.Data {
		prefix := fmt.Sprintf("uploadthing: data[%d]", idx)

		if f.Key == "" {
			return fmt.Errorf("%s.key missing", prefix)
		}
		if f.URL == "" {
			return fmt.Errorf("%s.url missing", prefix)
		}
	}

	return nil
}
