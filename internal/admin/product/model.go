package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Description     string    `json:"description"`
	NutritionalInfo string    `json:"nutritional_info"`
	DietaryLabels   string    `json:"dietary_labels"`
	Allergens       string    `json:"allergens"`
	Quantity        int32     `json:"quantity"`
	PriceCents      int64     `json:"price_cents"`
	Currency        string    `json:"currency"`
	Image_URL       string    `json:"image_url"`
	Image_ALT       string    `json:"image_alt"`
}

type AddProductInput struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	NutritionalInfo string    `json:"nutritional_info"`
	DietaryLabels   string    `json:"dietary_labels"`
	Allergens       string    `json:"allergens"`
	Quantity        int32     `json:"quantity"`
	PriceCents      int64     `json:"price_cents"`
	Currency        string    `json:"currency"`
	Image_ALT       string    `json:"image_alt"`
}