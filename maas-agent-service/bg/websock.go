package bg

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fasthttp/websocket"

	"github.com/netcracker/qubership-core-lib-go/v3/security"
	"github.com/netcracker/qubership-core-lib-go/v3/serviceloader"
)

func SecureWebSocketDial(ctx context.Context, webSocketURL url.URL, dialer *websocket.Dialer, requestHeaders http.Header) (*websocket.Conn, *http.Response, error) {
	tokenProvider := serviceloader.MustLoad[security.TokenProvider]()
	m2mToken, err := tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Can't get m2m token: %w", err)
	}
	if requestHeaders == nil {
		requestHeaders = http.Header{}
	}
	requestHeaders = addHeaderIfAbsent(ctx, requestHeaders, "Host", webSocketURL.Host)
	requestHeaders = addHeaderIfAbsent(ctx, requestHeaders, "Authorization", "Bearer "+m2mToken)
	return dialer.DialContext(ctx, webSocketURL.String(), requestHeaders)
}

func addHeaderIfAbsent(ctx context.Context, requestHeaders http.Header, headerName, headerValue string) http.Header {
	if _, ok := requestHeaders[headerName]; !ok {
		requestHeaders.Add(headerName, headerValue)
	}
	return requestHeaders
}
