package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	productsService "github.com/AleksK1NG/products-microservice/proto/product"
)

// Product models
type Product struct {
	ProductID   primitive.ObjectID `json:"productId" bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID `json:"CategoryId" bson:"CategoryId,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string             `json:"description" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64            `json:"price" bson:"price,omitempty" validate:"required"`
	ImageURL    string             `json:"imageUrl" bson:"imageUrl,omitempty"`
	Photos      []string           `json:"photos" bson:"photos,omitempty"`
	Quantity    int64              `json:"quantity" bson:"quantity,omitempty" validate:"required"`
	Rating      int                `json:"rating" bson:"rating,omitempty" validate:"required,min=0,max=10"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
}

// ToProto Convert product to proto
func (p *Product) ToProto() *productsService.Product {
	return &productsService.Product{
		ProductID:   p.ProductID.String(),
		CategoryID:  p.CategoryID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		ImageURL:    p.ImageURL,
		Photos:      p.Photos,
		Quantity:    p.Quantity,
		Rating:      int64(p.Rating),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}
