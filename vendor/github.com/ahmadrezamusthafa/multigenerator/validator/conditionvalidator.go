package validator

import (
	"github.com/ahmadrezamusthafa/multigenerator/shared/consts"
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"github.com/ahmadrezamusthafa/multigenerator/shared/utils"
	"strings"
	"time"
)

type Condition struct {
	*types.Condition
}

func (c *Condition) ValidateCondition(condition types.Condition) (isValid bool, err error) {
	referenceAttrMap := make(map[string]bool)
	inputAttrMap := make(map[string]bool)

	readAllAttributes(c.Condition, referenceAttrMap)
	readAllAttributes(&condition, inputAttrMap)
	setNonExistAttributeDefaultValue(&condition, referenceAttrMap, inputAttrMap)
	return c.validateConditionAttribute(condition)
}

func readAllAttributes(condition *types.Condition, attrMap map[string]bool) {
	if len(condition.Conditions) > 0 {
		for _, condition := range condition.Conditions {
			readAllAttributes(condition, attrMap)
		}
	} else {
		if condition.Attribute != nil {
			if _, ok := attrMap[condition.Attribute.Name]; !ok {
				attrMap[condition.Attribute.Name] = true
			}
		}
	}
}

func (c *Condition) validateConditionAttribute(inputCondition types.Condition) (isValid bool, err error) {
	if len(c.Conditions) > 0 {
		for i, subCondition := range c.Conditions {
			con := Condition{Condition: subCondition}
			isSubValid, err := con.validateConditionAttribute(inputCondition)
			if err != nil {
				return false, err
			}
			if i == 0 {
				isValid = isSubValid
			} else {
				if subCondition.Operator == consts.LogicalOperatorOr {
					isValid = isValid || isSubValid
				} else {
					isValid = isValid && isSubValid
				}
			}
		}
	} else {
		isValid, _, err = c.validateConditionValue("", inputCondition)
	}
	return
}

func (c *Condition) validateConditionValue(prefix string, condition types.Condition) (isValid, isSkip bool, err error) {
	isValid = true
	if len(condition.Conditions) > 0 {
		for i, subCondition := range condition.Conditions {
			isSubValid, isSkip, err := c.validateConditionValue(prefix, *subCondition)
			if err != nil {
				return false, false, err
			}
			if isSkip {
				continue
			}
			if i == 0 {
				isValid = isSubValid
			} else {
				if subCondition.Operator == consts.LogicalOperatorOr {
					isValid = isValid || isSubValid
				} else {
					isValid = isValid && isSubValid
				}
			}
		}
	} else {
		if c.Attribute == nil || condition.Attribute == nil {
			return false, false, nil
		}
		if condition.Attribute.Name == c.Attribute.Name {
			operator := c.Attribute.Operator
			switch operator {
			case consts.OperatorEqual:
				isValid = strings.EqualFold(condition.Attribute.Value, c.Attribute.Value)
			default:
				value := condition.Attribute.Value
				secondValue := c.Attribute.Value
				valueType := getValueType(c.Attribute.Value)

				switch valueType {
				case valuetype.Date:
					isValid = validateTime(utils.StringToTime(value), operator, utils.StringToTime(secondValue))
				default:
					isValid = validateNumeric(utils.StringToFloat64(value), operator, utils.StringToFloat64(secondValue))
				}
			}
		} else {
			return false, true, nil
		}
	}
	return
}

func setNonExistAttributeDefaultValue(condition *types.Condition, referenceAttrMap, inputAttrMap map[string]bool) {
	for attrName, _ := range referenceAttrMap {
		if _, ok := inputAttrMap[attrName]; !ok {
			condition.Conditions = append(condition.Conditions, &types.Condition{
				Operator: consts.LogicalOperatorAnd,
				Attribute: &types.Attribute{
					Name:     attrName,
					Operator: "=",
					Value:    "",
				},
			})
		}
	}
}

func getValueType(value string) valuetype.ValueType {
	varType, indexVal, dotCount := valuetype.Alphanumeric, 0, 0
	for _, char := range value {
		if char == ',' {
			continue
		}
		if '0' <= char && char <= '9' {
			if indexVal == 0 || (indexVal > 0 && dotCount == 1) {
				varType = valuetype.Numeric
			}
		} else if char == '.' {
			if indexVal > 0 && varType == valuetype.Numeric {
				dotCount++
				varType = valuetype.Alphanumeric
			}
			if dotCount > 1 {
				varType = valuetype.Alphanumeric
				break
			}
		} else {
			varType = valuetype.Alphanumeric
			break
		}
		indexVal++
	}
	if varType == valuetype.Alphanumeric {
		if _, err := time.Parse(consts.DateTimeFormat, value); err == nil {
			varType = valuetype.Date
		}
	}
	return varType
}
