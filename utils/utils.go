package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/diogenes-moreira/dokan-go-sdk/errors"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RequestOptions contains options for making HTTP requests
type RequestOptions struct {
	Method  string
	Path    string
	Query   interface{}
	Body    interface{}
	Headers map[string]string
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// MakeRequest makes an HTTP request with the given options
func MakeRequest(ctx context.Context, client HTTPClient, baseURL string, opts RequestOptions) (*Response, error) {
	// Build URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	
	u.Path = strings.TrimSuffix(u.Path, "/") + "/" + strings.TrimPrefix(opts.Path, "/")
	
	// Add query parameters
	if opts.Query != nil {
		queryParams, err := StructToURLValues(opts.Query)
		if err != nil {
			return nil, fmt.Errorf("failed to encode query parameters: %w", err)
		}
		u.RawQuery = queryParams.Encode()
	}
	
	// Prepare request body
	var body io.Reader
	if opts.Body != nil {
		jsonBody, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(jsonBody)
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, opts.Method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	if opts.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}
	
	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewNetworkError(err)
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
	}
	
	// Handle HTTP errors
	if resp.StatusCode >= 400 {
		// Try to parse Dokan error response
		var dokanErr errors.DokanError
		if err := json.Unmarshal(respBody, &dokanErr); err == nil && dokanErr.Code != "" {
			dokanErr.StatusCode = resp.StatusCode
			return response, &dokanErr
		}
		
		// Fall back to generic HTTP error
		return response, errors.HandleHTTPError(resp.StatusCode, respBody)
	}
	
	return response, nil
}

// StructToURLValues converts a struct to url.Values using struct tags
func StructToURLValues(v interface{}) (url.Values, error) {
	values := url.Values{}
	
	if v == nil {
		return values, nil
	}
	
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %T", v)
	}
	
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)
		
		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}
		
		// Get the tag
		tag := fieldType.Tag.Get("url")
		if tag == "" || tag == "-" {
			continue
		}
		
		// Parse tag options
		tagParts := strings.Split(tag, ",")
		name := tagParts[0]
		omitEmpty := false
		for _, option := range tagParts[1:] {
			if option == "omitempty" {
				omitEmpty = true
			}
		}
		
		// Skip empty values if omitempty is set
		if omitEmpty && isEmptyValue(field) {
			continue
		}
		
		// Convert field value to string
		value, err := fieldToString(field)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %s: %w", fieldType.Name, err)
		}
		
		if value != "" {
			values.Add(name, value)
		}
	}
	
	return values, nil
}

// isEmptyValue checks if a reflect.Value is empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// fieldToString converts a reflect.Value to string
func fieldToString(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.String {
			var strs []string
			for i := 0; i < v.Len(); i++ {
				strs = append(strs, v.Index(i).String())
			}
			return strings.Join(strs, ","), nil
		}
		if v.Type().Elem().Kind() == reflect.Int {
			var strs []string
			for i := 0; i < v.Len(); i++ {
				strs = append(strs, strconv.FormatInt(v.Index(i).Int(), 10))
			}
			return strings.Join(strs, ","), nil
		}
		return "", fmt.Errorf("unsupported slice type: %s", v.Type())
	case reflect.Ptr:
		if v.IsNil() {
			return "", nil
		}
		return fieldToString(v.Elem())
	case reflect.Interface:
		if v.IsNil() {
			return "", nil
		}
		// Handle time.Time specifically
		if t, ok := v.Interface().(time.Time); ok {
			return t.Format(time.RFC3339), nil
		}
		return fieldToString(v.Elem())
	default:
		// Handle time.Time specifically
		if t, ok := v.Interface().(time.Time); ok {
			return t.Format(time.RFC3339), nil
		}
		return "", fmt.Errorf("unsupported type: %s", v.Type())
	}
}

// ParseJSON parses JSON response into the given interface
func ParseJSON(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

// BuildPath builds a URL path with parameters
func BuildPath(path string, params ...interface{}) string {
	for i, param := range params {
		placeholder := fmt.Sprintf("{%d}", i)
		path = strings.ReplaceAll(path, placeholder, fmt.Sprintf("%v", param))
	}
	return path
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	Multiplier float64
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 3,
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
		Multiplier: 2.0,
	}
}

// WithRetry executes a function with retry logic
func WithRetry(ctx context.Context, config RetryConfig, fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Calculate delay with exponential backoff
			delay := time.Duration(float64(config.BaseDelay) * float64(attempt) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
			
			// Wait with context cancellation support
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
		
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		
		// Don't retry certain types of errors
		if errors.IsDokanError(lastErr) {
			dokanErr := lastErr.(*errors.DokanError)
			// Don't retry client errors (4xx) except rate limiting
			if dokanErr.StatusCode >= 400 && dokanErr.StatusCode < 500 && dokanErr.StatusCode != 429 {
				return lastErr
			}
		}
	}
	
	return lastErr
}

