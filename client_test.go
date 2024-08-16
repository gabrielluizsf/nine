package nine

import (
	"context"
	"testing"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	client := New(ctx)
	if client.Context() != ctx{
		t.Fatal("invalid context")
	}
}
