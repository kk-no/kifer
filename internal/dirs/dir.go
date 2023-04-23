package dirs

import "os"

func CreateDirIfNotExist(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path, os.ModePerm); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return path, nil
}
