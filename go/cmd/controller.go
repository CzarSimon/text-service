package main

import (
	"net/http"

	"github.com/CzarSimon/text-service/go/pkg/utils/context"
	"github.com/CzarSimon/text-service/go/pkg/utils/httputil"
	"github.com/CzarSimon/text-service/go/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

var controllerLog = logger.GetDefaultLogger("cmd/controller").Sugar()

func (e *env) getTextByKey(c *gin.Context) {
	ctx := createContext(c)
	controllerLog.Debugw("getTextByKey", "ctx", ctx)

	if ctx.Language == "" {
		c.Error(httputil.BadRequest("No language specified"))
		return
	}

	texts, err := e.textGetter.Get(ctx, c.Param("key"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, texts)
}

func (e *env) getTextGroup(c *gin.Context) {
	ctx := createContext(c)
	controllerLog.Debugw("getTextByKey", "ctx", ctx)

	if ctx.Language == "" {
		c.Error(httputil.BadRequest("No language specified"))
		return
	}

	texts, err := e.textGetter.GetGroup(ctx, c.Param("groupId"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, texts)
}

func createContext(c *gin.Context) *context.Context {
	requestID := httputil.GetRequestID(c)
	locale := httputil.GetLocale(c)

	return context.New(c.Request.Context(), requestID, locale)
}
