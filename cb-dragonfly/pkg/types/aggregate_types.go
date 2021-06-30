package types

type AggregateType string

const (
	MIN                   AggregateType = "min"
	MAX                   AggregateType = "max"
	AVG                   AggregateType = "avg"
	LAST                  AggregateType = "last"
	ReadConnectionTimeout               = 6
	DELTOPICS                           = "delTopics/"
)

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
