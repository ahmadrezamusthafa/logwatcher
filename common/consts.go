package common

import (
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/common/errors"
)

type ServiceName string

const (
	ASIPCNT        ServiceName = "ASIPCNT"
	ASIPRSV        ServiceName = "ASIPRSV"
	ASIPSRC        ServiceName = "ASIPSRC"
	ASIPSRC_DETAIL ServiceName = "ASIPSRC_DETAIL"
)

type TypeName string

const (
	STANDARD TypeName = "STANDARD"
	DETAIL   TypeName = "DETAIL"
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
	ASIPSRC: `
SELECT 
  timestamp, 
  message, 
  flowid, 
  type, 
  hostname, 
  part, 
  context.correlationid, 
  context.providerid, 
  context.event, 
  context.sourcemarket, 
  context.checkindate, 
  context.checkoutdate, 
  context.locale, 
  context.currency,
  context.noofadult,
  context.noofchild,
  context.noofroom
FROM 
  "s3log"."psrclogrqrs_init"
`,
	ASIPSRC_DETAIL: `
SELECT 
  timestamp, 
  message, 
  flowid, 
  type, 
  hostname, 
  part, 
  context.correlationid
FROM 
  "s3log"."psrclogdetail_init"
`,
}

var MapS3Bucket = map[ServiceName]string{
	ASIPCNT: "asipcnt-logging-775451169198-dc88e4c4e4897c39",
	ASIPSRC: `asipsrc-logging-775451169198-87365897c857f2cd`,
}

func GetBaseQuery(serviceName, typeName string) (string, error) {
	if TypeName(typeName) != STANDARD {
		serviceName += "_" + typeName
	}
	svcName := ServiceName(serviceName)
	if _, ok := mapQuery[svcName]; ok {
		return mapQuery[svcName], nil
	}
	return "", errors.AddTrace(fmt.Errorf("%s service doesn't have %s log", serviceName, typeName))
}

func (c ServiceName) ToString() string {
	return string(c)
}

const DEFAULT_LIMIT = 10
