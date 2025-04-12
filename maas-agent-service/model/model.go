package model

import (
	"context"
	"encoding/base64"
)

type AuthCredentials struct {
	Username string
	Password string
}

func (a AuthCredentials) AuthHeaderProvider(_ context.Context) (string, error) {
	encodedCreds := base64.StdEncoding.EncodeToString([]byte(a.Username + ":" + a.Password))
	return "Basic " + encodedCreds, nil
}

type CpWatcherMessageDto struct {
	State   CpVersionsDto `json:"state,omitempty"`
	Changes []CpChange    `json:"changes,omitempty"`
}

type CpVersionsDto []CpDeploymentVersion

type CpChange struct {
	New *CpDeploymentVersion `json:"new"`
	Old *CpDeploymentVersion `json:"old"`
}

type CpDeploymentVersion struct {
	Version     string `json:"version"`
	Stage       string `json:"stage"`
	CreatedWhen string `json:"createdWhen"`
	UpdatedWhen string `json:"updatedWhen"`
}
