package handler

import (
	"net/http"

	"bookkeeping/internal/model"
	"bookkeeping/pkg/httpx"
)

// registerCategoryRoutes 注册分类相关路由。
func (s *Server) registerCategoryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/categories", s.createCategory)
	mux.HandleFunc("GET /api/categories", s.listCategories)
	mux.HandleFunc("GET /api/categories/{id}", s.getCategory)
	mux.HandleFunc("PUT /api/categories/{id}", s.updateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}", s.deleteCategory)
}

type categoryRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCategory(model.Category{Name: req.Name, Type: req.Type})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

// listCategories 分类列表：GET /api/categories?type=income|expense
func (s *Server) listCategories(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListCategories(r.URL.Query().Get("type"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetCategory(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCategory(r.PathValue("id"), model.Category{Name: req.Name, Type: req.Type})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteCategory(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
