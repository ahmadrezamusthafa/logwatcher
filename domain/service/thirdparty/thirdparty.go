package thirdparty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/common"
	"github.com/ahmadrezamusthafa/logwatcher/common/errors"
	"github.com/ahmadrezamusthafa/logwatcher/common/util"
	"github.com/ahmadrezamusthafa/multigenerator"
	"github.com/ahmadrezamusthafa/multigenerator/shared/consts"
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"io/ioutil"
	"reflect"
	"strings"
)

func (svc *Service) GetLogAttributes(ctx context.Context, serviceName, sourceName string) (attributes []LogAttribute, err error) {
	fileName := ""
	switch common.ServiceName(serviceName) {
	case common.ASIPCNT:
		if common.TypeName(sourceName) == common.STANDARD {
			fileName = "asipcnt_logattribute"
		} else {
			return attributes, nil
		}
	case common.ASIPSRC:
		if common.TypeName(sourceName) == common.STANDARD {
			fileName = "asipsrc_logattribute"
		} else {
			fileName = "asipsrc_detail_logattribute"
		}
	default:
		return attributes, nil
	}

	data, err := ioutil.ReadFile("domain/service/thirdparty/resource/" + fileName + ".json")
	if err != nil {
		return attributes, errors.AddTrace(err)
	}
	err = json.Unmarshal(data, &attributes)
	if err != nil {
		return attributes, errors.AddTrace(err)
	}
	return attributes, nil
}

func (svc *Service) GenerateQuery(ctx context.Context, serviceName, typeName string, query QueryInput, limit int) (generated string, err error) {
	baseConditions := []*types.Condition{}
	queryContexts := []*types.Condition{}
	condition := types.Condition{
		Conditions: []*types.Condition{},
	}
	if query.ContextQuery != "" {
		condition, err = multigenerator.GenerateCondition(query.ContextQuery)
		if err != nil {
			return "", errors.AddTrace(err)
		}
	}
	if query.MessageQuery != "" {
		condition.Conditions = append(condition.Conditions, parseMessageQuery(query.MessageQuery)...)
	}
	queryContexts = append(queryContexts, &condition)
	baseConditions = append(baseConditions, queryContexts...)
	if limit <= 0 {
		limit = common.DEFAULT_LIMIT
	}
	baseCondition := types.BaseCondition{
		Conditions: baseConditions,
	}
	svcName := common.ServiceName(serviceName)

	switch svcName {
	case common.ASIPCNT:
		if common.TypeName(typeName) != common.STANDARD {
			return "", errors.AddTrace(fmt.Errorf("%s service doesn't have %s log", serviceName, typeName))
		}
	case common.ASIPSRC:
	default:
		return "", errors.AddTrace(fmt.Errorf("%s log service is not registered yet", serviceName))
	}

	baseQuery, err := common.GetBaseQuery(serviceName, typeName)
	if err != nil {
		return "", errors.AddTrace(err)
	}
	generatedQuery, err := multigenerator.GenerateQuery(baseQuery, baseCondition)
	if err != nil {
		return "", errors.AddTrace(err)
	}
	generatedQuery = strings.TrimSpace(generatedQuery)
	generatedQuery += fmt.Sprint(" LIMIT ", limit)
	return generatedQuery, nil
}

