package header

import "context"

const (
	requestID = "Request-Id"
)

type key string

func MutateContext(c CTX) context.Context {
	rid := c.GetHeader(requestID)
	ctx := context.Background()
	return context.WithValue(ctx, key(requestID), rid)
}

func GetRequestIDKV(ctx context.Context) (string, string) {
	i := ctx.Value(requestID)
	rid, ok := i.(string)
	if ok {
		return requestID, rid
	}
	return requestID, "unexpected type"
}
