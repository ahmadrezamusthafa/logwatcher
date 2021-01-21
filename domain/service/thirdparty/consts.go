package thirdparty

type ServiceName string

const (
	ASIPCNT ServiceName = "ASIPCNT"
	ASIPRSV ServiceName = "ASIPRSV"
	ASIPSRC ServiceName = "ASIPSRC"
)

var mapQuery = map[ServiceName]string{
	ASIPCNT: `SELECT timestamp,message,flowid,type,hostname,part FROM "s3log"."pcntapirqrs_init"`,
	ASIPRSV: ``,
	ASIPSRC: ``,
}

func (c ServiceName) GetBaseQuery() string {
	if _, ok := mapQuery[c]; ok {
		return mapQuery[c]
	}
	return ""
}

func (c ServiceName) ToString() string {
	return string(c)
}

