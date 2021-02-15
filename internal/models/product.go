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
	ImageURL    *string            `json:"imageUrl" bson:"imageUrl,omitempty"`
	Photos      []string           `json:"photos" bson:"photos,omitempty"`
	Quantity    int64              `json:"quantity" bson:"quantity,omitempty" validate:"required"`
	Rating      int                `json:"rating" bson:"rating,omitempty" validate:"required,min=0,max=10"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func (p *Product) GetImage() string {
	var img string
	if p.ImageURL != nil {
		img = *p.ImageURL
	}
	return img
}

// ToProto Convert product to proto
func (p *Product) ToProto() *productsService.Product {
	return &productsService.Product{
		ProductID:   p.ProductID.String(),
		CategoryID:  p.CategoryID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		ImageURL:    p.GetImage(),
		Photos:      p.Photos,
		Quantity:    p.Quantity,
		Rating:      int64(p.Rating),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

// ProductFromProto Get Product from proto
func ProductFromProto(product *productsService.Product) (*Product, error) {
	prodID, err := primitive.ObjectIDFromHex(product.GetCategoryID())
	if err != nil {
		return nil, err
	}
	catID, err := primitive.ObjectIDFromHex(product.GetCategoryID())
	if err != nil {
		return nil, err
	}

	return &Product{
		ProductID:   prodID,
		CategoryID:  catID,
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       product.GetPrice(),
		ImageURL:    &product.ImageURL,
		Photos:      product.GetPhotos(),
		Quantity:    product.GetQuantity(),
		Rating:      int(product.GetRating()),
		CreatedAt:   product.GetCreatedAt().AsTime(),
		UpdatedAt:   product.GetUpdatedAt().AsTime(),
	}, nil
}
