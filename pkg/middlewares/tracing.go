package middleware

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go-micro.dev/v5/metadata"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	HeaderXRequestID    = "X-Request-ID"
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-IP"

	RequestIDKey = "request_id"
	ClientIPKey  = "client_ip"
	TimestampKey = "timestamp"
	HostIDKey    = "host_id"
)

// RequestContext holds all request-specific information
type RequestContext struct {
	RequestID string
	ClientIP  string
	Timestamp time.Time
	HostID    string
}

// RequestID middleware
type RequestID struct {
	handler http.Handler
	prefix  string
	hostID  string
}

func NewRequestID(prefix string) func(http.Handler) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	hostID := fmt.Sprintf("%s-%s", hostname, uuid.New().String()[:8])

	return func(next http.Handler) http.Handler {
		return &RequestID{
			handler: next,
			prefix:  prefix,
			hostID:  hostID,
		}
	}
}

func (m *RequestID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().UTC()
	clientIP := m.getClientIP(r)

	// Generate or get request ID
	requestID := m.getRequestID(r, clientIP, timestamp)

	// Set headers
	w.Header().Set(HeaderXRequestID, requestID)

	// Create new context with all values
	ctx := r.Context()
	md := make(metadata.Metadata)

	md.Set(RequestIDKey, requestID)
	md.Set(ClientIPKey, clientIP)
	md.Set(TimestampKey, cast.ToString(timestamp))
	md.Set(HostIDKey, m.hostID)
	ctx = metadata.MergeContext(ctx, md, true)

	// Call next handler with enriched context
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}

// getRequestID generates or gets a request ID
// Format: ip-prefix.service-prefix.timestamp.host.uuid
func (m *RequestID) getRequestID(r *http.Request, clientIP string, timestamp time.Time) string {
	// Check if request ID already exists
	if reqID := r.Header.Get(HeaderXRequestID); reqID != "" {
		return reqID
	}

	// Create IP prefix - using last octet for IPv4 or last segment for IPv6
	ipPrefix := m.createIPPrefix(clientIP)

	// Generate UUID
	uid := strings.ReplaceAll(uuid.New().String(), "-", "")

	// Format: ip-prefix.service-prefix.timestamp.host.uuid
	return fmt.Sprintf("%s.%s.%s.%s",
		ipPrefix,
		m.prefix,
		//m.hostID,
		timestamp.Format("20060102150405"),
		uid,
	)
}

// createIPPrefix creates a searchable prefix from IP
func (m *RequestID) createIPPrefix(ip string) string {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return "unknown"
	}

	// Convert to 4-byte format for both IPv4 and IPv6
	ipv4 := ipAddr.To4()
	if ipv4 != nil {
		// For IPv4, use the format: 192-168-1-2
		return strings.ReplaceAll(ip, ".", "-")
	}

	// For IPv6, use the last 4 segments
	ipv6Segments := strings.Split(ipAddr.String(), ":")
	if len(ipv6Segments) > 4 {
		return strings.Join(ipv6Segments[len(ipv6Segments)-4:], "-")
	}

	return strings.Join(ipv6Segments, "-")
}

// getClientIP extracts the real client IP
func (m *RequestID) getClientIP(r *http.Request) string {
	// Check X-Real-IP header
	if ip := r.Header.Get(HeaderXRealIP); ip != "" {
		return m.cleanIP(ip)
	}

	// Check X-Forwarded-For header
	if forwardedFor := r.Header.Get(HeaderXForwardedFor); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return m.cleanIP(ips[0])
		}
	}

	// Fall back to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return m.cleanIP(ip)
}

func (m *RequestID) cleanIP(ip string) string {
	return strings.TrimSpace(ip)
}

// Helper functions to get values from context
func GetRequestContextFromCtx(ctx context.Context) *RequestContext {
	if ctx == nil {
		return nil
	}
	return &RequestContext{
		RequestID: GetRequestIDFromCtx(ctx),
		ClientIP:  GetClientIPFromCtx(ctx),
		Timestamp: GetTimestampFromCtx(ctx),
		HostID:    GetHostIDFromCtx(ctx),
	}
}

func GetRequestIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return v
	}
	return ""
}

func GetClientIPFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ClientIPKey).(string); ok {
		return v
	}
	return ""
}

func GetTimestampFromCtx(ctx context.Context) time.Time {
	if v, ok := ctx.Value(TimestampKey).(time.Time); ok {
		return v
	}
	return time.Time{}
}

func GetHostIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(HostIDKey).(string); ok {
		return v
	}
	return ""
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := GetRequestContextFromCtx(r.Context())

		// Now you can easily search logs by IP prefix
		log.Printf("[%s] Request from %s",
			reqCtx.RequestID, // Contains IP as prefix
			reqCtx.ClientIP,
		)

		next.ServeHTTP(w, r)
	})
}
