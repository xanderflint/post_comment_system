package post_comment_system

import (
	"errors"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
)

var db *gorm.DB

type Post struct {
	ID              string `gorm:"primary_key"`
	Title           string
	Content         string
	CommentsEnabled bool
	Comments        []Comment `gorm:"foreignkey:PostID"`
}

type Comment struct {
	ID       string `gorm:"primary_key"`
	PostID   string
	ParentID string
	Content  string
}

func init() {
	var err error
	dbType := os.Getenv("DB_TYPE")

	if dbType == "sqlite" {
		db, err = gorm.Open("sqlite3", ":memory:")
	} else {
		db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	}

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Post{}, &Comment{})
}

func generateID() string {
	return uuid.New().String()
}

func resolvePost(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(string)
	var post Post
	if err := db.Preload("Comments").First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func resolvePosts(p graphql.ResolveParams) (interface{}, error) {
	var posts []Post
	if err := db.Preload("Comments").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func resolveCreatePost(p graphql.ResolveParams) (interface{}, error) {
	post := Post{
		ID:              generateID(),
		Title:           p.Args["title"].(string),
		Content:         p.Args["content"].(string),
		CommentsEnabled: p.Args["commentsEnabled"].(bool),
	}
	if err := db.Create(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func resolveCreateComment(p graphql.ResolveParams) (interface{}, error) {
	comment := Comment{
		ID:      generateID(),
		PostID:  p.Args["postId"].(string),
		Content: p.Args["content"].(string),
	}
	if parentId, ok := p.Args["parentId"].(string); ok {
		comment.ParentID = parentId
	}
	if len(comment.Content) > 2000 {
		return nil, errors.New("comment content exceeds 2000 characters")
	}
	if err := db.Create(&comment).Error; err != nil {
		return nil, err
	}

	// Notify subscribers about the new comment
	notifySubscribers(comment.PostID, comment)

	return comment, nil
}
