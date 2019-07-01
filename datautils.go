package gemini

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"golang.org/x/exp/rand"
)

func randIntRange(rnd *rand.Rand, min int, max int) int {
	if max <= min {
		return min
	}
	return rnd.Intn(max-min) + min
}

func randInt64Range(rnd *rand.Rand, min int64, max int64) int64 {
	if max <= min {
		return min
	}
	return rnd.Int63n(max-min) + min
}

func randFloat32Range(rnd *rand.Rand, min float32, max float32) float32 {
	if max <= min {
		return min
	}
	return rnd.Float32() * (max - min)
}

func randFloat64Range(rnd *rand.Rand, min float64, max float64) float64 {
	if max <= min {
		return min
	}
	return rnd.Float64() * (max - min)
}

func randBlobWithTime(rnd *rand.Rand, len int, t time.Time) []byte {
	id, _ := ksuid.NewRandomWithTime(t)

	var buf bytes.Buffer
	buf.Write(id.Bytes())

	if buf.Len() >= len {
		return buf.Bytes()[:len]
	}

	// Pad some extra random data
	buff := make([]byte, len-buf.Len())
	rnd.Read(buff)
	buf.WriteString(base64.StdEncoding.EncodeToString(buff))

	return buf.Bytes()[:len]

}

func randStringWithTime(rnd *rand.Rand, len int, t time.Time) string {
	id, _ := ksuid.NewRandomWithTime(t)

	var buf strings.Builder
	buf.WriteString(id.String())
	if buf.Len() >= len {
		return buf.String()[:len]
	}

	// Pad some extra random data
	buff := make([]byte, len-buf.Len())
	rnd.Read(buff)
	buf.WriteString(base64.StdEncoding.EncodeToString(buff))

	return buf.String()[:len]
}

func nonEmptyRandBlobWithTime(rnd *rand.Rand, len int, t time.Time) []byte {
	if len <= 0 {
		len = 1
	}
	return randBlobWithTime(rnd, len, t)
}

func nonEmptyRandStringWithTime(rnd *rand.Rand, len int, t time.Time) string {
	if len <= 0 {
		len = 1
	}
	return randStringWithTime(rnd, len, t)
}

func randDate(rnd *rand.Rand) string {
	time := randTime(rnd)
	return time.Format("2006-01-02")
}

func randTime(rnd *rand.Rand) time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()

	sec := rnd.Int63n(max-min) + min
	return time.Unix(sec, 0)
}

func randTimeNewer(rnd *rand.Rand, d time.Time) time.Time {
	min := time.Date(d.Year()+1, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()

	sec := rnd.Int63n(max-min+1) + min
	return time.Unix(sec, 0)
}

func randIpV4Address(rnd *rand.Rand, v, pos int) string {
	if pos < 0 || pos > 4 {
		panic(fmt.Sprintf("invalid position for the desired value of the IP part %d, 0-3 supported", pos))
	}
	if v < 0 || v > 255 {
		panic(fmt.Sprintf("invalid value for the desired position %d of the IP, 0-255 suppoerted", v))
	}
	var blocks []string
	for i := 0; i < 4; i++ {
		if i == pos {
			blocks = append(blocks, strconv.Itoa(v))
		} else {
			blocks = append(blocks, strconv.Itoa(rnd.Intn(255)))
		}
	}
	return strings.Join(blocks, ".")
}

func appendValue(columnType Type, r *rand.Rand, p PartitionRangeConfig, values []interface{}) []interface{} {
	return append(values, columnType.GenValue(r, p)...)
}

/*
func appendValueRange(columnType Type, r *rand.Rand, p PartitionRangeConfig, values []interface{}) []interface{} {
	values = append(values, columnType.GenValue(r, p)...)
	values = append(values, columnType.GenValue(r, p)...)
	return values
}
*/
