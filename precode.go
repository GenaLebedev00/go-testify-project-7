package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
    "moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
    countStr := req.URL.Query().Get("count")
    if countStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("count missing"))
        return
    }

    count, err := strconv.Atoi(countStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong count value"))
        return
    }

    city := req.URL.Query().Get("city")

    cafe, ok := cafeList[city]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong city value"))
        return
    }

    if count > len(cafe) {
        count = len(cafe)
    }

    answer := strings.Join(cafe[:count], ",")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    totalCount := 4
    expectedCafes := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
    req := httptest.NewRequest("GET","/cafe",nil)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    // здесь нужно добавить необходимые проверки
    assert.Equal(t, http.StatusOK, responseRecorder.Code)

    assert.Equal(t,expectedCafes, responseRecorder.Body.String())

    cafeList := strings.Split(responseRecorder.Body.String(), ",")
    assert.Len(t,cafeList,totalCount)

    assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
    assert.Equal(t, "wrong city value", responseRecorder.Body.String())

    assert.Equal(t, "count missing", responseRecorder.Body.String())

}
