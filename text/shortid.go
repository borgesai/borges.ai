package text

import (
	"github.com/speps/go-hashids"
	"os"
)

const (
	alphabet string = "abcdefghijklmnopqrstuvwxyz"
)

func NumberToHash(num int64) string {
	hd := hashids.NewData()
	hd.Alphabet = alphabet
	hd.Salt = os.Getenv("SHORTID_SALT")
	h, _ := hashids.NewWithData(hd)
	id, _ := h.EncodeInt64([]int64{num})
	return id
}

func HashToNumber(hash string) int64 {
	hd := hashids.NewData()
	hd.Alphabet = alphabet
	hd.Salt = os.Getenv("SHORTID_SALT")
	h, _ := hashids.NewWithData(hd)
	dec, _ := h.DecodeInt64WithError(hash)
	if len(dec) > 0 {
		return dec[0]
	}
	return 0
}
