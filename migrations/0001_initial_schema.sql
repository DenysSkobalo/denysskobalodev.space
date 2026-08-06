PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    email TEXT UNIQUE NOT NULL,
    role TEXT CHECK(role IN ('admin', 'editor', 'user')) NOT NULL DEFAULT 'user',
    created_at INTEGER NOT NULL DEFAULT (UNIXEPOCH())
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    icon TEXT NOT NULL,
    category TEXT NOT NULL,
    status TEXT CHECK(status IN ('draft', 'published')) NOT NULL DEFAULT 'published',
    tagline TEXT NOT NULL,
    description TEXT NOT NULL,
    tech_stack TEXT NOT NULL,
    demo_url TEXT,
    github_url TEXT,
    created_at INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
    updated_at INTEGER NOT NULL DEFAULT (UNIXEPOCH())
);

CREATE INDEX IF NOT EXISTS idx_projects_slug_status ON projects(slug, status);
CREATE INDEX IF NOT EXISTS idx_projects_created ON projects(created_at DESC) WHERE status = 'published';
