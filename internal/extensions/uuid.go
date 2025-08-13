package extensions

import (
	"database/sql/driver"
	"strings"

	"github.com/google/uuid"
)

type UUID struct {
	UUID uuid.UUID
}

func (id UUID) String() string {
	return id.UUID.String()
}

func (id *UUID) UnmarshalJSON(bytes []byte) error {
	text := string(bytes)
	text = strings.Replace(text, "\"", "", -1)
	return id.UUID.UnmarshalText([]byte(text))
}

func (id *UUID) UnmarshalBinary(data []byte) error {
	return id.UUID.UnmarshalBinary(data)
}

func (id UUID) Value() (driver.Value, error) {
	return id.UUID.Value()
}

func (id *UUID) MarshalBinary() (data []byte, err error) {
	return id.UUID.MarshalBinary()
}

func (id *UUID) MarshalText() (text []byte, err error) {
	return id.UUID.MarshalText()
}

func (id *UUID) Scan(src any) error {
	return id.UUID.Scan(src)
}

func (id *UUID) UnmarshalParam(param string) error {
	parsedId, err := uuid.Parse(param)
	if err != nil {
		return err
	}
	id.UUID = parsedId
	return nil
}

func (id *UUID) UnmarshalText(text []byte) error {
	return id.UnmarshalParam(string(text))
}

func NewUUID() UUID {
	return UUID{
		UUID: uuid.New(),
	}
}
