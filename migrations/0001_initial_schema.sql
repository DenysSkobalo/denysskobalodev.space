-- D1 Database Schema Optimization for Read-Heavy Edge Operations
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    role TEXT CHECK(role IN ('admin', 'editor', 'user')) NOT NULL DEFAULT 'user',
    created_at INTEGER NOT NULL DEFAULT (UNIXEPOCH())
);

-- DevLog / Posts Table (Read-Optimized)
CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    status TEXT CHECK(status IN ('draft', 'published')) NOT NULL DEFAULT 'draft',
    views_count INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    updated_at INTEGER NOT NULL DEFAULT (UNIXEPOCH())
);

-- Covered Indexes for Fast Read Operations
CREATE INDEX IF NOT EXISTS idx_posts_slug_status ON posts(slug, status);
CREATE INDEX IF NOT EXISTS idx_posts_created_published ON posts(created_at DESC) WHERE status = 'published';

-- Read-Optimized View for Frontend Listing
CREATE VIEW IF NOT EXISTS v_published_posts AS
SELECT id, slug, title, views_count, created_at
FROM posts
WHERE status = 'published'
ORDER BY created_at DESC;
