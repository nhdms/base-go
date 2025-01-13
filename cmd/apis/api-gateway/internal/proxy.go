package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/rehttp"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/nhdms/base-go/internal"
	"github.com/nhdms/base-go/internal/token"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/common"
	"github.com/nhdms/base-go/pkg/logger"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"github.com/nhdms/base-go/pkg/utils/token_helper"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/selector"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ReverseProxy struct {
	reg            registry.Registry
	balancer       selector.Selector
	tokenProcessor token.Processor
}

const (
	MaxIdleConnectPerHost = 100
	MaxRetries            = 5
	RetryDelay            = 100 * time.Millisecond
)

var pxyTransport = getDefaultTransport()

func NewReverseProxy(redis *redis.Client, reg registry.Registry, balancer selector.Selector) *ReverseProxy {
	if balancer == nil {
		balancer = selector.NewSelector(selector.Registry(reg))
	}
	userClient := internal.CreateNewUserServiceClient(nil)
	tp := token.NewTokenProcessor(redis, userClient)

	return &ReverseProxy{reg: reg, balancer: balancer, tokenProcessor: tp}
}

func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) == 0 {
		transhttp.RespondJSONFull(w, http.StatusBadRequest, common.NewErrorHTTPResponse("Service name not provided"))
		return
	}
	serviceName := pathParts[0]

	start := time.Now()
	node, err := p.getNode(serviceName)
	timeToDiscovery := time.Since(start)
	if err != nil {
		logger.DefaultLogger.Errorw("Error when get service and node, details", "error", err)
		transhttp.RespondJSONFull(w, http.StatusInternalServerError, common.NewErrorHTTPResponse("no service found"))
		return
	}

	svcAddr := p.getServiceAddress(node, r.URL.Path)
	rp, err := url.Parse(svcAddr)
	if err != nil {
		logger.DefaultLogger.Errorw("Error when parse service url, details", "error", err)
		transhttp.RespondJSONFull(w, http.StatusInternalServerError, common.NewErrorHTTPResponse("server error"))
		return
	}

	if transhttp.IsWebSocket(r) {
		//p.serveWebSocket(rp.Host, w, r)
		transhttp.RespondJSONFull(w, http.StatusInternalServerError, common.NewErrorHTTPResponse("maintaining"))
		return
	}

	matchedEndpoint := GetMatchedEndpoints(r, node.Metadata)
	if matchedEndpoint == nil {
		transhttp.RespondJSONFull(w, http.StatusNotFound, common.NewErrorHTTPResponse("no route found"))
		return
	}

	p.cleanPrivateRequestHeader(r) // to prevent user fake header
	err = p.extractAndVerifyTokenInfo(r, matchedEndpoint)
	if err != nil {
		if errors.Is(err, common.UnauthorizedError) {
			transhttp.RespondJSONFull(w, http.StatusUnauthorized, common.NewErrorCodeHTTPResponse(common.AuthCodeUnauthorized))
			return
		}
		transhttp.RespondJSONFull(w, http.StatusInternalServerError, common.NewErrorHTTPResponse(err.Error()))
		return
	}

	responseStatusCode := http.StatusOK
	rw := transhttp.NewRecorderResponseWriter(w, responseStatusCode)

	// trim service path from request path
	trimmedPath := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("/%s", serviceName))
	r.URL.Path = trimmedPath

	// add user info to request
	pxy := httputil.NewSingleHostReverseProxy(rp)

	pxy.ModifyResponse = func(resp *http.Response) error {
		// todo handle error rate (alerts)
		responseStatusCode = resp.StatusCode
		return nil
	}

	pxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, context.Canceled) {
			// the client has closed the connection
			logger.DefaultLogger.Infow("The client has closed the connection", "service affected", serviceName,
				"clientIP", p.getOriginClientIP(r),
				"requestUri", r.RequestURI,
				"processing_time (ms)", time.Since(timeStart).Milliseconds())

			transhttp.RespondError(rw, http.StatusInternalServerError, err.Error())
			return
		}

		logger.DefaultLogger.Errorw("Error when sending request to dest service",
			"error", err.Error(),
			"service_name", serviceName,
			"request_uri", r.RequestURI,
			"request_method", r.Method,
			"addr", svcAddr,
		)
		transhttp.RespondError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	pxy.Transport = rehttp.NewTransport(
		pxyTransport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(MaxRetries),
			rehttp.RetryIsErr(func(e error) bool {
				return e != nil && (strings.Contains(e.Error(), "timeout") || strings.Contains(e.Error(), "connection reset by peer"))
			}),
		),
		rehttp.ConstDelay(RetryDelay),
	)

	beforeProxy := time.Now()
	pxy.ServeHTTP(rw, r)
	proxyDuration := time.Since(beforeProxy)
	totalProcessed := time.Since(timeStart)
	logger.DefaultLogger.Debugw("Processed request",
		"request_uri", r.RequestURI,
		"status", responseStatusCode,
		"service_looking", timeToDiscovery.Milliseconds(),
		"took", proxyDuration.Milliseconds(),
		"total_processed", totalProcessed.Milliseconds(),
	)
}

