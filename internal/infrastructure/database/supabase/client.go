package supabase

import (
	"context"
	"fmt"

	"github.com/nedpals/supabase-go"
)

type SupabaseService struct {
	client *supabase.Client
}

func NewSupabaseService(ctx context.Context, url string, key string) (*SupabaseService, error) {
	if url == "" || key == "" {
		return nil, fmt.Errorf("supabase url and key are required")
	}

	return &SupabaseService{
		client: supabase.CreateClient(url, key),
	}, nil
}
