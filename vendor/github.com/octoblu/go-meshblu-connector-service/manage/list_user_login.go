package manage

import "io/ioutil"

// ListUserLogin lists all service uuids installed in as user-login
func ListUserLogin(localAppData string) ([]string, error) {
	var empty []string

	files, err := ioutil.ReadDir(userLoginDirectory(localAppData))
	if err != nil {
		return empty, err
	}

	var uuids []string
	for _, file := range files {
		if !file.IsDir() || file.Name() == "bin" {
			continue
		}

		uuids = append(uuids, file.Name())
	}

	return uuids, nil
}
