package model

import (
	net_url "net/url"
)

// //////////////////////////////////////////////////
// url

type Url string

func (o Url) IsValid(pathValidator PathValidator) bool {
	return o.Validate(pathValidator) == nil
}

func (o Url) Validate(pathValidator PathValidator) error {
	if o == "" {
		return nil
	} else if parsed, err := net_url.Parse(string(o)); err != nil {
		return err
	} else if IsLocalPath(parsed) && pathValidator != nil {
		if err := pathValidator(string(o)); err != nil {
			return err
		}
	}
	return nil
}

func (o Url) IsEmpty() bool {
	return o == ""
}

func (o Url) IsRemote() bool {
	if o.IsEmpty() {
		return false
	}
	parsed, _ := net_url.Parse(string(o))
	return IsRemotePath(parsed)
}

func IsRemotePath(url *net_url.URL) bool {
	return url != nil && url.Scheme != "" && url.Host != ""
}

func (o Url) IsLocal() bool {
	if o.IsEmpty() {
		return false
	}
	parsed, _ := net_url.Parse(string(o))
	return IsLocalPath(parsed)
}

func IsLocalPath(url *net_url.URL) bool {
	return url != nil && (url.Scheme == "" || url.Host == "")
}
