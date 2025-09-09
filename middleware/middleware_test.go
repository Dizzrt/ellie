package middleware

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestChain(t *testing.T) {
	next := func(_ context.Context, req any) (any, error) {
		if req != "hello ellie!" {
			t.Errorf("unexpected req: %v", req)
		}

		i += 10
		return "null", nil
	}

	got, err := Chain(middleware1, middleware2, middleware3)(next)(context.Background(), "hello ellie!")
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if !reflect.DeepEqual(got, "null") {
		t.Errorf("got %v, want %v", got, "null")
	}

	if !reflect.DeepEqual(i, 16) {
		t.Errorf("got %v, want %v", i, 16)
	}
}

var i int

func middleware1(handler Handler) Handler {
	return func(ctx context.Context, req any) (any, error) {
		fmt.Println("enter middleware 1")
		i++
		res, err := handler(ctx, req)
		fmt.Println("exit middleware 1")
		return res, err
	}
}

func middleware2(handler Handler) Handler {
	return func(ctx context.Context, req any) (any, error) {
		fmt.Println("enter middleware 2")
		i += 2
		res, err := handler(ctx, req)
		fmt.Println("exit middleware 2")
		return res, err
	}
}

func middleware3(handler Handler) Handler {
	return func(ctx context.Context, req any) (any, error) {
		fmt.Println("enter middleware 3")
		i += 3
		res, err := handler(ctx, req)
		fmt.Println("exit middleware 3")
		return res, err
	}
}
