-- =============================================
-- 1. ДОМЕНЫ
-- =============================================
CREATE DOMAIN email_type AS VARCHAR(255)
    CHECK (VALUE ~ '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- =============================================
-- 2. ТАБЛИЦА ПОЛЬЗОВАТЕЛЕЙ
-- =============================================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email email_type NOT NULL,
    is_active BOOLEAN DEFAULT false,
    email_token VARCHAR(64), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================
-- 3. ТАБЛИЦА ПОМЕЩЕНИЙ (С МЕТРО И АДРЕСОМ)
-- =============================================
CREATE TABLE apartments (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT,
    name TEXT NOT NULL,
    description TEXT,
    capacity SMALLINT,
    price_per_hour INT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- ✅ НОВЫЕ КОЛОНКИ
    metro TEXT,
    address TEXT,
    
    -- ✅ КОЛОНКА ДЛЯ УДОБСТВ (JSONB)
    amenities JSONB DEFAULT '[]'::jsonb,

    CHECK (price_per_hour > 0),

    FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для apartments
CREATE INDEX idx_apartments_seller ON apartments(seller_id);
CREATE INDEX idx_apartments_active ON apartments(is_active);
CREATE INDEX idx_apartments_metro ON apartments(metro);
CREATE INDEX idx_apartments_address ON apartments(address);

-- =============================================
-- 4. ТАБЛИЦА ИЗОБРАЖЕНИЙ
-- =============================================
CREATE TABLE apartment_images (
    id BIGSERIAL PRIMARY KEY,
    apartment_id BIGINT NOT NULL REFERENCES apartments(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    position INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_apartment_images_apartment ON apartment_images(apartment_id);

-- =============================================
-- 5. ТАБЛИЦА ОТЗЫВОВ
-- =============================================
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

-- =============================================
-- 6. ТАБЛИЦА БРОНИРОВАНИЙ
-- =============================================
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
CREATE INDEX idx_booking_status ON booking(status);
CREATE INDEX idx_booking_user ON booking(user_id);

-- =============================================
-- 7. ДОПОЛНИТЕЛЬНЫЕ ИНДЕКСЫ ДЛЯ ПРОИЗВОДИТЕЛЬНОСТИ
-- =============================================

-- Составной индекс для быстрого поиска активных помещений по цене
CREATE INDEX idx_apartments_active_price ON apartments(is_active, price_per_hour);

-- Индекс для поиска по названию (для поиска)
CREATE INDEX idx_apartments_name_trgm ON apartments USING gin (name gin_trgm_ops);

-- Индекс для поиска по метро
CREATE INDEX idx_apartments_metro_trgm ON apartments USING gin (metro gin_trgm_ops);

-- =============================================
-- 8. ФУНКЦИЯ ДЛЯ АВТОМАТИЧЕСКОГО ЗАВЕРШЕНИЯ БРОНИРОВАНИЙ
-- =============================================

-- Создаем функцию для автоматического завершения прошедших бронирований
CREATE OR REPLACE FUNCTION complete_expired_bookings()
RETURNS void AS $$
BEGIN
    UPDATE booking 
    SET status = 'completed' 
    WHERE status = 'confirmed' 
    AND time_to < NOW();
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- 9. ТРИГГЕР ДЛЯ АВТОМАТИЧЕСКОГО ОБНОВЛЕНИЯ СТАТУСА БРОНИРОВАНИЙ
-- =============================================

-- Создаем триггер, который будет вызываться при вставке или обновлении бронирований
CREATE OR REPLACE FUNCTION check_booking_overlap()
RETURNS TRIGGER AS $$
DECLARE
    overlap_count INT;
BEGIN
    -- Проверяем пересечение с другими бронированиями
    SELECT COUNT(*) INTO overlap_count
    FROM booking
    WHERE apartment_id = NEW.apartment_id
    AND status NOT IN ('cancelled', 'rejected')
    AND id != NEW.id
    AND (
        (time_from <= NEW.time_from AND time_to > NEW.time_from) OR
        (time_from < NEW.time_to AND time_to >= NEW.time_to) OR
        (time_from >= NEW.time_from AND time_to <= NEW.time_to)
    );
    
    IF overlap_count > 0 THEN
        RAISE EXCEPTION 'Бронирование пересекается с существующим';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер на таблицу booking
DROP TRIGGER IF EXISTS check_booking_overlap_trigger ON booking;
CREATE TRIGGER check_booking_overlap_trigger
BEFORE INSERT OR UPDATE ON booking
FOR EACH ROW
EXECUTE FUNCTION check_booking_overlap();