--
-- SQLiteStudio v3.4.17 生成的文件，周五 3月 28 16:03:28 2025
--
-- 所用的文本编码：UTF-8
--
PRAGMA foreign_keys = off;

BEGIN TRANSACTION;

-- 表：sessionstore
DROP TABLE IF EXISTS sessionstore;

CREATE TABLE
    IF NOT EXISTS sessionstore (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        created_at DATETIME,
        updated_at DATETIME,
        deleted_at DATETIME,
        name TEXT NOT NULL,
        cmd TEXT NOT NULL,
        args TEXT NOT NULL,
        dir TEXT NOT NULL,
        CONSTRAINT uni_sessionstore_name UNIQUE (name)
    );

INSERT INTO
    sessionstore (
        id,
        created_at,
        updated_at,
        deleted_at,
        name,
        cmd,
        args,
        dir
    )
VALUES
    (
        1,
        '2025-03-22 15:05:36.3167149+08:00',
        '2025-03-22 15:05:36.3167149+08:00',
        NULL,
        'pwsh',
        'pwsh',
        '["-noProfile"]',
        'C:\Users\Public'
    );

-- 索引：idx_sessionstore_dir
DROP INDEX IF EXISTS idx_sessionstore_dir;

CREATE INDEX IF NOT EXISTS idx_sessionstore_dir ON sessionstore (dir);

-- 索引：idx_sessionstore_args
DROP INDEX IF EXISTS idx_sessionstore_args;

CREATE INDEX IF NOT EXISTS idx_sessionstore_args ON sessionstore (args);

-- 索引：idx_sessionstore_cmd
DROP INDEX IF EXISTS idx_sessionstore_cmd;

CREATE INDEX IF NOT EXISTS idx_sessionstore_cmd ON sessionstore (cmd);

-- 索引：idx_sessionstore_name
DROP INDEX IF EXISTS idx_sessionstore_name;

CREATE INDEX IF NOT EXISTS idx_sessionstore_name ON sessionstore (name);

-- 索引：idx_sessionstore_deleted_at
DROP INDEX IF EXISTS idx_sessionstore_deleted_at;

CREATE INDEX IF NOT EXISTS idx_sessionstore_deleted_at ON sessionstore (deleted_at);

COMMIT TRANSACTION;

PRAGMA foreign_keys = on;