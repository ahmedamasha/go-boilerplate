-- +goose Up
CREATE TABLE IF NOT EXISTS offers (
    offer_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    segment_id UUID NOT NULL,
    user_id VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    discount_percent DECIMAL(5,2) CHECK (discount_percent >= 0 AND discount_percent <= 100),
    discount_amount DECIMAL(10,2) CHECK (discount_amount >= 0),
    free_shipping BOOLEAN NOT NULL DEFAULT false,
    min_order_value DECIMAL(10,2) CHECK (min_order_value >= 0),
    max_discount DECIMAL(10,2) CHECK (max_discount >= 0),
    coupon_code VARCHAR(50) UNIQUE,
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    usage_limit INTEGER CHECK (usage_limit >= 0),
    usage_count INTEGER NOT NULL DEFAULT 0 CHECK (usage_count >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (segment_id) REFERENCES segments(segment_id) ON DELETE CASCADE,
    CHECK (valid_until > valid_from),
    CHECK (usage_limit IS NULL OR usage_count <= usage_limit)
);

-- Create indexes
CREATE INDEX idx_offers_segment_id ON offers(segment_id);
CREATE INDEX idx_offers_user_id ON offers(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_offers_coupon_code ON offers(coupon_code) WHERE coupon_code IS NOT NULL;
CREATE INDEX idx_offers_valid_from ON offers(valid_from);
CREATE INDEX idx_offers_valid_until ON offers(valid_until);
CREATE INDEX idx_offers_active ON offers(is_active);
CREATE INDEX idx_offers_validity ON offers(valid_from, valid_until, is_active);

-- +goose Down
DROP TABLE IF EXISTS offers; 