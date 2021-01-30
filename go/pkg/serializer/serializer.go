package serializer

import (
	"encoding/base64"
	"fmt"
	"time"
)

// Serialize will return a base64 encoded string with all the data we need.
func Serialize(content string) string {
	return fmt.Sprintf("-----%d\n%s\n", time.Now().Unix(), base64.StdEncoding.EncodeToString([]byte(content)))
}
