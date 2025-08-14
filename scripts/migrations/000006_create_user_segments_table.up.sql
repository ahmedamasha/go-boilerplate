-- +goose Up
CREATE TABLE IF NOT EXISTS user_segments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    segment_id UUID NOT NULL,
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true,
    FOREIGN KEY (segment_id) REFERENCES segments(segment_id) ON DELETE CASCADE,
    UNIQUE(user_id, segment_id)
);

-- Create indexes
CREATE INDEX idx_user_segments_user_id ON user_segments(user_id);
CREATE INDEX idx_user_segments_segment_id ON user_segments(segment_id);
CREATE INDEX idx_user_segments_active ON user_segments(is_active);
CREATE INDEX idx_user_segments_assigned_at ON user_segments(assigned_at);

-- +goose Down
DROP TABLE IF EXISTS user_segments; 