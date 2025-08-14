-- +goose Up
CREATE TABLE IF NOT EXISTS offer_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL,
    category VARCHAR(255) NOT NULL,
    FOREIGN KEY (offer_id) REFERENCES offers(offer_id) ON DELETE CASCADE,
    UNIQUE(offer_id, category)
);

-- Create indexes
CREATE INDEX idx_offer_categories_offer_id ON offer_categories(offer_id);
CREATE INDEX idx_offer_categories_category ON offer_categories(category);

-- +goose Down
DROP TABLE IF EXISTS offer_categories; 