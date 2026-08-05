PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    icon TEXT NOT NULL DEFAULT '🚀',
    category TEXT NOT NULL,
    status TEXT CHECK(status IN ('Live SaaS', 'Open Source', 'In Development')) NOT NULL,
    tagline TEXT NOT NULL,
    description TEXT NOT NULL,
    tech_stack TEXT NOT NULL,
    demo_url TEXT DEFAULT '',
    github_url TEXT DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_slug ON projects(slug);
