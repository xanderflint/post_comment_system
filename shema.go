package post_comment_system

import (
	"github.com/graphql-go/graphql"
)

var commentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Comment",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"postId": &graphql.Field{
			Type: graphql.String,
		},
		"parentId": &graphql.Field{
			Type: graphql.String,
		},
		"content": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"content": &graphql.Field{
			Type: graphql.String,
		},
		"comments": &graphql.Field{
			Type: graphql.NewList(commentType),
		},
		"commentsEnabled": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"post": &graphql.Field{
			Type:        postType,
			Description: "Get post by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: resolvePost,
		},
		"posts": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "Get list of posts",
			Resolve:     resolvePosts,
		},
	},
})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createPost": &graphql.Field{
			Type:        postType,
			Description: "Create a new post",
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"content": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"commentsEnabled": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Boolean),
				},
			},
			Resolve: resolveCreatePost,
		},
		"createComment": &graphql.Field{
			Type:        commentType,
			Description: "Create a new comment",
			Args: graphql.FieldConfigArgument{
				"postId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"parentId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"content": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: resolveCreateComment,
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    queryType,
	Mutation: mutationType,
})
