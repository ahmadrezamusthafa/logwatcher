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
		r           = c.Request
		w           = c.Writer
		err         error
		serviceName string
		sourceName  string
		ctx         = r.Context()
		queries     = r.URL.Query()
		respWriter  = respwriter.New()
	)
	defer func() {
		if err != nil {
			respWriter.ErrorWriter(w, errors.GetHttpStatus(err), "en", err)
		}
	}()

	if val, ok := queries["service"]; ok {
		serviceName = val[0]
	}
	if val, ok := queries["source"]; ok {
		sourceName = val[0]
	}

	resp, err := h.ThirdPartyService.GetLogAttributes(ctx, serviceName, sourceName)
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
	generatedQuery, err := h.ThirdPartyService.GenerateQuery(ctx, param.Service, param.Source, queryInput, param.Limit)
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
	output, err := h.ThirdPartyService.Query(ctx, param.Service, param.Source, queryInput, param.Limit)
	if err != nil {
		err = errors.AddTrace(err)
		return
	}
	respWriter.SuccessWriter(w, http.StatusOK, output)
}
