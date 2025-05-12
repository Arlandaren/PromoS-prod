package b2b

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("не удалось распарсить дату '%s': %v", s, err)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	s := d.Time.Format("2006-01-02")
	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	var s string
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("не удаётся преобразовать %T в Date", value)
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("не удалось распарсить дату '%s': %v", s, err)
	}
	d.Time = t
	return nil
}

func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}
