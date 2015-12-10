package temporal

type TimePie interface {
     Lower() *TimeSlice
     Upper() *TimeSlice
     InnerRange() (int, int)
     OuterRange() (int, int)
     Stringer() string
}
