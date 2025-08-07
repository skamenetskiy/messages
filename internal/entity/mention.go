package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Mention struct {
	ID   uint64
	Type uint32
	Pos  [2]uint32
}

type Mentions []Mention

func (m *Mentions) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("mentions are not []byte")
	}
	return json.Unmarshal(b, &m)
}

func (m Mentions) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Mention) MarshalJSON() ([]byte, error) {
	return json.Marshal([4]uint64{m.ID, uint64(m.Type), uint64(m.Pos[0]), uint64(m.Pos[1])})
}

func (m *Mention) UnmarshalJSON(data []byte) error {
	var v [4]uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	(*m).ID = v[0]
	(*m).Type = uint32(v[1])
	(*m).Pos = [2]uint32{
		uint32(v[2]),
		uint32(v[3]),
	}
	return nil
}
