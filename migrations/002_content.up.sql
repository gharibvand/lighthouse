CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS directors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS content (
    id UUID PRIMARY KEY ,
    type VARCHAR(20) NOT NULL CHECK (type IN ('movie', 'tv_show')),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_year INTEGER,
    duration INTEGER,
    rating DECIMAL(3,2) DEFAULT 0.0,
    poster_url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    trailer_url VARCHAR(500),
    video_path VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS episodes (
    id UUID PRIMARY KEY ,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    season INTEGER NOT NULL,
    episode_num INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INTEGER, -- Minutes
    thumbnail_url VARCHAR(500),
    video_path VARCHAR(500),
    release_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(content_id, season, episode_num)
);

CREATE TABLE IF NOT EXISTS content_genres (
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    genre_id INTEGER NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, genre_id)
);

CREATE TABLE IF NOT EXISTS content_actors (
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    actor_id INTEGER NOT NULL REFERENCES actors(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, actor_id)
);

CREATE TABLE IF NOT EXISTS content_directors (
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    director_id INTEGER NOT NULL REFERENCES directors(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, director_id)
);

CREATE TABLE IF NOT EXISTS subtitles (
    id UUID PRIMARY KEY ,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    language VARCHAR(10) NOT NULL,
    file_path VARCHAR(500) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_content_type ON content(type);
CREATE INDEX IF NOT EXISTS idx_content_is_active ON content(is_active);
CREATE INDEX IF NOT EXISTS idx_content_title ON content(title);
CREATE INDEX IF NOT EXISTS idx_content_rating ON content(rating DESC);
CREATE INDEX IF NOT EXISTS idx_content_release_year ON content(release_year DESC);
CREATE INDEX IF NOT EXISTS idx_content_created_at ON content(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_type_active_rating ON content(type, is_active, rating DESC);

CREATE INDEX IF NOT EXISTS idx_episodes_content_id ON episodes(content_id);
CREATE INDEX IF NOT EXISTS idx_episodes_content_season ON episodes(content_id, season, episode_num);

CREATE INDEX IF NOT EXISTS idx_content_genres_content_id ON content_genres(content_id);
CREATE INDEX IF NOT EXISTS idx_content_genres_genre_id ON content_genres(genre_id);
CREATE INDEX IF NOT EXISTS idx_content_actors_content_id ON content_actors(content_id);
CREATE INDEX IF NOT EXISTS idx_content_actors_actor_id ON content_actors(actor_id);
CREATE INDEX IF NOT EXISTS idx_content_directors_content_id ON content_directors(content_id);
CREATE INDEX IF NOT EXISTS idx_content_directors_director_id ON content_directors(director_id);

CREATE INDEX IF NOT EXISTS idx_subtitles_content_id ON subtitles(content_id);
CREATE INDEX IF NOT EXISTS idx_subtitles_episode_id ON subtitles(episode_id);
CREATE INDEX IF NOT EXISTS idx_subtitles_language ON subtitles(language);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_content_title_trgm ON content USING gin(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_content_description_trgm ON content USING gin(description gin_trgm_ops);
