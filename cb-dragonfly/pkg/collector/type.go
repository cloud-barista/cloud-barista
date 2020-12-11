package collector

type AggregateType string

const (
	MIN                     AggregateType = "min"
	MAX                     AggregateType = "max"
	AVG                     AggregateType = "avg"
	LAST                    AggregateType = "last"
	READ_CONNECTION_TIMEOUT               = 6
	MINIMUM                               = "min"
	MAXIMUM                               = "max"
	AVERAGE                               = "avg"
	LATEST                                = "last"
	DELTOPICS                             = "delTopics/"
)

func (a AggregateType) toString() string {
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
