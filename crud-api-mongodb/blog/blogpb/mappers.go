package blogpb

import "crud-api-mongodb/blog/model"

func (b *Blog) AsBlogItem() *model.BlogItem {
	blogItem := &model.BlogItem{
		AuthorID: b.GetAuthorId(),
		Content:  b.GetContent(),
		Title:    b.GetTitle(),
	}
	return blogItem
}

func CreateBlogResponseFromBlogItem(blogItem *model.BlogItem) *CreateBlogResponse {
	createBlogResponse := &CreateBlogResponse{Blog: &Blog{
		Id:       blogItem.ID.Hex(),
		AuthorId: blogItem.AuthorID,
		Title:    blogItem.Title,
		Content:  blogItem.Content,
	}}
	return createBlogResponse
}
