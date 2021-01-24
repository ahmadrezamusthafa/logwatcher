package multigenerator

import (
	"github.com/ahmadrezamusthafa/multigenerator/querygen"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"github.com/ahmadrezamusthafa/multigenerator/structgen"
	"github.com/ahmadrezamusthafa/multigenerator/validator"
)

/*
GenerateCondition
-----------------------------------------------------------------------
is a function to generate condition object as validator that used by
 - Validate
 - ValidateObjects
 - ValidateCondition

Param:
@astQuery is abstract syntax tree query
*/
func GenerateCondition(astQuery string) (types.Condition, error) {
	var gen structgen.StructGen
	return gen.GenerateCondition(astQuery)
}

func Validate(referenceCondition types.Condition, data interface{}) (isValid bool, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.Validate(data)
}

func ValidateObjects(referenceCondition types.Condition, data ...interface{}) (isValid bool, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.ValidateObjects(data...)
}

func ValidateCondition(referenceCondition types.Condition, inputCondition types.Condition) (isValid bool, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.ValidateCondition(inputCondition)
}

func FilterSlice(referenceCondition types.Condition, data interface{}) (result interface{}, err error) {
	con := validator.Condition{Condition: &referenceCondition}
	return con.FilterSlice(data)
}

/*
GenerateQuery
-----------------------------------------------------------------------
is a function to generate SQL query

Param:
@mainQuery is a parent query
@baseCondition is a condition object with header and footer info
*/
func GenerateQuery(mainQuery string, baseCondition types.BaseCondition) (string, error) {
	var gen querygen.QueryGen
	return gen.GenerateQuery(mainQuery, baseCondition)
}
