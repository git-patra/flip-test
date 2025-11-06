package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/constant"
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"boilerplate-go/internal/pkg/statements/infrastructure/repo/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGetIssuesUsecase(t *testing.T) {
	mockUserRepo := mocks.NewInMemoryRepo(t)
	uc := NewGetIssuesUsecase(mockUserRepo)
	assert.NotNil(t, uc)
}

func TestGetIssuesUsecase_Execute(t *testing.T) {
	type args struct {
		uploadID string
		statuses []string
		page     int
		size     int
	}

	type getByUpload struct {
		returnTxs []entity.Transaction
	}

	tests := []struct {
		name        string
		args        args
		getByUpload *getByUpload
		want        []entity.Transaction
		total       int
	}{
		{
			name: "no status filter and data is empty",
			args: args{uploadID: "no-data", statuses: []string{}, page: 0, size: 10},
			getByUpload: &getByUpload{
				returnTxs: []entity.Transaction{},
			},
			want:  []entity.Transaction{},
			total: 0,
		},
		{
			name: "filter by FAILED status",
			args: args{uploadID: "upload-123", statuses: []string{string(constant.FAILED)}, page: 1, size: 10},
			getByUpload: &getByUpload{
				returnTxs: []entity.Transaction{
					{
						Description: "Test Transaction 1",
						Status:      constant.FAILED,
					},
					{
						Description: "Test Transaction 2",
						Status:      constant.PENDING,
					},
					{
						Description: "Test Transaction 3",
						Status:      constant.SUCCESS,
					},
					{
						Description: "Test Transaction 4",
						Status:      constant.FAILED,
					},
				},
			},
			want: []entity.Transaction{
				{
					Description: "Test Transaction 1",
					Status:      constant.FAILED,
				},
				{
					Description: "Test Transaction 4",
					Status:      constant.FAILED,
				},
			},
			total: 2,
		},
		{
			name: "no status filter with single page",
			args: args{uploadID: "upload-456", statuses: []string{}, page: 1, size: 10},
			getByUpload: &getByUpload{
				returnTxs: []entity.Transaction{
					{
						Description: "Test Transaction 1",
						Status:      constant.FAILED,
					},
					{
						Description: "Test Transaction 2",
						Status:      constant.PENDING,
					},
					{
						Description: "Test Transaction 3",
						Status:      constant.SUCCESS,
					},
				},
			},
			want: []entity.Transaction{
				{
					Description: "Test Transaction 1",
					Status:      constant.FAILED,
				},
				{
					Description: "Test Transaction 2",
					Status:      constant.PENDING,
				},
			},
			total: 2,
		},
		{
			name: "no status filter with multiple pages",
			args: args{uploadID: "upload-456", statuses: []string{}, page: 2, size: 2},
			getByUpload: &getByUpload{
				returnTxs: []entity.Transaction{
					{
						Description: "Test Transaction 1",
						Status:      constant.FAILED,
					},
					{
						Description: "Test Transaction 2",
						Status:      constant.PENDING,
					},
					{
						Description: "Test Transaction 3",
						Status:      constant.SUCCESS,
					},
					{
						Description: "Test Transaction 4",
						Status:      constant.FAILED,
					},
					{
						Description: "Test Transaction 5",
						Status:      constant.PENDING,
					},
				},
			},
			want: []entity.Transaction{
				{
					Description: "Test Transaction 4",
					Status:      constant.FAILED,
				},
				{
					Description: "Test Transaction 5",
					Status:      constant.PENDING,
				},
			},
			total: 4,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := mocks.NewInMemoryRepo(t)
			if tt.getByUpload != nil {
				mockRepo.
					On("GetByUpload", tt.args.uploadID).
					Return(tt.getByUpload.returnTxs)
			}
			uc := &getIssuesUsecase{
				repo: mockRepo,
			}
			got, total := uc.Execute(tt.args.uploadID, tt.args.statuses, tt.args.page, tt.args.size)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.total, total)
		})
	}
}
