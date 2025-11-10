package codehost

import "context"

type CodeHostFactory interface {
	Create(ctx context.Context, accessToken string) CodeHost
}
