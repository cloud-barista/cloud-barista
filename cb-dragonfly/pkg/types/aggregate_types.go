package types

type AggregateType string

const (
	MIN  AggregateType = "min"
	MAX  AggregateType = "max"
	AVG  AggregateType = "avg"
	LAST AggregateType = "last"
)

const ReadConnectionTimeout = 5

func (a AggregateType) ToString() string {
	switch a {
	case MIN:
		return "min"
	case MAX:
		return "max"
	case AVG:
		return "avg"
	case LAST:
		return "last"
	default:
		return ""
	}
}
