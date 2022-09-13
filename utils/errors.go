package utils

import "os"

func LogError(err error, errorFile *os.File) {
	_, err = errorFile.WriteString(err.Error() + "\n")
}