func (p *ReverseProxy) getServiceAddress(node *registry.Node, path string) string {
	return fmt.Sprintf("http://%s", node.Address)
}

func (p *ReverseProxy) getOriginClientIP(r *http.Request) string {
	// todo handle IP from Cloudflare here
	if realIP := r.Header.Get("X-Real-IP"); len(realIP) > 0 {
		return realIP
	}

	if realIP := r.Header.Get("X-Forwarded-For"); len(realIP) > 0 {
		if strings.Contains(realIP, ",") {
			tmp := strings.Split(realIP, ",")
			return strings.TrimSpace(tmp[len(tmp)-1])
		}

		return realIP
	}
	return ""
}

func (p *ReverseProxy) getNode(serviceName string) (*registry.Node, error) {
	//var service *registry.Service
	next, err := p.balancer.Select(app.GetAPIName(serviceName))
	if err != nil {
		return nil, err
	}

	node, err := next()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (p *ReverseProxy) extractAndVerifyTokenInfo(r *http.Request, endpoint *transhttp.Route) error {
	if !endpoint.AuthInfo.Enable {
		return nil
	}

	if p.isRequestBeWhiteListed(r) {
		return nil
	}

	tokenString := token_helper.ExtractTokenFromRequest(r)
	if len(tokenString) == 0 {
		return common.UnauthorizedError
	}

	jwtToken, err := p.tokenProcessor.GetToken(r.Context(), tokenString)
	if err != nil {
		return err
	}

	if jwtToken == nil {
		return common.UnauthorizedError
	}

	if granted := p.tokenProcessor.CheckPermissions(jwtToken, endpoint.AuthInfo.RequirePermissions); !granted {
		return common.UnauthorizedError
	}

	metadataToSet := p.tokenProcessor.ExtractMetadata(jwtToken)
	p.setMetadataToHeader(r, metadataToSet)

	return nil
}

func (p *ReverseProxy) isRequestBeWhiteListed(r *http.Request) bool {
	return false
}

func (p *ReverseProxy) setMetadataToHeader(r *http.Request, set map[string]string) {
	for k, v := range set {
		r.Header.Set(k, v)
	}
}

func (p *ReverseProxy) cleanPrivateRequestHeader(r *http.Request) {
	for _, header := range common.ExtraDataHeaders {
		r.Header.Set(header, "")
	}
}

func GetMatchedEndpoints(request *http.Request, meta map[string]string) *transhttp.Route {
	tmpRouter := mux.NewRouter().StrictSlash(true)
	endpoints := make([]*transhttp.Route, 0)
	_ = json.Unmarshal([]byte(meta["endpoints"]), &endpoints)
	basePath := meta["base_path"]

	for _, e := range endpoints {
		path := basePath + e.Pattern
		method := e.Method
		route := tmpRouter.Path(path).Methods(method)

		var match mux.RouteMatch
		ok := route.Match(request, &match)
		if ok {
			return e
		}
	}
	return nil
}
