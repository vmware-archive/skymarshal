package skycmd

import (
	"encoding/json"
	"errors"

	"github.com/concourse/dex/connector/bitbucket"
	"github.com/hashicorp/go-multierror"
)

func init() {
	RegisterConnector(&Connector{
		id:         "bitbucket",
		config:     &BitbucketFlags{},
		teamConfig: &BitbucketTeamFlags{},
	})
}

type BitbucketFlags struct {
	ClientID     string `long:"client-id" description:"(Required) Client id"`
	ClientSecret string `long:"client-secret" description:"(Required) Client secret"`
}

func (self *BitbucketFlags) Name() string {
	return "Bitbucket"
}

func (self *BitbucketFlags) Validate() error {
	var errs *multierror.Error

	if self.ClientID == "" {
		errs = multierror.Append(errs, errors.New("Missing client-id"))
	}

	if self.ClientSecret == "" {
		errs = multierror.Append(errs, errors.New("Missing client-secret"))
	}

	return errs.ErrorOrNil()
}

func (self *BitbucketFlags) Serialize(redirectURI string) ([]byte, error) {
	if err := self.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(bitbucket.Config{
		ClientID:     self.ClientID,
		ClientSecret: self.ClientSecret,
		RedirectURI:  redirectURI,
	})
}

type BitbucketTeamFlags struct {
	Users []string `long:"user" description:"List of whitelisted Bitbucket users" value-name:"USERNAME"`
	Teams []string `long:"team" description:"List of whitelisted Bitbucket teams" value-name:"ORG_NAME:TEAM_NAME"`
}

func (self *BitbucketTeamFlags) IsValid() bool {
	return len(self.Users) > 0 || len(self.Teams) > 0
}

func (self *BitbucketTeamFlags) GetUsers() []string {
	return self.Users
}

func (self *BitbucketTeamFlags) GetGroups() []string {
	return self.Teams
}
