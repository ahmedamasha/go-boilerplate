-- +goose Up
CREATE TABLE IF NOT EXISTS user_interests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    segment_id UUID NOT NULL,
    product_id VARCHAR(255),
    category VARCHAR(255),
    score DECIMAL(3,2) NOT NULL DEFAULT 0.0 CHECK (score >= 0.0 AND score <= 1.0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (segment_id) REFERENCES segments(segment_id) ON DELETE CASCADE,
    CHECK (product_id IS NOT NULL OR category IS NOT NULL)
);

-- Create indexes
CREATE INDEX idx_user_interests_user_id ON user_interests(user_id);
CREATE INDEX idx_user_interests_segment_id ON user_interests(segment_id);
CREATE INDEX idx_user_interests_product_id ON user_interests(product_id) WHERE product_id IS NOT NULL;
CREATE INDEX idx_user_interests_category ON user_interests(category) WHERE category IS NOT NULL;
CREATE INDEX idx_user_interests_score ON user_interests(score);

-- +goose Down
DROP TABLE IF EXISTS user_interests; 