package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type PeriodNode struct {
	Duration time.Duration
	Cur      Cursor
}

func (n PeriodNode) Cursor() Cursor {
	return n.Cur
}

func (n PeriodNode) String() string {
	return fmt.Sprint(n.Duration)
}

type periodStr struct{}

var Period = periodStr{}

func (p periodStr) String() string {
	return "<period>"
}

func (p periodStr) Parse(input Cursor) (Node, Cursor, error) {
	pat := regexp.MustCompile(`^(?:\d+[hms])+`)
	indices := pat.FindStringSubmatchIndex(input.String())
	if indices == nil {
		return nil, input, nil
	}
	match := input.String()[indices[0]:indices[1]]
	d, err := time.ParseDuration(match)
	if err != nil {
		return nil, input, fmt.Errorf("error while parsing duration %s: %w", d, err)
	}
	rest := input.Advance(indices[1])
	return PeriodNode{Duration: d, Cur: input}, rest, nil
}

type EpochTimeNode struct {
	// ts is epoch timestamp in seconds.
	ts     float64
	cursor Cursor
}

func (n EpochTimeNode) Cursor() Cursor {
	return n.cursor
}

const isoFormat = "2006-01-02T15:04:05-07:00"

// func (t EpochTimeNode) FormatISO() string {
// 	sec, frac := math.Modf(float64(t))
// 	return time.Unix(int64(sec), int64(1_000_000_000*frac)).UTC().Format(isoFormat)
// }

func (n EpochTimeNode) ToIsoTimeNode() IsoTimeNode {
	t := time.UnixMicro(int64(1_000_000 * float64(n.ts)))
	return IsoTimeNode{Time: t, Cur: n.Cursor()}
}

func (n EpochTimeNode) String() string {
	return fmt.Sprintf("%f", n.ts)
}

var EpochTime = epochTimeStr{}

type epochTimeStr struct{}

func (t epochTimeStr) String() string {
	return "<epoch-time>"
}

func (p epochTimeStr) Parse(input Cursor) (Node, Cursor, error) {
	Logf("EpochTime on: %s$", input)
	pat := regexp.MustCompile(`^\d+(\.\d+)?`)
	indices := pat.FindStringIndex(input.String())
	if indices == nil {
		return nil, input, nil
	}
	match := input.String()[indices[0]:indices[1]]
	t, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return nil, input, fmt.Errorf("error while parsing %s: %w", input, err)
	}
	return EpochTimeNode{ts: t, cursor: input}, input.Advance(indices[1]), nil
}

type IsoTimeNode struct {
	Time time.Time
	Cur  Cursor
}

func (n IsoTimeNode) Cursor() Cursor {
	return n.Cur
}

func (n IsoTimeNode) String() string {
	return n.Time.UTC().Format(isoFormat)
}

func (n IsoTimeNode) ToEpochTimeNode() EpochTimeNode {
	ts := n.Time.Unix()
	return EpochTimeNode{ts: float64(ts), cursor: n.Cursor()}
}

type isoTimeStr struct{}

var IsoTime = isoTimeStr{}

func (t isoTimeStr) String() string {
	return "<iso-time>"
}

func (p isoTimeStr) Parse(input Cursor) (Node, Cursor, error) {
	pat := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}`)
	indices := pat.FindStringIndex(input.String())
	if indices == nil {
		return nil, input, nil
	}
	match := input.String()[indices[0]:indices[1]]
	t, err := time.Parse(isoFormat, match)
	if err != nil {
		return nil, input, fmt.Errorf("error while parsing %s: %w", input, err)
	}
	return IsoTimeNode{Time: t, Cur: input}, input.Advance(indices[1]), nil
}
