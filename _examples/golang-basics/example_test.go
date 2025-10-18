package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	client ExampleClient
)

// func TestMain()

func init() {
	go func() {
		startServer()
	}()

	client = NewExampleClient("http://0.0.0.0:4242", &http.Client{
		Timeout: time.Duration(10 * time.Second),
	})
	time.Sleep(time.Millisecond * 500)

}

func TestPing(t *testing.T) {
	err := client.Ping(context.Background())
	assert.NoError(t, err)
}

func TestStatus(t *testing.T) {
	resp, err := client.Status(context.Background())
	assert.Equal(t, true, resp)
	assert.NoError(t, err)
}

func TestVersion(t *testing.T) {
	version, err := client.Version(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, version.ClientGenVersion)
	assert.NotNil(t, version.ServerGenVersion)
}

func TestGetUser(t *testing.T) {
	{
		arg1 := map[string]string{"a": "1"}
		user, err := client.GetUser(context.Background(), arg1, 12)
		assert.Equal(t, &User{ID: 12, Username: "hihi", Nicknames: []Nickname{}}, user)
		assert.NoError(t, err)
	}

	{ // userID == 911, expect not found err
		user, err := client.GetUser(context.Background(), nil, 911)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Contains(t, err.Error(), "not found")

		rpcErr, ok := err.(WebRPCError)
		assert.True(t, ok)
		assert.Equal(t, rpcErr.HTTPStatus, 400)
		assert.Contains(t, rpcErr.Unwrap().Error(), "911")
	}

	{ // userID == 31337, expect unauthorized
		user, err := client.GetUser(context.Background(), nil, 31337)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUnauthorized)
	}

	{ // userID == 666, expect panic
		user, err := client.GetUser(context.Background(), nil, 666)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrWebrpcServerPanic)

		rpcErr, ok := err.(WebRPCError)
		assert.True(t, ok)
		assert.Equal(t, rpcErr.HTTPStatus, 500)
		assert.Contains(t, rpcErr.Unwrap().Error(), "oh no")
	}

	{
		name, user, err := client.FindUser(context.Background(), &SearchFilter{Q: "joe"})
		assert.Equal(t, "joe", name)
		assert.Equal(t, &User{ID: 123, Username: "joe", Nicknames: []Nickname{}}, user)
		assert.NoError(t, err)
	}

	{
		articleResp, err := client.GetArticle(context.Background(), GetArticleRequest{ArticleID: 42})
		assert.NoError(t, err)
		expectedContent := "This is the content of the article."
		assert.Equal(t, &GetArticleResponse{
			Title:   "Article 42",
			Content: &expectedContent,
		}, articleResp)
	}

	{ // streaming
		stream, err := client.StreamNewArticles(context.Background())
		assert.NoError(t, err)

		var articles []*GetArticleResponse
		for {
			article, err := stream.Read()
			if err != nil {
				if errors.Is(err, ErrWebrpcStreamFinished) {
					break
				} else {
					t.Fatal(err)
				}
			}
			articles = append(articles, article)
		}
		require.Len(t, articles, 4)
		require.True(t, articles[0].Title == "Article 0")
	}
}
