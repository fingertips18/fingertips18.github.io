package domain

import (
	"errors"
	"fmt"
	"strconv"
)

type Files struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Type     string  `json:"type"`
	CustomID *string `json:"customId,omitempty"`
}

type UploadRequest struct {
	Files              []Files `json:"files"`
	ACL                *string `json:"acl,omitempty"`
	Metadata           any     `json:"metadata,omitempty"`
	ContentDisposition *string `json:"contentDisposition,omitempty"`
}

func (i UploadRequest) Validate() error {
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

type UploadData struct {
	Key      string  `json:"key"`
	URL      string  `json:"url"`
	AppURL   string  `json:"appUrl"`
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	CustomId *string `json:"customId,omitempty"`
}

type UploadError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UploadFileResponse struct {
	Data  *UploadData  `json:"data"`
	Error *UploadError `json:"error"`
}

type UploadResponse struct {
	Data []UploadFileResponse `json:"data"`
}

func (r UploadResponse) Validate() error {
	if len(r.Data) == 0 {
		return errors.New("uploadthing: response returned no files")
	}

	for idx, f := range r.Data {
		prefix := "uploadthing: data[" + strconv.Itoa(idx) + "]."

		// Check if this file had an error
		if f.Error != nil {
			return fmt.Errorf("%s error: %s (code: %s)",
				prefix, f.Error.Message, f.Error.Code)
		}

		// Validate success data
		if f.Data == nil {
			return errors.New(prefix + ": missing data")
		}
		if f.Data.Key == "" {
			return errors.New(prefix + ".key missing")
		}
		if f.Data.URL == "" {
			return errors.New(prefix + ".url missing")
		}

	}

	return nil
}
