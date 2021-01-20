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

func (svc *Service) GenerateQuery(ctx context.Context, serviceCode string, query QueryInput) (string, error) {
	baseConditions := []*types.Condition{}
	queryContexts := []*types.Condition{}
	if query.ContextQuery != "" {
		condition, err := multigenerator.GenerateCondition(query.ContextQuery)
		if err != nil {
			return "", errors.AddTrace(err)
		}
		queryContexts = append(queryContexts, &condition)
	}
	if query.MessageQuery != "" {
		queryContexts = append(queryContexts, parseMessageQuery(query.MessageQuery)...)
	}
	baseConditions = append(baseConditions, queryContexts...)
	baseCondition := types.BaseCondition{
		Conditions: []*types.Condition{
			{
				Conditions: baseConditions,
			},
		},
	}
	svcName := ServiceName(serviceCode)
	generatedQuery, err := multigenerator.GenerateQuery(svcName.GetBaseQuery(), baseCondition)
	if err != nil {
		return "", errors.AddTrace(err)
	}
	fmt.Println(generatedQuery)
	return generatedQuery, nil
}

func (svc *Service) Query(ctx context.Context, serviceCode string, query QueryInput) error {
	generatedQuery, err := svc.GenerateQuery(ctx, serviceCode, query)
	if err != nil {
		return errors.AddTrace(err)
	}
	//TODO: exec query to athena
	fmt.Println(generatedQuery)
	return nil
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
			isOpenQuote = !isOpenQuote
		default:
			buffer.WriteRune(char)
		}
		if !isOpenQuote {
			if buffer.Len() > 0 {
				tokenAttributes = appendAttribute(tokenAttributes, buffer, buffer.String())
			}
		}
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
