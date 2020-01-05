package model

import (
	"regexp"
)

type formatCheck struct{}

var FormatCheck formatCheck

func (formatCheck) Username(username string) bool {
	re, _ := regexp.Compile("[^'\"\\\\]{5,30}")
	return re.MatchString(username)
}

func (formatCheck) Password(password string) bool {
	re, _ := regexp.Compile("[a-zA-Z0-9+\\-*/=!@#$%^&]{8,30}")
	return re.MatchString(password)
}

func (formatCheck) PhoneNumber(phoneNumber string) bool {
	re, _ := regexp.Compile("(0|86|17951)?(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}")
	return re.MatchString(phoneNumber)
}

func (formatCheck) UserProfile(profile string) bool {
	return len(profile) <= 255
}

func (formatCheck) VideoTitle(title string) bool {
	return len(title) >= 5 && len(title) <= 100
}

func (formatCheck) VideoDesc(desc string) bool {
	return len(desc) <= 255
}
