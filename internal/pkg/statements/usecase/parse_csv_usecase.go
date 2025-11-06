package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/response"
	"context"
	"io"
)

type ParseCSVUsecase interface {
	Execute(ctx context.Context, r io.Reader) (result response.StatementResponse, err error)
}
