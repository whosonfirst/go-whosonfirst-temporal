package temporal

import (
	// "errors"
	"fmt"
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
	Stringer() string
}

type Range interface {
	Upper() *Date
	Lower() *Date
	Stringer() string
}

type Date interface {
	IsBCE() bool
	AsInt() int
	Stringer() string
}

func NewTimePie(name string, upper Range, lower Range) *TimePie {

	tp := TimePie{name: name, upper: upper, lower: lower}
	return &tp
}

/*
func NewTimeWedgeFromString(s string) *TimeWedge {

     // This is the complicated string parser...
}
*/

func NewTimeWedge(lower Date, upper Date) *TimeWedge {

	tw := TimeWedge{lower: lower, upper: upper}
	return &tw
}

func NewTimeSliceFromInt(i int) *TimeSlice {

	t := IntToTime(i)
	bce := false // check flag here...

	return NewTimeSlice(t, bce)
}

func NewTimeSliceFromString(s string) *TimeSlice {

     	// please to be parsing me...

	t := time.Now()
	bce := false

	return NewTimeSlice(t, bce)
}

func NewTimeSlice(t time.Time, bce bool) *TimeSlice {

	ts := TimeSlice{t: t, bce: bce}
	return &ts
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

func (tp *TimePie) Stringer() string {
	return tp.Name()
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

func (tw *TimeWedge) Stringer() string {
	return "TIME WEDGE"
}

type TimeSlice struct {
	Date
	t   time.Time
	bce bool
}

func (ts *TimeSlice) IsBCE() bool {
	return ts.bce
}

func (ts *TimeSlice) AsInt() int {
	return TimeToInt(ts.t, ts.bce)
}

func (ts *TimeSlice) Stringer() string {
	return fmt.Sprintf("%v", ts.t)
}

func Parse(lower string, upper string) (int, int) {

	t_lower, _ := time.Parse(ISO_8601, lower)
	t_upper, _ := time.Parse(ISO_8601, upper)

	// TO DO BCE

	l := TimeToInt(t_lower, false)
	u := TimeToInt(t_upper, false)

	u = (u | UPPER_FLAG)

	return l, u
}

func UnParse(lower int, upper int) (string, string) {

	t_lower := IntToTime(lower)
	t_upper := IntToTime(upper)

	return t_lower.Format(ISO_8601), t_upper.Format(ISO_8601)
}

func TimeToInt(t time.Time, bce bool) int {

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

	return i
}

func IntToTime(i int) time.Time {

	year := i >> 16
	month := (i & GET_MONTH) >> 12
	day := (i & GET_DAY) >> 7

	/*
		flag := (i & BCE_FLAG)
		period_flag := (i & PERIOD_EXPR_BIT)
	*/

	// TO DO: BCE

	ymd := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	t, _ := time.Parse(ISO_8601, ymd)

	return t
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
