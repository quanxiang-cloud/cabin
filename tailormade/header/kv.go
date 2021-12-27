package header

import "context"

const (
	requestID = "Request-Id"

	timezone = "timezone"
)

type key string

func MutateContext(c CTX) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, key(requestID), c.GetHeader(requestID))
	ctx = context.WithValue(ctx, key(requestID), c.GetHeader(timezone))

	return ctx
}

func GetRequestIDKV(ctx context.Context) (string, string) {
	i := ctx.Value(requestID)
	rid, ok := i.(string)
	if ok {
		return requestID, rid
	}
	return requestID, "unexpected type"
}

func GetTimezone(ctx context.Context) (string, string) {
	i := ctx.Value(timezone)
	tz, ok := i.(string)
	if ok {
		return timezone, tz
	}
	return timezone, "unexpected type"
}
