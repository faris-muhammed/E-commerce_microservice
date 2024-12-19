package handlers

import (
	"context"

	"github.com/faris-muhammed/e-commerce/seller-service/service"
	sellerpb "github.com/faris-muhammed/e-protofiles/seller"
)

type CategoryHandler struct {
	sellerpb.UnimplementedCategoryServiceServer
	service service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: svc}
}

func (h *CategoryHandler) AddCategory(ctx context.Context, req *sellerpb.AddCategoryRequest) (*sellerpb.CategoryResponse, error) {
	return h.service.AddCategory(req)
}

func (h *CategoryHandler) EditCategory(ctx context.Context, req *sellerpb.EditCategoryRequest) (*sellerpb.CategoryResponse, error) {
	return h.service.EditCategory(req)
}

func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *sellerpb.DeleteCategoryRequest) (*sellerpb.CategoryResponse, error) {
	return h.service.DeleteCategory(req)
}

func (h *CategoryHandler) GetCategories(ctx context.Context, req *sellerpb.Empty) (*sellerpb.CategoryListResponse, error) {
	return h.service.GetCategories()
}
