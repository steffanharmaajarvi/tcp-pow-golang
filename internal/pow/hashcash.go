package pow

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"
)

const hashcashDuration = 3600

type Hashcash struct {
	Bits      uint // Number of zero Bits
	Zeros     uint // Number of zero digits
	SaltLen   uint
	Counter   int
	Datetime  int64
	BaseValue string // Base value
	Extra     string // Extension to add to the minted stamp
	Client    string
}

// New creates a new Hash with specified options.
func New(baseValue string, client string, bits uint, saltLen uint, extra string) *Hashcash {

	h := &Hashcash{
		Bits:      bits,
		SaltLen:   saltLen,
		Datetime:  time.Now().Unix(),
		BaseValue: baseValue,
		Client:    client,
		Counter:   0,
		Extra:     extra,
	}
	h.Zeros = uint(math.Ceil(float64(h.Bits) / 4.0))
	return h
}

// NewStd creates a new Hash with 20 Bits of collision and 8 bytes of salt chars.
func NewStd(baseValue string, client string) *Hashcash {
	return New(baseValue, client, 3, 8, "")
}

func hashFromString(str string) []byte {
	h := sha1.New()
	h.Write([]byte(str))

	return h.Sum(nil)
}

func base64EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func base64EncodeInt(n int) string {
	return base64EncodeBytes([]byte(strconv.Itoa(n)))
}

func (h *Hashcash) getHeader() string {
	return fmt.Sprintf("%s:%s:%d:%d:%s:%s", h.Extra,
		h.BaseValue,
		h.Datetime,
		h.Zeros,
		h.Client,
		base64EncodeInt(h.Counter),
	)
}

func (h *Hashcash) Check() bool {
	hashSum := hashFromString(h.getHeader())

	sumUint64 := binary.BigEndian.Uint64(hashSum)
	sumBits := strconv.FormatUint(sumUint64, 2)

	zeroes := 64 - len(sumBits)

	return uint(zeroes) >= h.Bits && h.checkDate()
}

func (h *Hashcash) checkDate() bool {
	return time.Now().Unix()-h.Datetime < hashcashDuration
}

func (h *Hashcash) ComputeHashcash(maxIterations int) (*Hashcash, error) {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		if h.Check() {
			fmt.Println(h)
			return h, nil
		}
		// if hash don't have needed count of leading zeros, we are increasing counter and try next hash
		h.Counter++
	}
	return h, fmt.Errorf("maximum iterations exceeded")
}
