package goerror

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoError_Reasons(t *testing.T) {
	t.Run("GetEmptyReason", func(t *testing.T) {
		e := DefineBadRequest("UserExist", "user is already exist")

		require.Len(t, e.reasons, 0)
		require.Nil(t, e.reasons)
	})

	t.Run("AddAndGetReason", func(t *testing.T) {
		e := DefineBadRequest("UserExist", "user is already exist")
		e.AddReason("username", "username already exist", nil)

		reasons := e.GetReasons()
		require.Len(t, reasons, 1)
		require.Equal(t, reasons[0], &Reason{FieldName: "username", Reason: "username already exist", Value: nil})
	})

	t.Run("AddAndGetReasons", func(t *testing.T) {
		e := DefineBadRequest("UserExist", "user is already exist")
		e.AddReason("username", "username already exist", nil)
		e.AddReason("phone", "phone number already exist", "0598881111")

		reasons := e.GetReasons()
		require.Len(t, reasons, 2)
		require.Equal(t, reasons[0], &Reason{FieldName: "username", Reason: "username already exist", Value: nil})
		require.Equal(t, reasons[1], &Reason{FieldName: "phone", Reason: "phone number already exist", Value: "0598881111"})
	})
}
