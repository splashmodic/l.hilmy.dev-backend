package converter

import (
	"github.com/bytedance/sonic"
	"l.hilmy.dev/backend/helpers/errorhandler"
)

func MapToStruct[T comparable](m *map[string]interface{}, s *T) *error {
	marshalled, err := sonic.MarshalString(&m)
	if err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return &err
	}
	if err := sonic.UnmarshalString(marshalled, &s); err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return &err
	}
	return nil
}

func StructToMap[T comparable](s *T, m *map[string]interface{}) error {
	marshalled, err := sonic.MarshalString(&s)
	if err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return err
	}
	if err := sonic.UnmarshalString(marshalled, &m); err != nil {
		errorhandler.LogErrorThenContinue(&err)
		return err
	}
	return nil
}
