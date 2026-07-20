package reviews

import "github.com/gin-gonic/gin"

// RegisterProductReviews mounts product-review routes. The public listing
// endpoint sits at GET /products/:id/reviews while the caller-scoped endpoints
// require auth.
func RegisterProductReviews(rg *gin.RouterGroup, h *Handler, auth gin.HandlerFunc) {
	g := rg.Group("/products/:id/reviews")
	g.GET("", h.List)
	authed := rg.Group("/products/:id/reviews")
	authed.Use(auth)
	authed.GET("/me", h.GetMine)
	authed.PUT("/me", h.UpsertMe)
}
