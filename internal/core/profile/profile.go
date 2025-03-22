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

func GetProfiles(q *database.Queries) ([]database.Profile, error) {
	return q.ListProfiles(context.Background())
}

func UpdateProfile(q *database.Queries, oldName string, newName string) error {
	if err := validateProfileName(q, newName); err != nil {
		return err
	}

	return q.UpdateProfileByName(context.Background(), database.UpdateProfileByNameParams{
		Name:   newName,
		Name_2: oldName,
	})
}

func DeleteProfile(q *database.Queries, profileName string) error {
	return q.DeleteProfileByName(context.Background(), profileName)
}
