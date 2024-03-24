package models

import "encoding/json"

type AgrRepo struct {
	FullName   *string         `json:"full_name,omitempty"`
	Owner      *string         `json:"owner,omitempty"`
	Repository *string         `json:"repository,omitempty"`
	Languages  json.RawMessage `json:"languages,omitempty"`
}
