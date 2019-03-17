package test

import (
	"testing"

	"github.com/davebryson/menta/app"
	sdk "github.com/davebryson/menta/types"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	assert := assert.New(t)

	router := app.NewRouter()
	router.Add("dave", func(ctx sdk.Context) sdk.Result {
		return sdk.Result{Log: "dave"}
	})
	router.Add("bob", func(ctx sdk.Context) sdk.Result {
		return sdk.Result{Log: "bob"}
	})
	router.Add("carl", func(ctx sdk.Context) sdk.Result {
		return sdk.Result{Log: "carl"}
	})

	ch := router.GetHandler("carl")
	assert.NotNil(ch)
	assert.Equal("carl", ch(sdk.Context{}).Log)

	assert.Nil(router.GetHandler("cl"))

	dh := router.GetHandler("dave")
	assert.NotNil(dh)
	assert.Equal("dave", dh(sdk.Context{}).Log)
}
