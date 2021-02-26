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
	CategoryID  primitive.ObjectID `json:"categoryId,omitempty" bson:"categoryId,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	ImageURL    *string            `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	Photos      []string           `json:"photos,omitempty" bson:"photos,omitempty"`
	Quantity    int64              `json:"quantity,omitempty" bson:"quantity,omitempty" validate:"required"`
	Rating      int                `json:"rating,omitempty" bson:"rating,omitempty" validate:"required,min=0,max=10"`
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

// ProductsList All Products response with pagination
type ProductsList struct {
	TotalCount int64      `json:"totalCount"`
	TotalPages int64      `json:"totalPages"`
	Page       int64      `json:"page"`
	Size       int64      `json:"size"`
	HasMore    bool       `json:"hasMore"`
	Products   []*Product `json:"products"`
}

// ToProtoList convert products list to proto
func (p *ProductsList) ToProtoList() []*productsService.Product {
	productsList := make([]*productsService.Product, 0, len(p.Products))
	for _, product := range p.Products {
		productsList = append(productsList, product.ToProto())
	}
	return productsList
}
