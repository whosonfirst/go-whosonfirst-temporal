# go-whosonfirst-temporal

A Go package for converting `year-month-day` expressions to and from 8-byte integers following the model of the CIDOC‐CRM Temporal representation specification.

## Caveats

This package is under active development. It is incomplete and probably still has bugs. Among other things:

* It does not handle BCE yet
* It does not handle period "expressions" yet
* It does not handle years before 1000 yet _because Go's date parser appears to be built on top of MADNESS_
* It does not implement temporal operators yet

It does not implement complete CIDOC-CRM (textual) temporal expressions nor will it. That will be left to another package to implement string to year-month-day conversions.

## Example

```
import(
	"fmt"
	"temporal"
)

func main() {

     lower := "1914-08-04"
     upper := "1918-11-11"

     fmt.Println(lower, upper)

     x, y := temporal.Parse(lower, upper)
     fmt.Println(x, y)

     lower, upper = temporal.UnParse(x, y)
     fmt.Println(lower, upper)
}

// This yields:
// 1914-08-04 1918-11-11
// 125469184 125744576
// 1914-08-04 1918-11-11
```

## See also

* http://www.cidoc-crm.org/downloads/CIDOC-CRM_temporal_representation.pdf
