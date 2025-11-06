package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/constant"
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"boilerplate-go/internal/pkg/statements/infrastructure/repo"
	"strings"
)

type getIssuesUsecase struct {
	repo repo.InMemoryRepo
}

func NewGetIssuesUsecase(r repo.InMemoryRepo) GetIssuesUsecase {
	return &getIssuesUsecase{repo: r}
}

func (u *getIssuesUsecase) Execute(uploadID string, statuses []string, page, size int) ([]entity.Transaction, int) {
	txs := u.repo.GetByUpload(uploadID)

	if len(txs) == 0 {
		return []entity.Transaction{}, 0
	}

	filter := map[string]struct{}{}
	for _, s := range statuses {
		filter[strings.ToUpper(s)] = struct{}{}
	}

	var issues []entity.Transaction
	for _, t := range txs {
		if len(filter) > 0 {
			if _, ok := filter[string(t.Status)]; !ok {
				continue
			}
		} else if t.Status != constant.FAILED && t.Status != constant.PENDING {
			continue
		}
		issues = append(issues, t)
	}

	total := len(issues)
	start := (page - 1) * size
	if start > total {
		return []entity.Transaction{}, total
	}
	end := start + size
	if end > total {
		end = total
	}

	return issues[start:end], total
}
