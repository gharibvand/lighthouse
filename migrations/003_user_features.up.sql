CREATE TABLE IF NOT EXISTS watch_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    progress INTEGER DEFAULT 0,
    duration INTEGER DEFAULT 0,
    watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS watchlist (
    id UUID PRIMARY KEY ,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, profile_id, content_id)
);

CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY ,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, profile_id, content_id)
);

CREATE INDEX IF NOT EXISTS idx_watch_history_user_id ON watch_history(user_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_profile_id ON watch_history(profile_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_content_id ON watch_history(content_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_user_content ON watch_history(user_id, content_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_profile_content ON watch_history(profile_id, content_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_watched_at ON watch_history(watched_at DESC);

CREATE INDEX IF NOT EXISTS idx_watchlist_user_id ON watchlist(user_id);
CREATE INDEX IF NOT EXISTS idx_watchlist_profile_id ON watchlist(profile_id);
CREATE INDEX IF NOT EXISTS idx_watchlist_user_profile ON watchlist(user_id, profile_id);
CREATE INDEX IF NOT EXISTS idx_watchlist_user_content ON watchlist(user_id, content_id);
CREATE INDEX IF NOT EXISTS idx_watchlist_added_at ON watchlist(added_at DESC);

CREATE INDEX IF NOT EXISTS idx_ratings_content_id ON ratings(content_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user_id ON ratings(user_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user_content ON ratings(user_id, content_id);
CREATE INDEX IF NOT EXISTS idx_ratings_content_created ON ratings(content_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ratings_user_created ON ratings(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ratings_rating ON ratings(rating DESC);
