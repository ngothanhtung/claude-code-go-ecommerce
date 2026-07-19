package catalog

import "github.com/gin-gonic/gin"

func Register(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/catalog")
	g.GET("/categories", h.ListCategories)
	g.GET("/products", h.ListProducts)
	g.GET("/products/:id", h.GetProduct)
	g.GET("/products/:id/related", h.GetRelated)
	g.GET("/promos", h.ListPromos)
}
