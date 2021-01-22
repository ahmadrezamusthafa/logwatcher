package common

type ServiceName string

const (
	ASIPCNT ServiceName = "ASIPCNT"
	ASIPRSV ServiceName = "ASIPRSV"
	ASIPSRC ServiceName = "ASIPSRC"
)

var mapQuery = map[ServiceName]string{
	ASIPCNT: `
SELECT 
  timestamp, 
  message, 
  flowid, 
  type, 
  hostname, 
  part, 
  context.correlationid, 
  context.event, 
  context.uri, 
  context.providerid, 
  context.providerhotelid, 
  context.locale, 
  context.providerbrandid, 
  context.providerchainid 
FROM 
  "s3log"."pcntapirqrs_init"
`,
	ASIPRSV: ``,
	ASIPSRC: ``,
}

var mapS3Bucket = map[ServiceName]string{
	ASIPCNT: "asipcnt-logging-775451169198-dc88e4c4e4897c39",
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

const DEFAULT_LIMIT = 10