func (svc *Service) Query(ctx context.Context, serviceName, typeName string, query QueryInput, limit int) (outputs []QueryOutput, err error) {
	generatedQuery, err := svc.GenerateQuery(ctx, serviceName, typeName, query, limit)
	if err != nil {
		return outputs, errors.AddTrace(err)
	}
	fmt.Println(generatedQuery)
	svcName := common.ServiceName(serviceName)
	rows, err := svc.DB.GetDB(svcName).Query(generatedQuery)
	if err != nil {
		return outputs, errors.AddTrace(err)
	}
	for rows.Next() {
		row := QueryOutput{}
		var finalContext interface{}

		switch svcName {
		case common.ASIPCNT:
			if common.TypeName(typeName) == common.STANDARD {
				contextInfo := ASIPCNTContext{}
				err = rows.Scan(
					&row.Timestamp,
					&row.Message,
					&row.FlowID,
					&row.Type,
					&row.Hostname,
					&row.Part,
					&contextInfo.CorrelationID,
					&contextInfo.Event,
					&contextInfo.Uri,
					&contextInfo.ProviderID,
					&contextInfo.ProviderHotelID,
					&contextInfo.Locale,
					&contextInfo.ProviderBrandID,
					&contextInfo.ProviderChainID,
				)
				if err != nil {
					return outputs, errors.AddTrace(err)
				}
				finalContext = contextInfo
			} else {
				return outputs, errors.AddTrace(fmt.Errorf("%s service doesn't have %s log", serviceName, typeName))
			}
		case common.ASIPSRC:
			if common.TypeName(typeName) == common.STANDARD {
				contextInfo := ASIPSRCContext{}
				err = rows.Scan(
					&row.Timestamp,
					&row.Message,
					&row.FlowID,
					&row.Type,
					&row.Hostname,
					&row.Part,
					&contextInfo.CorrelationID,
					&contextInfo.ProviderID,
					&contextInfo.Event,
					&contextInfo.SourceMarket,
					&contextInfo.CheckinDate,
					&contextInfo.CheckoutDate,
					&contextInfo.Locale,
					&contextInfo.Currency,
					&contextInfo.NoOfAdult,
					&contextInfo.NoOfChild,
					&contextInfo.NoOfRoom,
				)
				if err != nil {
					return outputs, errors.AddTrace(err)
				}
				finalContext = contextInfo
			} else {
				contextInfo := ASIPSRCDetailContext{}
				err = rows.Scan(
					&row.Timestamp,
					&row.Message,
					&row.FlowID,
					&row.Type,
					&row.Hostname,
					&row.Part,
					&contextInfo.CorrelationID,
				)
				if err != nil {
					return outputs, errors.AddTrace(err)
				}
				finalContext = contextInfo
			}
		default:
			return outputs, errors.AddTrace(fmt.Errorf("%s log service is not registered yet", serviceName))
		}

		contextHtml := generateContextHtml(finalContext)
		row.Context = &contextHtml
		outputs = append(outputs, row)
	}
	return outputs, nil
}

func parseMessageQuery(message string) []*types.Condition {
	attributes := getTokenAttributes(message)
	if attributes != nil {
		conditions := []*types.Condition{}
		for _, value := range attributes {
			conditions = append(conditions, &types.Condition{
				Attribute: &types.Attribute{Name: "message", Operator: consts.OperatorLike, Value: fmt.Sprint("%", value.Value, "%"), Type: valuetype.Alphanumeric},
			})
		}
		return conditions
	}
	return nil
}

func getTokenAttributes(query string) []*types.TokenAttribute {
	var tokenAttributes []*types.TokenAttribute
	buffer := &bytes.Buffer{}
	isQuoteFound := false
	isOpenQuote := false
	for _, char := range query {
		switch char {
		case ' ', '\n', '\'':
			if !isOpenQuote {
				continue
			} else {
				buffer.WriteRune(char)
			}
		case '"':
			isQuoteFound = true
			isOpenQuote = !isOpenQuote
		default:
			buffer.WriteRune(char)
		}
		if !isOpenQuote && isQuoteFound {
			if buffer.Len() > 0 {
				tokenAttributes = appendAttribute(tokenAttributes, buffer, buffer.String())
			}
		}
	}
	if len(tokenAttributes) == 0 {
		tokenAttributes = appendAttribute(tokenAttributes, buffer, query)
	}
	return tokenAttributes
}

func appendAttribute(tokenAttributes []*types.TokenAttribute, buffer *bytes.Buffer, value string) []*types.TokenAttribute {
	tokenAttributes = append(tokenAttributes, &types.TokenAttribute{
		Value: value,
	})
	buffer.Reset()
	return tokenAttributes
}

func generateContextHtml(contexts interface{}) string {
	v := reflect.ValueOf(contexts)
	typeOfS := v.Type()

	buffer := bytes.Buffer{}
	for i := 0; i < v.NumField(); i++ {
		typeField := typeOfS.Field(i)
		field := v.Field(i)
		name := typeField.Name
		value := field.Interface()

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		switch value.(type) {
		case int:
			value = util.InterfaceToInt(value)
		case *int:
			value = util.InterfacePtrToInt(value)
		case string:
			value = util.InterfaceToString(value)
		case *string:
			value = util.InterfacePtrToString(value)
		}

		buffer.WriteString(fmt.Sprintf(`<span style="background-color: #33cc66; color: #fff; display: inline-block; padding: 3px 10px; font-weight: bold; border-radius: 5px; margin-bottom: 5px;">%s : %v</span><br/>`,
			name, value))
	}
	return buffer.String()
}
