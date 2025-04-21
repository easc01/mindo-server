-- Ensure UUID extension is enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Playlist Table
CREATE TABLE "playlist" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "name" VARCHAR(255),
    "description" TEXT,
    "thumbnail_url" TEXT,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- user saved playlist
CREATE TABLE "user_playlist" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "playlist_id" uuid NOT NULL,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Topic Table
CREATE TABLE "topic" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "number" int,
    "name" VARCHAR(255),
    "playlist_id" uuid NOT NULL,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Study Material Table
CREATE TABLE "study_material" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "topic_id" uuid NOT NULL UNIQUE,
    "title" VARCHAR(255),
    "content" TEXT,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Saved Study Material Table
CREATE TABLE "saved_study_material" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "study_material_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- YouTube Video Table
CREATE TABLE "youtube_video" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "topic_id" uuid NOT NULL,
    "title" VARCHAR(255),
    "video_date" timestamp,
    "video_views" VARCHAR(255),
    "thumbnail_url" TEXT,
    "expiry_at" timestamp,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Watched Video Table
CREATE TABLE "watched_video" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "youtube_video_id" uuid NOT NULL,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Watched Playlist Table
CREATE TABLE "watched_playlist" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "playlist_id" uuid NOT NULL,
    "progress" int,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Quiz Table
CREATE TABLE "quiz" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "name" VARCHAR(255),
    "thumbnail_url" TEXT,
    "play_count" int,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Quiz Question Table
CREATE TABLE "quiz_question" (
  "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  "quiz_id" uuid NOT NULL,
  "question" TEXT,
  "options" TEXT[],  -- Arrays of strings for options
  "correct_option" int,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_by" uuid
);

-- Quiz Result Table
CREATE TABLE "quiz_result" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "quiz_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "score" decimal,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Quiz Result Question Table
CREATE TABLE "quiz_result_question" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "quiz_result_id" uuid NOT NULL,
    "quiz_question_id" uuid NOT NULL,
    "attempted_option" int,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Community Table
CREATE TABLE "community" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "title" VARCHAR(255),
    "about" TEXT,
    "thumbnail_url" TEXT,
    "logo_url" TEXT,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Message Table
CREATE TABLE "message" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "community_id" uuid NOT NULL,
    "content" TEXT,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Joined Community Table
CREATE TABLE "joined_community" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "user_id" uuid NOT NULL,
    "community_id" uuid NOT NULL,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- User Table (with Enum)
CREATE TYPE user_type AS ENUM ('app_user', 'admin_user');

CREATE TABLE "user" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "user_type" user_type,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- User Token Table
CREATE TABLE "user_token" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "user_id" uuid NOT NULL UNIQUE,
    "refresh_token" TEXT NOT NULL,
    "expires_at" TIMESTAMP NOT NULL,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- App User Interest Table
CREATE TABLE "app_user_interest" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "app_user_id" uuid NOT NULL,
    "interest_id" uuid,
    "name" VARCHAR(255),
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- App User Table
CREATE TABLE "app_user" (
    "user_id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "oauth_client_id" VARCHAR(255) UNIQUE,
    "username" VARCHAR(255) UNIQUE,
    "profile_picture_url" TEXT,
    "bio" TEXT,
    "name" VARCHAR(255),
    "mobile" VARCHAR(50),
    "email" VARCHAR(255),
    "password_hash" TEXT,
    "last_login_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Admin User Table
CREATE TABLE "admin_user" (
    "user_id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "name" VARCHAR(255),
    "email" VARCHAR(255) UNIQUE,
    "password_hash" TEXT,
    "last_login_at" timestamp,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

CREATE TABLE "interest" (
    "id" uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    "name" VARCHAR(255),
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_by" uuid
);

-- Add Foreign Keys
ALTER TABLE "user_playlist"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_playlist"
ADD FOREIGN KEY ("playlist_id") REFERENCES "playlist" ("id");

ALTER TABLE "user_token"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "topic"
ADD FOREIGN KEY ("playlist_id") REFERENCES "playlist" ("id");

ALTER TABLE "study_material"
ADD FOREIGN KEY ("topic_id") REFERENCES "topic" ("id");

ALTER TABLE "saved_study_material"
ADD FOREIGN KEY ("study_material_id") REFERENCES "study_material" ("id");

ALTER TABLE "saved_study_material"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "youtube_video"
ADD FOREIGN KEY ("topic_id") REFERENCES "topic" ("id");

ALTER TABLE "watched_video"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "watched_video"
ADD FOREIGN KEY ("youtube_video_id") REFERENCES "youtube_video" ("id");

ALTER TABLE "watched_playlist"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "watched_playlist"
ADD FOREIGN KEY ("playlist_id") REFERENCES "playlist" ("id");

ALTER TABLE "quiz_question"
ADD FOREIGN KEY ("quiz_id") REFERENCES "quiz" ("id");

ALTER TABLE "quiz_result"
ADD FOREIGN KEY ("quiz_id") REFERENCES "quiz" ("id");

ALTER TABLE "quiz_result"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "quiz_result_question"
ADD FOREIGN KEY ("quiz_result_id") REFERENCES "quiz_result" ("id");

ALTER TABLE "quiz_result_question"
ADD FOREIGN KEY ("quiz_question_id") REFERENCES "quiz_question" ("id");

ALTER TABLE "message"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "message"
ADD FOREIGN KEY ("community_id") REFERENCES "community" ("id");

ALTER TABLE "joined_community"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "joined_community"
ADD FOREIGN KEY ("community_id") REFERENCES "community" ("id");

ALTER TABLE "app_user_interest"
ADD FOREIGN KEY ("interest_id") REFERENCES "interest" ("id");

ALTER TABLE "app_user_interest"
ADD FOREIGN KEY ("app_user_id") REFERENCES "app_user" ("user_id");

ALTER TABLE "app_user"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "admin_user"
ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");