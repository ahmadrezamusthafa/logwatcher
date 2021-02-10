package shared

import (
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/errors"
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
