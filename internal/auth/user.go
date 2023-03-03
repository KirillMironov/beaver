package auth

import (
	"bytes"
	"encoding/gob"
)

type User struct {
	ID         string
	Passphrase string
}

func (u *User) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	if err = enc.Encode(u); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (u *User) UnmarshalBinary(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	return dec.Decode(u)
}
