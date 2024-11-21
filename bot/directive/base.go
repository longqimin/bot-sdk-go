package directive

import (
	uuid "github.com/satori/go.uuid"
)

type BaseDirective struct {
	Type string `json:"type"`
}

func (this *BaseDirective) GenToken() string {
	var err error
	return uuid.Must(uuid.NewV4(), err).String()
}
