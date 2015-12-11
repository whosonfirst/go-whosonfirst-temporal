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

type Flags interface {
     GetBoolean(string) (bool, error)
     SetBoolean(string, bool) (bool, error)
}

func StringToTime (s string) (time.Time, Flags, error) {

     re, err := regexp.Compile(`(?i)^(\d{1,}-\d{2}-\d{2})(?:\s?(BCE))?$`)

     if err != nil {
        nil_time := time.Time{}
     	return nil_time, nil, err
     }

     m := re.FindStringSubmatch(s)

     if len(m) == 0 {
        nil_time := time.Time{}
     	return nil_time, nil, errors.New("Invalid string")
     }     

     t, err := time.Parse(ISO_8601, m[1])

     if err != nil {
     	return nil_time, nil, err
     }

     flags := NewDefaultTimeFlags()

     if m[2] != "" {
     	flags.SetBoolean("bce", true)
     }

     return t, flags, nil
}

// see below inre notes about flags (and bce and upper)

func TimeToInt(t time.Time, flags Flags) int {

	var i int
	i = ClearTime(i)

	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	bce, _ := flags.GetBoolean("bce")
	upper, _ := flags.GetBoolean("upper")

	if bce == true {
		year = -year
	}

	i = SetYear(i, year)
	i = SetMonth(i, month) // Go is weird...
	i = SetDay(i, day)

	if upper == true {
	   i = (i | UPPER_FLAG)
	}

	return i
}

func IntToTime(i int) (time.Time, Flags) {

	year := i >> 16

	flags := NewDefaultTimeFlags()

	if (i & BCE_FLAG) != 0 {
	   flags.SetBoolean("bce", true)
	}

	if (i & UPPER_FLAG) != 0 {
	   flags.SetBoolean("upper", true)
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

     lower_time, lower_flags, lower_err := StringToTime(dates[0])

     if lower_err != nil {
     	return nil, lower_err
     }

     upper_time, upper_flags, upper_err := StringToTime(dates[0])

     if upper_err != nil {
     	return nil, upper_err
     }

     // Do some sanity checking around dates here and set BCE flags
     // accordingly (20151211/thisisaaronland)

     // Hey look - see what we're doing here? There is no way for the
     // computer (or more specifically the TimeSlice) to "know" it is
     // the upper bounds of a range since TimeSlices don't even know
     // about ranges (20151211/thisisaaronland)
 
     upper_flags.SetBoolean("upper", true)

     lower, err := NewTimeSlice(lower_time, lower_flags)

     if err != nil {
     	return nil, err
     }

     upper, err := NewTimeSlice(upper_time, upper_flags)

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

	return NewTimeSlice(t, flags)
}

func NewTimeSlice(t time.Time, flags Flags) (*TimeSlice, error) {

	ts := TimeSlice{time: t, flags: flags}
	return &ts, nil
}

func NewDefaultTimeFlags() *TimeFlags {

     booleans := make(map[string]bool)
     booleans["bce"] = false
     booleans["upper"] = false

     return NewTimeFlags(booleans)
}

func NewTimeFlags(b map[string]bool) *TimeFlags {
     fl := TimeFlags{booleans: b}
     return &fl
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
	time   time.Time
	flags Flags

}

func (ts *TimeSlice) Flags() Flags {
	return ts.flags
}

func (ts *TimeSlice) AsInt() int {
	return TimeToInt(ts.time, ts.flags)
}

func (ts *TimeSlice) String() string {

	year := ts.time.Year()
	month := int(ts.time.Month())
	day := ts.time.Day()

     	s := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	bce, _ := ts.flags.GetBoolean("bce")

	if bce == true {
	   s = fmt.Sprintf("%s BCE", s)
	}

	return s
}

type TimeFlags struct {
     Flags
     booleans map[string]bool
}

func (tf *TimeFlags) GetBoolean (k string) (bool, error) {

     v, ok := tf.booleans[k]

     if !ok{
     	return false, errors.New("Unknown flag")
     }

     return v, nil
}

func (tf *TimeFlags) SetBoolean (k string, v bool) (bool, error) {

     _, ok := tf.booleans[k]

     if !ok{
     	return false, errors.New("Unknown flag")
     }

     tf.booleans[k] = v
     return true, nil
}
