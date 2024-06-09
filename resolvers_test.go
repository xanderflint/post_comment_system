package post_comment_system

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initTestDB() *gorm.DB {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.AutoMigrate(&Post{}, &Comment{})
	return db
}

func TestCreatePost(t *testing.T) {
	// Инициализация in-memory базы данных
	db = initTestDB()

	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			"title":           "Test Post",
			"content":         "This is a test post",
			"commentsEnabled": true,
		},
	}

	result, err := resolveCreatePost(params)
	assert.Nil(t, err)
	post := result.(Post)

	var savedPost Post
	db.First(&savedPost, "id = ?", post.ID)

	assert.Equal(t, post.ID, savedPost.ID)
	assert.Equal(t, post.Title, savedPost.Title)
	assert.Equal(t, post.Content, savedPost.Content)
	assert.Equal(t, post.CommentsEnabled, savedPost.CommentsEnabled)
}

func TestCreateComment(t *testing.T) {
	db = initTestDB()

	post := Post{
		ID:              generateID(),
		Title:           "Test Post",
		Content:         "This is a test post",
		CommentsEnabled: true,
	}
	db.Create(&post)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			"postId":  post.ID,
			"content": "This is a test comment",
		},
	}

	result, err := resolveCreateComment(params)
	assert.Nil(t, err)
	comment := result.(Comment)

	var savedComment Comment
	db.First(&savedComment, "id = ?", comment.ID)

	assert.Equal(t, comment.ID, savedComment.ID)
	assert.Equal(t, comment.PostID, savedComment.PostID)
	assert.Equal(t, comment.Content, savedComment.Content)
}

func TestResolvePosts(t *testing.T) {
	db = initTestDB()

	post1 := Post{
		ID:              generateID(),
		Title:           "Test Post 1",
		Content:         "This is a test post 1",
		CommentsEnabled: true,
	}
	db.Create(&post1)

	post2 := Post{
		ID:              generateID(),
		Title:           "Test Post 2",
		Content:         "This is a test post 2",
		CommentsEnabled: true,
	}
	db.Create(&post2)

	params := graphql.ResolveParams{}
	result, err := resolvePosts(params)
	assert.Nil(t, err)
	posts := result.([]Post)

	assert.Equal(t, 2, len(posts))
}

func TestResolvePost(t *testing.T) {
	db = initTestDB()

	post := Post{
		ID:              generateID(),
		Title:           "Test Post",
		Content:         "This is a test post",
		CommentsEnabled: true,
	}
	db.Create(&post)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			"id": post.ID,
		},
	}

	result, err := resolvePost(params)
	assert.Nil(t, err)
	resolvedPost := result.(Post)

	assert.Equal(t, post.ID, resolvedPost.ID)
	assert.Equal(t, post.Title, resolvedPost.Title)
	assert.Equal(t, post.Content, resolvedPost.Content)
	assert.Equal(t, post.CommentsEnabled, resolvedPost.CommentsEnabled)
}
