package utils

import (
	"strconv"
	"strings"
	"time"
)

type Time time.Time

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func StringToTime(s string) time.Time {
	value := strings.Trim(s, `"`) //get rid of "
	// if value == "" || value == "null" {
	// 	return nil
	// }

	t, _ := time.Parse("2006-01-02", value) //parse time
	// if err != nil {
	// 	return err
	// }
	// *c = radTime(t) //set result using the pointer
	return t
}

func (c *Time) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	*c = Time(t) //set result using the pointer
	return nil
}

func (c Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
	// return []byte(`"` + time.Time(c).Format("02/01/2006") + `"`), nil
}
