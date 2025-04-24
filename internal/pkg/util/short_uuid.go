package util

import (
	"github.com/btcsuite/btcutil/base58"

	uid "github.com/satori/go.uuid"
)

func GetId() string {
	id := uid.NewV4()
	return id.String()
}

func Encode(id string) string {
	if id == "00000000-0000-0000-0000-000000000000" || id == "" {
		return ""
	}
	uuid, _ := uid.FromString(id)
	return base58.Encode(uuid.Bytes())
}

func Decode(id string) (*uid.UUID, error) {
	decID, err := uid.FromBytes(base58.Decode(id))
	if err != nil {
		return nil, err
	}
	return &decID, nil
}

func DecodeGetUUID(id string) (uid.UUID, error) {
	idDecode, err := Decode(id)
	if err != nil {
		return uid.UUID{}, err
	}
	idUUID, _ := uid.FromString(idDecode.String())
	return idUUID, nil
}

func ParseUUIDString(id string) string {
	uuid, err := uid.FromBytes(base58.Decode(id))
	if err != nil {
		return ""
	}

	return uuid.String()
}
