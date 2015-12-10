package temporal

type TimeSlice interface {
     Range() (int, int)
     Stringer()	string
}
