package thirdparty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahmadrezamusthafa/logwatcher/common/errors"
	"github.com/ahmadrezamusthafa/multigenerator"
	"github.com/ahmadrezamusthafa/multigenerator/shared/consts"
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"io/ioutil"
	"strings"
)

func (svc *Service) GetLogAttributes(ctx context.Context) (attributes []LogAttribute, err error) {
	data, err := ioutil.ReadFile("domain/service/thirdparty/resource/asipcnt_logattribute.json")
	if err != nil {
		return attributes, errors.AddTrace(err)
	}
	err = json.Unmarshal(data, &attributes)
	if err != nil {
		return attributes, errors.AddTrace(err)
	}
	return attributes, nil
}

func (svc *Service) GenerateQuery(ctx context.Context, serviceCode string, query QueryInput, limit int) (generated string, err error) {
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
		limit = DEFAULT_LIMIT
	}
	baseCondition := types.BaseCondition{
		Conditions: baseConditions,
	}
	svcName := ServiceName(serviceCode)
	generatedQuery, err := multigenerator.GenerateQuery(svcName.GetBaseQuery(), baseCondition)
	if err != nil {
		return "", errors.AddTrace(err)
	}
	generatedQuery = strings.TrimSpace(generatedQuery)
	generatedQuery += fmt.Sprint(" LIMIT ", limit)
	return generatedQuery, nil
}

func (svc *Service) Query(ctx context.Context, serviceCode string, query QueryInput, limit int) (outputs []QueryOutput, err error) {
	generatedQuery, err := svc.GenerateQuery(ctx, serviceCode, query, limit)
	if err != nil {
		return outputs, errors.AddTrace(err)
	}
	fmt.Println(generatedQuery)
	rows, err := svc.DB.GetDB().Query(generatedQuery)
	if err != nil {
		return outputs, errors.AddTrace(err)
	}
	for rows.Next() {
		row := QueryOutput{}
		err = rows.Scan(&row.Timestamp, &row.Message, &row.FlowID, &row.Type, &row.Hostname, &row.Part)
		if err != nil {
			return outputs, errors.AddTrace(err)
		}
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
