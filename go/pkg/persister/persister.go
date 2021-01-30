package persister

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Persister interface has a Write func that takes domain and data and returns an error if writing fails.
type Persister interface {
	Write(string, string) (string, error)
}

// NewFilePersister returns a new file persister that writes data to files into specified directory.
func NewFilePersister(directory string) FilePersister {
	return FilePersister{directory: directory}
}

type FilePersister struct {
	directory string
}

func (f FilePersister) Write(fileName, data string) (string, error) {
	absPath, err := filepath.Abs(fmt.Sprintf("%s%s%s", f.directory, string(os.PathSeparator), fileName))
	if err != nil {
		return "", fmt.Errorf("figuring out absolute path failed: %w", err)
	}

	file, err := os.OpenFile(
		absPath,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0o644,
	)
	if err != nil {
		return "", fmt.Errorf("persister.FilePersister.Write os.OpenFile: %w", err)
	}

	if _, err = file.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("persister.FilePersister.Write file.Write: %w", err)
	}

	_ = file.Close()

	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("could not read file contents: %w", err)
	}

	return string(content), nil
}
