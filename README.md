# go-whosonfirst-temporal

A Go package for converting `year-month-day` expressions to and from 8-byte integers modeled after the model of the CIDOC‐CRM Temporal representation specification.

## Caveats

This package is under active development. It is incomplete and probably still has bugs. Among other things:

* It is not the CIDOC‐CRM Temporal specification – specifically the numeric representations of dates are _not_ the same
* It does not handle period "expressions" or any of signifiers yet
* It does not handle years greater than 9999 (or less than -9999) yet _because Go's date parser appears to be built on top of MADNESS..._
* It does not implement temporal operators yet

It does not implement complete CIDOC-CRM (textual) temporal expressions nor will it. Currently there are a handful of `NewThingFromString` functions which allow for simple `YYYY-MM-DD BCE?` strings but these will probably be replaced or at least superseded by equivalent functions that hide more complex string parsing from this package. We'll see.

## Example

```
import (
       "fmt"
       "github.com/whosonfirst/go-whosonfirst-temporal"
       )

func main (){

     lower := "1914-08-04"
     upper := "1918-11-11"

     s := fmt.Sprintf("%s,%s", lower, upper)
     w, _ := temporal.NewTimeWedgeFromString(s)

     fmt.Printf("wedge: %v\n", w)
     fmt.Printf("lower: %v\n", w.Lower())
     fmt.Printf("upper: %v\n", w.Upper())

     lower_int := w.Lower().AsInt()
     upper_int := w.Upper().AsInt()

     fmt.Printf("lower (as int): %d\n", lower_int)
     fmt.Printf("upper (as int): %d\n", upper_int)

     lower_slice, _ := temporal.NewTimeSliceFromInt(lower_int)
     upper_slice, _ := temporal.NewTimeSliceFromInt(upper_int)

     fmt.Printf("lower (from int): %v (%d)\n", lower_slice, lower_slice.AsInt())
     fmt.Printf("upper (from int): %v (%d)\n", upper_slice, upper_slice.AsInt())
}

// The output of this would be:
// wedge: 1914-08-04,1918-11-11
// lower: 1914-08-04
// upper: 1918-11-11
// lower (as int): 125469184
// upper (as int): 125744576
// lower (from int): 1914-08-04 (125469184)
// upper (from int): 1918-11-11 (125744576)
```

_Not discussed here is the `TimePie` which is a named pair of `TimeWedge` object-thingies._

## See also

* http://cidoc-crm.org/tools.html
* http://www.cidoc-crm.org/downloads/CIDOC-CRM_temporal_representation.pdf
