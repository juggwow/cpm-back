package rad

import (
	"strings"
	"time"
)

type Time time.Time

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
