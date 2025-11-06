package usecase

import (
	"boilerplate-go/internal/pkg/statements/domain/entity"
	"boilerplate-go/internal/pkg/statements/infrastructure/repo/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGetBalanceUsecase(t *testing.T) {
	mockUserRepo := mocks.NewInMemoryRepo(t)
	uc := NewGetBalanceUsecase(mockUserRepo)
	assert.NotNil(t, uc)
}

func TestGetBalanceUsecase_Execute(t *testing.T) {
	type args struct {
		uploadID string
	}

	type getByUpload struct {
		returnTxs []entity.Transaction
	}

	tests := []struct {
		name        string
		args        args
		getByUpload *getByUpload
		wantBalance int64
	}{
		{
			name: "success",
			args: args{
				uploadID: "upload-123",
			},
			getByUpload: &getByUpload{
				returnTxs: []entity.Transaction{
					{Type: "CREDIT", Amount: 1000, Status: "SUCCESS"},
					{Type: "DEBIT", Amount: 500, Status: "SUCCESS"},
					{Type: "CREDIT", Amount: 200, Status: "FAILED"},
				},
			},
			wantBalance: 500,
		},
		{
			name: "no transactions",
			args: args{
				uploadID: "upload-456",
			},
			getByUpload: &getByUpload{
				returnTxs: []entity.Transaction{},
			},
			wantBalance: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewInMemoryRepo(t)
			if tt.getByUpload != nil {
				mockRepo.
					On("GetByUpload", tt.args.uploadID).
					Return(tt.getByUpload.returnTxs)
			}
			uc := &getBalanceUsecase{
				repo: mockRepo,
			}
			gotBalance := uc.Execute(tt.args.uploadID)
			if gotBalance != tt.wantBalance {
				t.Errorf("GetBalanceUsecase.Execute() = %v, want %v", gotBalance, tt.wantBalance)
			}
		})
	}
}
