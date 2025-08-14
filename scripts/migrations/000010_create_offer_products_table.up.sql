-- +goose Up
CREATE TABLE IF NOT EXISTS offer_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (offer_id) REFERENCES offers(offer_id) ON DELETE CASCADE,
    UNIQUE(offer_id, product_id)
);

-- Create indexes
CREATE INDEX idx_offer_products_offer_id ON offer_products(offer_id);
CREATE INDEX idx_offer_products_product_id ON offer_products(product_id);

-- +goose Down
DROP TABLE IF EXISTS offer_products; 