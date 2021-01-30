package persister

import (
	"fmt"
	"log"
	"os"
)

// Persister interface has a Write func that takes domain and data and returns an error if writing fails.
type Persister interface {
	Write(string, string) error
}

// NewFilePersister returns a new file persister that writes data to files into specified directory.
func NewFilePersister(directory string) FilePersister {
	return FilePersister{directory: directory}
}

type FilePersister struct {
	directory string
}

func (f FilePersister) Write(fileName, data string) error {
	file, err := os.OpenFile(
		fmt.Sprintf("%s%s%s", f.directory, string(os.PathSeparator), fileName),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0o644,
	)
	if err != nil {
		return fmt.Errorf("persister.FilePersister.Write os.OpenFile: %w", err)
	}

	defer func() {
		if fErr := file.Close(); fErr != nil {
			log.Fatalf("persister.FilePersister.Write file.Close: %s", fErr)
		}
	}()

	if _, err = file.Write([]byte(data)); err != nil {
		return fmt.Errorf("persister.FilePersister.Write file.Write: %w", err)
	}

	return nil
}
