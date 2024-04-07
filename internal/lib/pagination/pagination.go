package pagination

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func DecodeCursor(encodedCursor string) (int, error) {
	bytes, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return 0, err
	}

	lastID, err := strconv.Atoi(string(bytes))
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func EncodeCursor(id int) string {
	key := fmt.Sprintf("%v", id)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
