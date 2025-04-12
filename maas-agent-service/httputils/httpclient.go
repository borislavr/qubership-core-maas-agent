package httputils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/netcracker/qubership-core-lib-go/v3/context-propagation/ctxhelper"
	"github.com/netcracker/qubership-core-lib-go/v3/logging"
	tlsUtils "github.com/netcracker/qubership-core-lib-go/v3/utils"
	"github.com/valyala/fasthttp"
)

const HEADER_X_REQUEST_ID = "X-Request-Id"

var (
	fasthttpClient = &fasthttp.Client{
		TLSConfig:                     tlsUtils.GetTlsConfig(),
		MaxIdleConnDuration:           30 * time.Second,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		DialDualStack:                 true,
		WriteTimeout:                  30 * time.Second,
		ReadTimeout:                   130 * time.Second, // should be greater than watch requests timeout 120sec
	}
	log logging.Logger
)

func init() {
	log = logging.GetLogger("utils/http-client")
}

type HttpRequest struct {
	method             string
	url                string
	headers            http.Header
	requestBody        *string
	requestBodyBytes   []byte
	authHeaderProvider func(ctx context.Context) (string, error)
}

func (c HttpRequest) String() string {
	return fmt.Sprintf("%s %s", c.method, c.url)
}

func Req(method string, url string, authHeaderProvider func(ctx context.Context) (string, error)) *HttpRequest {
	return &HttpRequest{method: method, url: url, headers: http.Header{}, authHeaderProvider: authHeaderProvider}
}

func (c *HttpRequest) SetRequestBody(body *string) *HttpRequest {
	c.requestBody = body
	return c
}

func (c *HttpRequest) SetRequestBodyBytes(body []byte) *HttpRequest {
	c.requestBodyBytes = body
	return c
}

func (c *HttpRequest) SetHeader(name string, value string) *HttpRequest {
	c.headers.Set(name, value)
	return c
}

func (c *HttpRequest) AddHeader(name string, value string) *HttpRequest {
	c.headers.Add(name, value)
	return c
}

func (c *HttpRequest) Execute(ctx context.Context) (code int, body []byte, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(c.method)
	req.SetRequestURI(c.url)

	if c.requestBody != nil {
		req.SetBodyString(*c.requestBody)
	}
	if len(c.requestBodyBytes) > 0 {
		req.SetBody(c.requestBodyBytes)
	}

	if err := ctxhelper.AddSerializableContextData(ctx, c.headers.Add); err != nil {
		return -1, nil, fmt.Errorf("cannot add context data: %w", err)
	}

	if c.headers.Get(HEADER_X_REQUEST_ID) == "" {
		requestId := uuid.New().String()
		log.WarnC(ctx, "No request id is set. Create new one: %s", requestId)
		c.headers.Add(HEADER_X_REQUEST_ID, requestId)
	}

	for name := range c.headers {
		for _, value := range c.headers.Values(name) {
			req.Header.Add(name, value)
		}
	}

	if c.authHeaderProvider != nil {
		headerValue, err := c.authHeaderProvider(ctx)
		if err != nil {
			return -1, nil, fmt.Errorf("error getting M2M token: %w", err)
		}

		log.DebugC(ctx, "Request authorization: %v", headerValue)
		req.Header.Set("Authorization", headerValue)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	log.InfoC(ctx, "Execute request: %+v", c)
	log.DebugC(ctx, "Underlying client: %+v", req)
	err = fasthttpClient.Do(req, resp)
	log.DebugC(ctx, "Response: %v", req)

	if err != nil {
		log.ErrorC(ctx, "Error execute http request: %v", err)
		return -1, nil, err
	}

	body = make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	return resp.StatusCode(), body, nil
}
