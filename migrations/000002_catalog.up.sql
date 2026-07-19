-- User profile: add photo_url
ALTER TABLE users ADD COLUMN IF NOT EXISTS photo_url TEXT;

-- Catalog: categories
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon_name VARCHAR(100) NOT NULL,
    color_hex VARCHAR(9) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Catalog: products
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    rating DECIMAL(3,2) NOT NULL DEFAULT 0,
    reviews_count INT NOT NULL DEFAULT 0,
    icon_name VARCHAR(100) NOT NULL,
    color_hex VARCHAR(9) NOT NULL,
    category_id VARCHAR(50) NOT NULL REFERENCES categories(id),
    description TEXT NOT NULL DEFAULT '',
    image_urls TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);

-- Catalog: promos
CREATE TABLE IF NOT EXISTS promos (
    id VARCHAR(100) PRIMARY KEY,
    badge VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) NOT NULL,
    icon_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Cart items (per user)
CREATE TABLE IF NOT EXISTS cart_items (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id VARCHAR(100) NOT NULL REFERENCES products(id),
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, product_id)
);

-- Wishlist (per user)
CREATE TABLE IF NOT EXISTS wishlists (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id VARCHAR(100) NOT NULL REFERENCES products(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, product_id)
);

-- Orders
CREATE TABLE IF NOT EXISTS store_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    total DECIMAL(12,2) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'processing',
    shipping_address TEXT NOT NULL DEFAULT '',
    payment_method VARCHAR(30) NOT NULL DEFAULT 'cod',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_store_orders_user ON store_orders(user_id);

-- Order items
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES store_orders(id) ON DELETE CASCADE,
    product_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    quantity INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);

-- Reviews
CREATE TABLE IF NOT EXISTS product_reviews (
    user_id UUID NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    product_id VARCHAR(100) NOT NULL,
    order_id VARCHAR(100) NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, product_id)
);

-- Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    subtitle TEXT NOT NULL DEFAULT '',
    type VARCHAR(30) NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, created_at DESC);