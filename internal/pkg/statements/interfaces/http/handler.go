package httpiface

import (
	"boilerplate-go/internal/delivery/rest/response"
	"boilerplate-go/internal/pkg/statements/usecase"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	parseCSV   *usecase.ParseCSVUsecase
	getBalance *usecase.GetBalanceUsecase
	getIssues  *usecase.GetIssuesUsecase
}

func NewHandler(p *usecase.ParseCSVUsecase, g *usecase.GetBalanceUsecase, i *usecase.GetIssuesUsecase) *Handler {
	return &Handler{parseCSV: p, getBalance: g, getIssues: i}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	var (
		apiResponse *response.ApiResponse
		ctx         = r.Context()
	)

	file, _, err := r.FormFile("file")
	if err != nil {
		logrus.Warn(err.Error())
		apiResponse = response.BuildErrorResponse(response.ValidationError)
		response.JSON(w, apiResponse.StatusCode, apiResponse)
		return
	}
	defer file.Close()

	result, err := h.parseCSV.Execute(ctx, file)
	if err != nil {
		logrus.Warn(err.Error())
		apiResponse = response.BuildErrorResponse(response.ValidationError)
		response.JSON(w, apiResponse.StatusCode, apiResponse)
		return
	}

	apiResponse = response.BuildSuccessResponseWithData(response.Created, result)
	response.JSON(w, apiResponse.StatusCode, apiResponse)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("upload_id")
	result := h.getBalance.Execute(id)

	apiResponse := response.BuildSuccessResponseWithData(response.Ok, result)
	response.JSON(w, apiResponse.StatusCode, apiResponse)
}

func (h *Handler) GetIssues(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("upload_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	statuses := []string{}
	if s := r.URL.Query().Get("status"); s != "" {
		statuses = strings.Split(s, ",")
	}

	items, total := h.getIssues.Execute(id, statuses, page, size)

	metadata := response.SetPagination(page, size, total)
	apiResponse := response.BuildSuccessResponseWithDataAndMetaData(response.Ok, items, metadata)
	response.JSON(w, apiResponse.StatusCode, apiResponse)
}
