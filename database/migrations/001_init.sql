-- Migration: 001_init
-- Description: Initial database setup for cats, types, and masters

-- Masters table
CREATE TABLE IF NOT EXISTS masters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    place TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Types table
CREATE TABLE IF NOT EXISTS types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Cats table
CREATE TABLE IF NOT EXISTS cats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cat_type_id INTEGER NOT NULL,
    master_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cat_type_id) REFERENCES types (id),
    FOREIGN KEY (master_id) REFERENCES masters (id)
);

-- Insert sample data
INSERT OR IGNORE INTO types (name) VALUES 
    ('Домашняя'),
    ('Персидская'),
    ('Сиамская'),
    ('Мейн-кун');

INSERT OR IGNORE INTO masters (first_name, last_name, place) VALUES 
    ('Галсанов', 'Солбон', 'Иркутск'),
    ('Громова', 'Альбина', 'Москва'),
    ('Серебренникова', 'Арина', 'Чита');

INSERT OR IGNORE INTO cats (name, cat_type_id, master_id) VALUES 
    ('Дымок', 1, 1),
    ('Изюм', 2, 2),
    ('Барсик', 3, 3);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_cats_type_id ON cats(cat_type_id);
CREATE INDEX IF NOT EXISTS idx_cats_master_id ON cats(master_id);
CREATE INDEX IF NOT EXISTS idx_masters_name ON masters(first_name, last_name);