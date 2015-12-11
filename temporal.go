package temporal

import (
       "errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (

	// Go is weird and this does not support years < 1000 (and presumably > 9999...)
	// See also: https://github.com/metakeule/fmtdate

	ISO_8601 string = "2006-01-02"

	SIS_DATE    int = 1
	YEAR        int = 2
	DECADE      int = 3
	IMPL_DECADE int = 4
	Century     int = 5
	PERIOD_EXPR int = 6
	IMPL_PERIOD int = 7
	Circa       int = 8

	UNDEFINED int = -1
	UNUSED    int = -1
	T_lower   int = 1
	T_upper   int = 2
	bce       int = 1
	ce        int = 2

	NEGATIVE_INF int = 0x80000000
	POSITIVE_INF int = 0x7fffffff

	RESET_TIME  int = 0x00000000
	RESET_YEAR  int = 0x0000ffff
	RESET_MONTH int = 0xffff0fff
	RESET_DAY   int = 0xfffff07f
	BCE_FLAG    int = 0x80000000
	UPPER_FLAG  int = 0x00000040
	CLEAR_DATE  int = 0x8000007f
	CLEAR_FLAGS int = 0xffffffa0
	GET_YEAR    int = 0xffff0000
	GET_MONTH   int = 0x0000f000
	GET_DAY     int = 0x00000f80

	EXPR_BITS       int = 0x0000003e
	PERIOD_EXPR_BIT int = 0x00000001
	IMPL_PERIOD_BIT int = 0x00000003
	MODE_BIT        int = 0x00000002
	EXPL_DEC_BIT    int = 0x00000004
	IMPL_DEC_BIT    int = 0x0000000c
	DATE_BIT        int = 0x00000008
	CENTURY_BIT     int = 0x00000010
	CIRCA_BIT       int = 0x00000020
	AAT_BIT         int = 0x00000018
)

type Period interface {
	Name() string
	InnerRange() (Date, Date)
	OuterRange() (Date, Date)
	String() string
}

type Range interface {
	Upper() *Date
	Lower() *Date
	String() string
}

type Date interface {
	IsBCE() bool
	IsUpper() bool
	AsInt() int
	String() string
}

// see below inre notes about flags (and bce and upper)

func TimeToInt(t time.Time, bce bool, upper bool) int {

	var i int
	i = ClearTime(i)

	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	if bce {
		year = -year
	}

	i = SetYear(i, year)
	i = SetMonth(i, month) // Go is weird...
	i = SetDay(i, day)

	if upper {
	   i = (i | UPPER_FLAG)
	}

	return i
}

func IntToTime(i int) (time.Time, map[string]bool) {

	year := i >> 16

	// Hey look - soon we will make this (all the flags stuff)
	// into a proper Flag interface and pass that around instead
	// of the kludge we are using now (20151210/thisisaaronland)

	flags := make(map[string]bool)
	flags["is_bce"] = false
	flags["is_upper"] = false

	if (i & BCE_FLAG) != 0 {
	   flags[ "is_bce" ] = true
	}

	if (i & UPPER_FLAG) != 0 {
	   flags[ "is_upper" ] = true
	}

	month := (i & GET_MONTH) >> 12
	day := (i & GET_DAY) >> 7

	ymd := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	// fmt.Printf("YMD %s\n", ymd)

	t, _ := time.Parse(ISO_8601, ymd)
	return t, flags
}

func ClearTime(time int) int {
	return time & RESET_TIME
}

func SetDay(time int, day int) int {
	return (time & RESET_DAY) | ((day | 0) << 7)
}

func SetMonth(time int, month int) int {
	return (time & RESET_MONTH) | ((month | 0) << 12)
}

func SetYear(time int, year int) int {
	return (time & RESET_YEAR) | ((year | 0) << 16)
}

