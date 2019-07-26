package context_test

import (
	"context"
	"fmt"
	"testing"

	myCtx "github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/CzarSimon/text-service/go/pkg/utils/id"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	parent := context.Background()

	ctxID := id.New()
	ctx := myCtx.New(parent, ctxID, "sv")
	var stdctx context.Context = ctx

	assert.Equal(ctxID, ctx.ID)

	val := stdctx.Value(myCtx.ContextIDKey)
	stdctxID, ok := val.(string)
	assert.True(ok, fmt.Sprintf("context.Value returned unexpected type. Expected: stirng Got: %v", val))

	assert.Equal(ctx.ID, stdctxID)
}
