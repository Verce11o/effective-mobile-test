package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func DecodeCursor(encodedCursor string) (time.Time, int, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, 0, err
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		return time.Time{}, 0, errors.New("invalid cursor")
	}

	res, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, 0, err
	}

	carID, err := strconv.Atoi(arrStr[1])
	if err != nil {
		return time.Time{}, 0, err
	}

	return res, carID, nil
}

func EncodeCursor(t time.Time, id int) string {
	key := fmt.Sprintf("%v,%v", t.Format(time.RFC3339Nano), id)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
