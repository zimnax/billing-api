package api

import (
	"billing-api/model"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

func respondWithError(w http.ResponseWriter, e model.PaymentError) {
	respondWithJSON(w, e.Code, map[string]string{"error": e.Message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func isCorrectMoneyFormat(m string) (bool, error) {
	match, err := regexp.MatchString("\\d[.]\\d\\d", m)
	if err != nil {
		return false, err
	}
	return match, nil
}

func getRequestFilter(r *http.Request) (*model.TransactionFilter, *model.PaymentError) {

	f := &model.TransactionFilter{
		CompanyId: r.FormValue("company_id"),
		URL:r.URL.Query(),
	}

	root := false //TODO ROOT STUB

	if !root && f.CompanyId == "" {
		return nil, &model.PermissionRestrictions
	}

	dfrom := r.FormValue("date_from")
	if dfrom != "" {
		uintDateFrom, err := strconv.ParseUint(dfrom, 0, 64)
		if err != nil {
			return nil, &model.InternalError
		}
		f.DateFrom = uintDateFrom
	}

	dto := r.FormValue("date_from")
	if dto != "" {
		uintDateTo, err := strconv.ParseUint(dto, 0, 64)
		if err != nil {
			return nil, &model.InternalError
		}
		f.DateTo = uintDateTo
	}

	page := r.FormValue("page")
	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return nil, &model.InternalError
		}
		f.Page = pageInt
	}

	size := r.FormValue("size")
	if page != "" {
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			return nil, &model.InternalError
		}
		f.Limit = sizeInt
	}

	return f, nil
}
