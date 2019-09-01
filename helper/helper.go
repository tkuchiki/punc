package helper

import "os"

func IsDisabled() bool {
	return os.Getenv("DISABLE_PUNC") != ""
}
