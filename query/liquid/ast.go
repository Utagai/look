package liquid

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/utagai/look/datum"
	"github.com/utagai/look/query/liquid/execution"
)

type ConstKind string

const (
	ConstKindString ConstKind = "string"
	ConstKindNumber ConstKind = "number"
	ConstKindBool   ConstKind = "bool"
)

type Const struct {
	Kind        ConstKind
	Stringified string
}

func (c *Const) String() string {
	return c.Stringified
}

func (c *Const) Interface() interface{} {
	switch c.Kind {
	case ConstKindString:
		return c.Stringified
	case ConstKindNumber:
		f64, _ := strconv.ParseFloat(c.Stringified, 64)
		return f64
	case ConstKindBool:
		return c.Stringified == "true"
	default:
		panic(fmt.Sprintf("unexpected const kind: %q", c.Kind))
	}
}

type BinaryOp string

const (
	BinaryOpEquals   BinaryOp = "="
	BinaryOpGeq               = ">"
	BinaryOpContains          = "contains"
)

func evaluateBinaryOp(left interface{}, right *Const, op BinaryOp) bool {
	switch op {
	case BinaryOpEquals:
		return execution.Compare(left, right.Interface()) == execution.Equal
	case BinaryOpGeq:
		return execution.Compare(left, right.Interface()) == execution.Greater
	case BinaryOpContains:
		return strings.Contains(left.(string), right.Stringified)
	default:
		panic(fmt.Sprintf("unrecognized BinaryOp: %q", op))
	}
}

type Check struct {
	Field string
	Value *Const
	Op    BinaryOp
}

func (c *Check) Evaluate(datum datum.Datum) bool {
	fieldValue, ok := datum[c.Field]
	if !ok {
		// If the field does not exist on this datum, evaluate to false.
		return false
	}

	return evaluateBinaryOp(fieldValue, c.Value, c.Op)
}

func (c *Check) String() string {
	return fmt.Sprintf("%s %s %v", c.Field, c.Op, c.Value)
}

type Find struct {
	Checks []*Check
}

type FindStream struct {
	*Find
	source datum.DatumStream
}

// Next implements the datum.DatumStream interface.
func (fs *FindStream) Next() (datum.Datum, error) {
	// Keep iterating the stream until something passes the checks.
outer:
	for {
		datum, err := fs.source.Next()
		if err != nil {
			return nil, err
		}

		for _, check := range fs.Find.Checks {
			// If we failed, move onto the next datum.
			if !check.Evaluate(datum) {
				continue outer
			}
		}

		// If we get here, we have successfully evaluated against every check, and
		// we can be returned.
		return datum, nil
	}
}

// Execute implements the Stage interface.
func (f *Find) Execute(datums datum.DatumStream) (datum.DatumStream, error) {
	return &FindStream{
		Find:   f,
		source: datums,
	}, nil
}

func (f *Find) Name() string {
	return "find"
}

type SortStream struct {
	*Sort
	sortedDatums []datum.Datum
	source       datum.DatumStream
	sortedSource datum.DatumStream
}

func (ss *SortStream) Next() (datum.Datum, error) {
	if ss.sortedDatums == nil {
		if err := ss.sortStream(); err != nil {
			return nil, fmt.Errorf("failed to sort data: %w", err)
		}
	}
	// TODO: Sort this...
	return ss.sortedSource.Next()
}

func (ss *SortStream) sortStream() error {
	// TODO: Dangerous if source is large.
	datums, err := datum.StreamToSlice(ss.source)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}
	sortable := sortableDatums{datums: datums, fieldName: ss.Field}
	sort.Sort(sortable)
	ss.sortedDatums = sortable.datums
	ss.sortedSource = datum.NewDatumSliceStream(sortable.datums)

	return nil
}

type sortableDatums struct {
	datums    []datum.Datum
	fieldName string
}

func (ds sortableDatums) Len() int {
	return len(ds.datums)
}

func (ds sortableDatums) Less(i, j int) bool {
	ithDoc := ds.datums[i]
	jthDoc := ds.datums[j]

	ithValue, ok := ithDoc[ds.fieldName]
	if !ok {
		// Treat documents where the field does not exist as being less than.
		return true
	}

	jthValue, ok := jthDoc[ds.fieldName]
	if !ok {
		// Ditto above.
		return true
	}

	return execution.Compare(ithValue, jthValue) == execution.Lesser
}

func (ds sortableDatums) Swap(i, j int) {
	ds.datums[i], ds.datums[j] = ds.datums[j], ds.datums[i]
}

type Sort struct {
	Descending bool
	Field      string
}

func (s *Sort) Execute(datums datum.DatumStream) (datum.DatumStream, error) {
	return &SortStream{
		Sort:   s,
		source: datums,
	}, nil
}

func (s *Sort) Name() string {
	return "sort"
}

type Stage interface {
	// TODO: This should be decoupled from the AST.
	Execute(datum.DatumStream) (datum.DatumStream, error)
	Name() string
}
