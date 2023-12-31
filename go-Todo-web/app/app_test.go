package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-Todo-web.example/model"
)

func TestTodos(t *testing.T) {
	// test code에서만 이런식으로 반환
	getSessionID = func(r *http.Request) string {
		return "test_session_id"
	}
	os.Remove("./test.db")
	assert := assert.New(t)
	ah := MakeHandler("./test.db") // app handler
	defer ah.Close()
	ts := httptest.NewServer(ah)
	defer ts.Close()
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo item"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	var todo model.Todo

	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo item")

	id1 := todo.ID

	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo item2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo item2")

	id2 := todo.ID

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var todoList []*model.Todo

	err = json.NewDecoder(resp.Body).Decode(&todoList)
	assert.NoError(err)
	assert.Equal(2, len(todoList))
	assert.Equal(id1, todoList[0].ID)
	assert.Equal(id2, todoList[1].ID)
	for _, t := range todoList {
		if t.ID == id1 {
			assert.Equal(t.Name, "Test todo item")
		} else if t.ID == id2 {
			assert.Equal(t.Name, "Test todo item2")
		}
	}

	resp, err = http.Get(ts.URL + "/todo-complete/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	var success Success
	err = json.NewDecoder(resp.Body).Decode(&success)
	assert.NoError(err)
	assert.Equal(success.Success, true)

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&todoList)
	assert.NoError(err)
	assert.Equal(2, len(todoList))
	assert.Equal(id1, todoList[0].ID)
	assert.Equal(id2, todoList[1].ID)
	for _, t := range todoList {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&todoList)
	assert.NoError(err)
	assert.Equal(1, len(todoList))
	assert.Equal(id2, todoList[0].ID)
}
