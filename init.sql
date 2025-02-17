-- Create advertisers table
CREATE TABLE IF NOT EXISTS advertisers (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    budget DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create advertisements table
CREATE TABLE IF NOT EXISTS advertisements (
    id BIGINT PRIMARY KEY,
    advertiser_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    cpm_count BIGINT NOT NULL DEFAULT 0,
    cpc_count BIGINT NOT NULL DEFAULT 0,
    cpm_rate DECIMAL(10,2) NOT NULL,
    cpc_rate DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (advertiser_id) REFERENCES advertisers(id)
);

-- Insert demo advertisers
INSERT INTO advertisers (id, name, budget, created_at, updated_at) VALUES
(1, 'Tech Corp', 1000.00, NOW(), NOW()),
(2, 'Fashion Brand', 800.00, NOW(), NOW()),
(3, 'Food Delivery', 500.00, NOW(), NOW());

-- Insert demo advertisements
INSERT INTO advertisements (id, advertiser_id, title, content, cpm_count, cpc_count, cpm_rate, cpc_rate, created_at, updated_at) VALUES
(1, 1, 'New Laptop Sale', 'Get 20% off on all laptops!', 0, 0, 2.50, 0.50, NOW(), NOW()),
(2, 1, 'Gaming PC Deal', 'Build your dream gaming PC today', 0, 0, 3.00, 0.75, NOW(), NOW()),
(3, 2, 'Summer Collection', 'Discover our new summer styles', 0, 0, 2.00, 0.40, NOW(), NOW()),
(4, 2, 'Premium Watches', 'Luxury watches at amazing prices', 0, 0, 4.00, 1.00, NOW(), NOW()),
(5, 3, 'Free Delivery', 'Free delivery on your first order', 0, 0, 1.50, 0.30, NOW(), NOW()),
(6, 3, 'Lunch Special', 'Get 15% off lunch orders', 0, 0, 1.75, 0.35, NOW(), NOW());