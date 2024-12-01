package types

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
    str := string(b)
    str = str[1 : len(str)-1]

    t, err := time.Parse("02.01.2006", str)
    if err != nil {
        return err
    }
    *d = Date(t)
    return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Time(d).Format("02.01.2006"))
}

func (d *Date) Scan(value interface{}) error {
    switch v := value.(type) {
    case time.Time:
        *d = Date(v)
        return nil
    case string:
        t, err := time.Parse("02.01.2006", v)
        if err != nil {
            return err
        }
        *d = Date(t)
        return nil
    default:
        return fmt.Errorf("cannot scan type %T into Date", value)
    }
}

func (d Date) Value() (driver.Value, error) {
    t := time.Time(d)
    return t, nil
}

func (d Date) String() string {
    return time.Time(d).Format("02.01.2006")
}
