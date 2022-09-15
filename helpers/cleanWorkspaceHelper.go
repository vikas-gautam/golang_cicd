package helpers

import (
	"fmt"
	"os"
)

// To remove older workspace
func CleanWorkspace(DestFolder string) error {
	if err := os.RemoveAll(DestFolder); err != nil {
		fmt.Println("not able to clean ws")
		return err
	}
	return nil
}
