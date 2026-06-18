CREATE DOMAIN email_type AS VARCHAR(255)
    CHECK (VALUE ~ '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email email_type NOT NULL,
    is_active BOOLEAN DEFAULT false,
    email_token VARCHAR(64), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE apartments (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT,
    name TEXT NOT NULL,
    description TEXT,
    capacity SMALLINT,
    price_per_hour INT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (price_per_hour > 0),

    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_apartments_seller ON apartments(seller_id);
CREATE INDEX idx_apartments_active ON apartments(is_active);

CREATE TABLE apartment_images (
    id BIGSERIAL PRIMARY KEY,
    apartment_id BIGINT NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    position INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_apartment_images_apartment ON apartment_images(apartment_id);

CREATE TABLE reviews (
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
CREATE INDEX idx_reviews_apartment ON reviews(apartment_id);

CREATE TABLE booking (
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
CREATE INDEX idx_booking_dates ON booking(time_from, time_to);
CREATE INDEX idx_booking_apartment ON booking(apartment_id);
