package core

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/vinitparekh17/syncsnipe/internal/database"
)

// ProfileService interface defines the methods for profile management
type ProfileService interface {
	validateProfileName(ctx context.Context, profileName string) error
	AddProfile(ctx context.Context, profileName string) error
	GetProfiles(ctx context.Context) ([]database.Profile, error)
	UpdateProfile(ctx context.Context, oldName string, newName string) error
	DeleteProfile(ctx context.Context, profileName string) error
}

// Profile implements the ProfileService interface
type Profile struct {
	Queries *database.Queries
}

func NewProfile(queries *database.Queries) *Profile {
	return &Profile{Queries: queries}
}

func (p *Profile) validateProfileName(ctx context.Context, profileName string) error {
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

	count, err := p.Queries.IsProfileExists(ctx, profileName)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("profile with %s name already exists", profileName)
	}

	return nil
}

func (p *Profile) AddProfile(ctx context.Context, profileName string) error {
	if err := p.validateProfileName(ctx, profileName); err != nil {
		return err
	}

	_, err := p.Queries.CreateProfile(ctx, profileName)
	return err
}

func (p *Profile) GetProfiles(ctx context.Context) ([]database.Profile, error) {
	return p.Queries.ListProfiles(ctx)
}

func (p *Profile) UpdateProfile(ctx context.Context, oldName string, newName string) error {
	if err := p.validateProfileName(ctx, newName); err != nil {
		return err
	}

	row, err := p.Queries.UpdateProfileByName(ctx, database.UpdateProfileByNameParams{
		Name:   newName,
		Name_2: oldName,
	})

	if err != nil {
		return err
	}

	if row == 0 {
		return fmt.Errorf("profile with %s name does not exist", oldName)
	}

	return nil
}

func (p *Profile) DeleteProfile(ctx context.Context, profileName string) error {
	row, err := p.Queries.DeleteProfileByName(ctx, profileName)
	if err != nil {
		return err
	}

	if row == 0 {
		return fmt.Errorf("profile with %s name does not exist", profileName)
	}

	return nil
}
