package custom_errors

import (
	"errors"
)

// user-defined errors
var (
	RequestTimeoutError    = errors.New("request timeout")
	RequestValidationError = errors.New("request validation error")
	RequestProcessingError = errors.New("request processing error")
)

// internal errors
var (
	CacheDriverGetError                = errors.New("cache driver get error")
	CacheDriverSetError                = errors.New("cache driver set error")
	KeyGenerationError                 = errors.New("key generation error")
	NodeRequestJsonMarshalError        = errors.New("node request json marshal error")
	NodeRequestNewRequestError         = errors.New("node request new request error")
	NodeRequestClientDoError           = errors.New("node request client do error")
	NodeRequestReadResponseBodyError   = errors.New("node request read response body error")
	NodeResponseResultMarshalError     = errors.New("node response result marshal error")
	ProxyFromCacheResultUnmarshalError = errors.New("proxy from cache result unmarshal error")
	NodeResultMarshalError             = errors.New("node result marshal error")
	CacheDriverSetExpireError          = errors.New("cache driver set expire error")
)
