-- ANR Platform Database Initialization Migration
-- Version: 000001
-- Description: Base schema for Core Platform (Apps, Users, Devices, Remote Config)

-- 1. Applications Registry
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    app_id VARCHAR(64) UNIQUE NOT NULL, -- Format: anr-NNN-slug (e.g. anr-001-wallpaper)
    name VARCHAR(128) NOT NULL,
    description TEXT,
    status VARCHAR(32) NOT NULL DEFAULT 'active', -- active, suspended, deprecated
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_applications_app_id ON applications(app_id);

-- 2. Users (First-party identity)
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'user', -- user, admin, developer
    status VARCHAR(32) NOT NULL DEFAULT 'active', -- active, inactive, banned
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 3. Devices
CREATE TABLE IF NOT EXISTS devices (
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(128) NOT NULL, -- Hardware UUID or Client generated UUID
    app_id VARCHAR(64) NOT NULL REFERENCES applications(app_id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    platform VARCHAR(32) NOT NULL, -- ios, android, web
    os_version VARCHAR(32),
    app_version VARCHAR(32),
    push_token TEXT,
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_device_app UNIQUE (device_id, app_id)
);

CREATE INDEX IF NOT EXISTS idx_devices_app_id ON devices(app_id);
CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);


