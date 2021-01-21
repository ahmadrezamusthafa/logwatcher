package thirdparty

import (
	"github.com/ahmadrezamusthafa/logwatcher/common/errors"
	"github.com/ahmadrezamusthafa/logwatcher/common/respwriter"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/ahmadrezamusthafa/logwatcher/domain/service/thirdparty"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	Config            config.Config       `inject:"config"`
	ThirdPartyService *thirdparty.Service `inject:"thirdPartyService"`
}

func (h *Handler) GetLogAttributes(c *gin.Context) {
	var (
		r          = c.Request
		w          = c.Writer
		err        error
		ctx        = r.Context()
		respWriter = respwriter.New()
	)
	defer func() {
		if err != nil {
			respWriter.ErrorWriter(w, errors.GetHttpStatus(err), "en", err)
		}
	}()

	resp, err := h.ThirdPartyService.GetLogAttributes(ctx)
	if err != nil {
		err = errors.AddTrace(err)
		return
	}
	respWriter.SuccessWriter(w, http.StatusOK, resp)
}

func (h *Handler) GenerateQuery(c *gin.Context) {
	var (
		r          = c.Request
		w          = c.Writer
		err        error
		param      QueryParam
		ctx        = r.Context()
		respWriter = respwriter.New()
	)
	defer func() {
		if err != nil {
			respWriter.ErrorWriter(w, errors.GetHttpStatus(err), "en", err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.ReadDataError
		return
	}
	err = jsoniter.Unmarshal(body, &param)
	if err != nil {
		err = errors.AddTrace(err)
		return
	}

	queryInput := thirdparty.QueryInput{
		MessageQuery: param.MessageQuery,
		ContextQuery: param.ContextQuery,
	}
	generatedQuery, err := h.ThirdPartyService.GenerateQuery(ctx, param.Service, queryInput)
	if err != nil {
		err = errors.AddTrace(err)
		return
	}
	respWriter.SuccessWriter(w, http.StatusOK, generatedQuery)
}

func (h *Handler) Query(c *gin.Context) {
	var (
		r          = c.Request
		w          = c.Writer
		err        error
		param      QueryParam
		ctx        = r.Context()
		respWriter = respwriter.New()
	)
	defer func() {
		if err != nil {
			respWriter.ErrorWriter(w, errors.GetHttpStatus(err), "en", err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.ReadDataError
		return
	}
	err = jsoniter.Unmarshal(body, &param)
	if err != nil {
		err = errors.AddTrace(err)
		return
	}

	queryInput := thirdparty.QueryInput{
		MessageQuery: param.MessageQuery,
		ContextQuery: param.ContextQuery,
	}
	output, err := h.ThirdPartyService.Query(ctx, param.Service, queryInput)
	if err != nil {
		err = errors.AddTrace(err)
		return
	}
	respWriter.SuccessWriter(w, http.StatusOK, output)
}
