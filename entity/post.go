package entity

// Post represents a blog post
type Post struct {
	ID      int    `bson:"id"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
	Author  string `bson:"author"`
}