/*
func SetExprBits(time int, expr int) (int, error) {

	var operand int

	switch expr {
	case SIS_DATE:
		operand = DATE_BIT
	case DECADE:
		operand = EXPL_DEC_BIT
	case IMPL_DECADE:
		operand = IMPL_DEC_BIT
	case Century:
		operand = CENTURY_BIT
	case Circa:
		operand = CIRCA_BIT
	case PERIOD_EXPR:
		operand = PERIOD_EXPR_BIT
	case IMPL_PERIOD:
		operand = IMPL_PERIOD_BIT
	default:
		return 0, errors.New("Unknown expression")
	}

	result := (time | operand)
	return result, nil
}
*/

func NewTimePie(name string, upper Range, lower Range) (*TimePie, error) {

	tp := TimePie{name: name, upper: upper, lower: lower}
	return &tp, nil
}

func NewTimeWedgeFromString(s string) (*TimeWedge, error) {

     // This is the complicated string parser but for now it is not complicated
     // and just assumes a pair of comma separated YYYY-MM-DD BCE? strings
     
     // As in: Please to write an EDTF -> YMD parser that we can use here...

     dates := strings.Split(s, ",")

     if len(dates) != 2 {
     	return nil, errors.New("Invalid string")
     }

     lower, err := NewTimeSliceFromString(dates[0], false)

     if err != nil {
     	return nil, err
     }

     upper, err := NewTimeSliceFromString(dates[1], true)

     if err != nil {
     	return nil, err
     }

     return NewTimeWedge(lower, upper)
}

func NewTimeWedge(lower Date, upper Date) (*TimeWedge, error) {

	tw := TimeWedge{lower: lower, upper: upper}
	return &tw, nil
}

func NewTimeSliceFromInt(i int) (*TimeSlice, error) {

	t, flags := IntToTime(i)

	return NewTimeSlice(t, flags["is_bce"], flags["is_upper"])
}

func NewTimeSliceFromString(s string, upper bool) (*TimeSlice, error) {

     re, err := regexp.Compile(`(?i)^(\d{1,}-\d{2}-\d{2})(?:\s?(BCE))?$`)

     if err != nil {
     	return nil, err
     }

     m := re.FindStringSubmatch(s)

     if len(m) == 0 {
     	return nil, errors.New("Invalid string")
     }     

     t, err := time.Parse(ISO_8601, m[1])

     if err != nil {
     	return nil, err
     }

     bce := false

     if m[2] != "" {
     	bce = true
     }

     return NewTimeSlice(t, bce, upper)
}

func NewTimeSlice(t time.Time, bce bool, upper bool) (*TimeSlice, error) {

	ts := TimeSlice{t: t, bce: bce, upper: upper}
	return &ts, nil
}

type TimePie struct {
	Period
	name  string
	upper Range
	lower Range
}

func (tp *TimePie) Name() string {

	return tp.name
}

func (tp *TimePie) InnerRange() (*Date, *Date) {

	return tp.lower.Upper(), tp.upper.Lower()
}

func (tp *TimePie) OuterRange() (*Date, *Date) {

	return tp.lower.Lower(), tp.upper.Upper()
}

func (tp *TimePie) String() string {
	return fmt.Sprintf("%s (%v - %v)", tp.name, tp.lower, tp.upper)
}

type TimeWedge struct {
	Range
	upper Date
	lower Date
}

func (tw *TimeWedge) Upper() Date {
	return tw.upper
}

func (tw *TimeWedge) Lower() Date {
	return tw.lower
}

func (tw *TimeWedge) String() string {

     	return fmt.Sprintf("%v,%v", tw.lower, tw.upper)
}

type TimeSlice struct {
	Date
	t   time.Time
	bce bool
	upper bool
}

func (ts *TimeSlice) IsBCE() bool {
	return ts.bce
}

func (ts *TimeSlice) IsUpper() bool {
	return ts.upper
}

func (ts *TimeSlice) AsInt() int {
	return TimeToInt(ts.t, ts.bce, ts.upper)
}

func (ts *TimeSlice) String() string {

	year := ts.t.Year()
	month := int(ts.t.Month())
	day := ts.t.Day()

     	s := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	if (ts.bce) {
	   s = fmt.Sprintf("%s BCE", s)
	}

	return s
}
