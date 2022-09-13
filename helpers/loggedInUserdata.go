package helpers

import (
	"io/ioutil"
	"os"
)

var FilePath = "/home/vikash/go_registerdApp/"
var LoggedInUsersfile = "loggedInUsersfile"

// var LoggedInUsersfilePath = FilePath + LoggedInUsersfile + "." + "json"

func LoggedInUserdata() ([]byte, error) {
	loggedInUsersfilePath := FilePath + LoggedInUsersfile + "." + "json"
	if _, err := os.Stat(loggedInUsersfilePath); err != nil {
		os.Create(loggedInUsersfilePath)
	}
	read_data, err := ioutil.ReadFile(loggedInUsersfilePath)
	if err != nil {
		return nil, err
	}
	return read_data, nil
}
