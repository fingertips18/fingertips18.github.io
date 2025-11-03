package domain

import "errors"

type PageView struct {
	PageLocation string `json:"location"`
	PageTitle    string `json:"title"`
}

func (p PageView) Validate() error {
	if p.PageLocation == "" {
		return errors.New("pageLocation missing")
	}
	if p.PageTitle == "" {
		return errors.New("pageTitle missing")
	}

	return nil
}
