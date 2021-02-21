package v1

// MapRoutes products routes
func (p *productHandlers) MapRoutes() {
	p.group.POST("", p.CreateProduct())
	p.group.PUT("/:product_id", p.CreateProduct())
	p.group.GET("/:product_id", p.GetByIDProduct())
	p.group.GET("/search", p.SearchProduct())
}
