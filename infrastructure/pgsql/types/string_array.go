package types

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	arr := pq.StringArray{}
	if err := arr.Scan(value); err != nil {
		return err
	}

	*a = StringArray(arr)
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	return pq.StringArray(a).Value()
}

func (a StringArray) GormDataType() string {
	return "text[]"
}

func (a StringArray) String() string {
	return fmt.Sprintf("[%s]", strings.Join(a, ", "))
}
