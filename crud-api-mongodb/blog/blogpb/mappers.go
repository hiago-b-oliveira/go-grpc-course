package blogpb

import "crud-api-mongodb/blog/model"

func (b *Blog) AsBlogItem() *model.BlogItem {
	return &model.BlogItem{
		AuthorID: b.GetAuthorId(),
		Content:  b.GetContent(),
		Title:    b.GetTitle(),
	}
}

func CreateBlogFromBlogItem(blogItem *model.BlogItem) *Blog {
	return &Blog{
		Id:       blogItem.ID.Hex(),
		AuthorId: blogItem.AuthorID,
		Title:    blogItem.Title,
		Content:  blogItem.Content,
	}
}
