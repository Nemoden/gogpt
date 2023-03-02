package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/nemoden/gogpt/config"
)

func Hangup(err error) {
	if errors.Is(err, config.ErrApiTokenFileDoesntExist) {
		fmt.Printf(config.TokenFileDoesntExistPrompt)
	} else if errors.Is(err, config.ErrCantOpenApiTokenFileForReading) {
		fmt.Printf(config.TokenFileNotReadablePrompt)
	} else {
		fmt.Printf("Error: %s", err)
	}
	os.Exit(0)
}
