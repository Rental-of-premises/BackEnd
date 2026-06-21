DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'email_type') THEN
        CREATE DOMAIN email_type AS VARCHAR(255)
            CHECK (VALUE ~ '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email email_type NOT NULL,
    is_active BOOLEAN DEFAULT false,
    email_token VARCHAR(64), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS avatar (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_avatar_user ON avatar(user_id);

    -- <=============================================>

CREATE TABLE IF NOT EXISTS apartments (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT,
    name TEXT NOT NULL,
    description TEXT,
    capacity SMALLINT,
    price_per_hour INT,
    is_active BOOLEAN DEFAULT true,
    address TEXT NOT NULL,
    metro TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (price_per_hour > 0),

    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_apartments_seller ON apartments(seller_id);
CREATE INDEX IF NOT EXISTS idx_apartments_active ON apartments(is_active);

CREATE TABLE IF NOT EXISTS apartment_images (
    id BIGSERIAL PRIMARY KEY,
    apartment_id BIGINT NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    position INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_apartment_images_apartment ON apartment_images(apartment_id);

CREATE TABLE IF NOT EXISTS amenities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    icon VARCHAR(50)
);

INSERT INTO amenities (name, icon) VALUES 
    ('Wi-Fi', 'wifi'),
    ('Кондиционер', 'ac'),
    ('Телевизор', 'tv'),
    ('Кухня', 'kitchen'),
    ('Стиральная машина', 'washing_machine'),
    ('Парковка', 'parking'),
    ('Бассейн', 'pool'),
    ('Бильярд', 'billiards'),
    ('Настольный теннис', 'table_tennis'),
    ('Караоке', 'karaoke'),
    ('Библиотека', 'library'),
    ('Игровая комната', 'game_room'),
    ('Кинотеатр', 'cinema'),
    ('Бар', 'bar'),
    ('Кальянная', 'hookah'),
    ('Спортзал', 'gym')
    ON CONFLICT (name) DO NOTHING;;

CREATE TABLE IF NOT EXISTS apartment_amenities (
apartment_id BIGINT NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
amenity_id BIGINT NOT NULL REFERENCES amenities(id) ON DELETE CASCADE,
PRIMARY KEY (apartment_id, amenity_id)
);
CREATE INDEX IF NOT EXISTS idx_apartment_amenities_apartment ON apartment_amenities(apartment_id);
CREATE INDEX IF NOT EXISTS idx_apartment_amenities_amenity ON apartment_amenities(amenity_id);

    -- <=============================================>

CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    apartment_id BIGINT,
    comment TEXT,
    stars SMALLINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (stars >= 1 AND stars <= 5),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_reviews_apartment ON reviews(apartment_id);

    -- <=============================================>

CREATE TABLE IF NOT EXISTS booking (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    apartment_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'confirmed',
    time_from TIMESTAMPTZ NOT NULL,
    time_to TIMESTAMPTZ NOT NULL,

    CHECK (status IN ('confirmed', 'cancelled', 'completed', 'waiting', 'rejected')),
    CHECK (time_to > time_from),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (apartment_id) REFERENCES apartments(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_booking_dates ON booking(time_from, time_to);
CREATE INDEX IF NOT EXISTS idx_booking_apartment ON booking(apartment_id);
