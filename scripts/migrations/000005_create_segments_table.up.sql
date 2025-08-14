-- +goose Up
CREATE TABLE IF NOT EXISTS segments (
    segment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    segment_name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    rules JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_segments_name ON segments(segment_name);
CREATE INDEX idx_segments_active ON segments(is_active);
CREATE INDEX idx_segments_rules ON segments USING GIN(rules);

-- +goose Down
DROP TABLE IF EXISTS segments; 