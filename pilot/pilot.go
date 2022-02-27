package pilot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Pilot() {
	filepath.Walk("cases", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sol") {
			fmt.Println(path)
		}

		return nil
	})
}
