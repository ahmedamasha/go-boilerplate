-- +goose Up
CREATE TABLE IF NOT EXISTS offer_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    order_id VARCHAR(255),
    used_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    FOREIGN KEY (offer_id) REFERENCES offers(offer_id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_offer_usage_offer_id ON offer_usage(offer_id);
CREATE INDEX idx_offer_usage_user_id ON offer_usage(user_id);
CREATE INDEX idx_offer_usage_order_id ON offer_usage(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX idx_offer_usage_used_at ON offer_usage(used_at);

-- +goose Down
DROP TABLE IF EXISTS offer_usage; 