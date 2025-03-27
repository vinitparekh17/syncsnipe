package utils

import "github.com/vinitparekh17/syncsnipe/internal/colorlog"

// VerifySuccess prints success message if err is nil and returns the error
func VerifySuccess(err error, successMessage string, args ...any) error {
	if err != nil {
		colorlog.CLIError(err.Error())
		return err
	}
	colorlog.CLISuccess(successMessage, args...)
	return nil
}
