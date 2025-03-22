package profile

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func validateProfileName(q *database.Queries, profileName string) error {
	profileName = strings.TrimSpace(profileName)
	profileName = regexp.MustCompile(`\s+`).ReplaceAllString(profileName, " ")

	if profileName == "" {
		return errors.New("empty profile's name is not allowed")
	}

	if len(profileName) < 2 || len(profileName) > 50 {
		return errors.New("profile's name must be between 2 and 50 characters")
	}

	var validProfileName = regexp.MustCompile(`^[a-zA-Z0-9 _-]+$`)
	if !validProfileName.MatchString(profileName) {
		return errors.New("profile's name can only contain letters, numbers, spaces, dashes (-), and underscores (_)")
	}

	count, err := q.IsProfileExists(context.Background(), profileName)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("profile with %s name already exists", profileName)
	}

	return nil
}

func AddProfile(q *database.Queries, profileName string) error {
	if err := validateProfileName(q, profileName); err != nil {
		return err
	}

	_, err := q.CreateProfile(context.Background(), profileName)
	return err
}

func GetProfiles(q *database.Queries) ([]string, error) {
	profiles, err := q.ListProfiles(context.Background())
	if err != nil {
		return nil, err
	}

	var profileNames []string
	for _, profile := range profiles {
		profileNames = append(profileNames, profile.Name)
	}

	return profileNames, nil
}
