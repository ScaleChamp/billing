package models

import uuid "github.com/satori/go.uuid"

type Project struct {
	Id          uuid.UUID `json:"id"`
	Usage       float64   `json:"-"`
	Credit      float64   `json:"-"`
	AccountType string    `json:"-"`
}
