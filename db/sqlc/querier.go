// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package db

import (
	"context"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateSnippet(ctx context.Context, arg CreateSnippetParams) (Snippet, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccount(ctx context.Context, id int32) error
	DeleteSnippet(ctx context.Context, id int32) error
	GetAccount(ctx context.Context, id int32) (Account, error)
	GetSnippet(ctx context.Context, id int32) (Snippet, error)
	GetUser(ctx context.Context, name string) (User, error)
	ListSnippets(ctx context.Context, arg ListSnippetsParams) ([]Snippet, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
