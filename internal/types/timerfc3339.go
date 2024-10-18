package types

import (
	"fmt"
	"strings"
	"time"
)

// --------------------------------------------------------------------------------
// кастомный тип для времени чтобы настроит маршаллинг в формат RFC3339
type TimeRFC3339 time.Time

const layout = "2006-01-02T15:04:05Z07:00"

func (t TimeRFC3339) String() string {
	tt := time.Time(t)
	return tt.Format(layout)
}

func (t *TimeRFC3339) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"") //убираем кавычки
	if s == "null" {
		return nil
	}
	tt, err := time.Parse(layout, s)
	*t = TimeRFC3339(tt)
	return err
}

func (t TimeRFC3339) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, tt.Format(layout))), nil
}
