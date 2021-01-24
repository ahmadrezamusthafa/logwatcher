package querygen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ahmadrezamusthafa/multigenerator/shared/consts"
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"strings"
)

type QueryGen struct {
}

var (
	operatorMap = map[string]interface{}{
		consts.OperatorEqual:            nil,
		consts.OperatorNotEqual:         nil,
		consts.OperatorLessThan:         nil,
		consts.OperatorLessThanEqual:    nil,
		consts.OperatorGreaterThan:      nil,
		consts.OperatorGreaterThanEqual: nil,
		consts.OperatorInclude:          nil,
		consts.OperatorExclude:          nil,
		consts.OperatorIsNull:           nil,
		consts.OperatorIsNotNull:        nil,
		consts.OperatorLike:             nil,
	}

	logicalOperatorMap = map[string]interface{}{
		consts.LogicalOperatorAnd: nil,
		consts.LogicalOperatorOr:  nil,
	}
)

func (g *QueryGen) GenerateQuery(mainQuery string, baseCondition types.BaseCondition) (string, error) {
	queries := strings.Split(strings.ToLower(mainQuery), " from ")
	if len(queries) > 1 && baseCondition.Fields != nil && len(baseCondition.Fields) > 0 {
		mainQuery = "SELECT " + strings.Trim(strings.Join(baseCondition.Fields, ", "), "[]") + " FROM " + queries[1]
	}
	return generateQueryParameter(mainQuery, baseCondition)
}

func generateQueryParameter(mainQuery string, baseCondition types.BaseCondition) (string, error) {
	conditionQuery, err := generateWhereParameter(baseCondition.Conditions)
	if err != nil {
		return "", err
	}

	sortQuery, limitQuery := generateSortLimit(baseCondition.Footer.Page, baseCondition.Footer.Limit, baseCondition.Footer.Sort)
	mainQuery += " " + conditionQuery + sortQuery + " " + limitQuery
	return mainQuery, nil
}

func generateWhereParameter(conditions []*types.Condition) (string, error) {
	var queryBuffer bytes.Buffer
	for i, condition := range conditions {
		isFirst := false
		if i == 0 {
			isFirst = true
		}
		err := buildWhereParameter(condition.Conditions, &queryBuffer, false, isFirst)
		if err != nil {
			return "", err
		}
	}
	return queryBuffer.String(), nil
}

func generateSortLimit(page int, limit int, sort map[string]string) (string, string) {
	querySort := ""
	sortCount := 0
	for key, value := range sort {
		sortCount++
		if sortCount == 1 {
			querySort = fmt.Sprintf(" ORDER BY %s %s", key, value)
			continue
		}
		querySort = fmt.Sprintf("%s ,%s %s", querySort, key, value)
	}

	queryLimit := ""
	if limit > 0 {
		offset := (page * limit) - limit
		queryLimit = fmt.Sprintf(` LIMIT %d OFFSET %d `, limit, offset)
	}
	return querySort, queryLimit
}

func buildWhereParameter(conditions []*types.Condition, queryBuffer *bytes.Buffer, isGroup, isFirst bool) error {
	logicalOperator := "WHERE"
	if isGroup {
		logicalOperator = ""
	}

	for i, condition := range conditions {
		if condition.Conditions != nil {
			conditionLength := len(condition.Conditions)
			if conditionLength > 0 {
				var operator string
				enableGroup := conditionLength > 1
				operator = condition.Operator
				if !isGroup && i == 0 && isFirst {
					operator = "WHERE"
				}
				var buffer bytes.Buffer
				err := buildWhereParameter(condition.Conditions, &buffer, true, isFirst)
				if err != nil {
					continue
				}

				queryBuffer.WriteString(operator)
				if enableGroup {
					queryBuffer.WriteByte(' ')
					queryBuffer.WriteByte('(')
				}
				queryBuffer.WriteString(buffer.String())
				if enableGroup {
					queryBuffer.WriteByte(')')
					queryBuffer.WriteByte(' ')
				}
				continue
			}
		}

		if condition.Attribute == nil {
			continue
		}
		err := assignAndValidateOperator(condition)
		if err != nil {
			return err
		}
		queryValue, err := assignQueryValue(condition.Attribute)
		if err != nil {
			return err
		}
		if i > 0 {
			logicalOperator = condition.Operator
		}

		queryBuffer.WriteString(logicalOperator)
		queryBuffer.WriteByte(' ')
		queryBuffer.WriteString(condition.Attribute.Name)
		queryBuffer.WriteByte(' ')
		queryBuffer.WriteString(condition.Attribute.Operator)
		queryBuffer.WriteByte(' ')
		queryBuffer.WriteString(queryValue)
		queryBuffer.WriteByte(' ')
	}
	return nil
}

func assignQueryValue(attribute *types.Attribute) (value string, err error) {
	if attribute == nil {
		return "", fmt.Errorf(consts.ErrorMessageInvalidParameter, "attribute")
	}
	switch attribute.Operator {
	case consts.OperatorInclude, consts.OperatorExclude:
		value = "(" + assignCollectionValueByAttributeType(attribute.Type, attribute.Value) + ")"
	case consts.OperatorIsNull, consts.OperatorIsNotNull:
		value = ""
	default:
		value = assignValueByAttributeType(attribute.Type, attribute.Value)
	}
	return
}

func assignCollectionValueByAttributeType(attrType valuetype.ValueType, attrValue string) string {
	if attrType != valuetype.Numeric {
		attrs := strings.Split(attrValue, ",")
		var strData []string
		for _, attr := range attrs {
			strData = append(strData, "'"+strings.ReplaceAll(attr, "'", "''")+"'")
		}
		attrValue = strings.Trim(strings.Join(strData, ","), "[]")
	}
	return attrValue
}

func assignValueByAttributeType(attrType valuetype.ValueType, attrValue string) string {
	if attrType != valuetype.Numeric {
		attrValue = "'" + strings.ReplaceAll(attrValue, "'", "''") + "'"
	}
	return attrValue
}

func assignAndValidateOperator(condition *types.Condition) error {
	if condition == nil || condition.Attribute == nil {
		return fmt.Errorf(consts.ErrorMessageInvalidParameter, "condition")
	}
	if condition.Operator == "" {
		condition.Operator = consts.LogicalOperatorAnd
	}
	if !isValidFilterLogicalOperator(condition.Operator) {
		return errors.New(fmt.Sprintf("Invalid logical operator: %s", condition.Operator))
	}

	if condition.Attribute.Operator == "" {
		condition.Attribute.Operator = consts.OperatorEqual
	}
	if !isValidFilterOperator(condition.Attribute.Operator) {
		return errors.New(fmt.Sprintf("Invalid operator: %s", condition.Attribute.Operator))
	}
	return nil
}

func isValidFilterOperator(currOperator string) bool {
	if _, ok := operatorMap[currOperator]; ok {
		return true
	}
	return false
}

func isValidFilterLogicalOperator(currOperator string) bool {
	if _, ok := logicalOperatorMap[currOperator]; ok {
		return true
	}
	return false
}
