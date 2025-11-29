-- CREATE TABLE Examples

-- Simple table
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE
);

-- Table with IF NOT EXISTS
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) DEFAULT 0.00,
    category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table with foreign keys
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT DEFAULT 1,
    total DECIMAL(10, 2),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE SET NULL
);

-- Table with composite primary key
CREATE TABLE user_roles (
    user_id INT,
    role_id INT,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);

-- DROP Examples
DROP TABLE IF EXISTS old_users;
DROP DATABASE IF EXISTS test_db;
DROP INDEX idx_users_email;

-- ALTER TABLE Examples
ALTER TABLE users ADD COLUMN age INT;
ALTER TABLE users ADD COLUMN phone VARCHAR(20) NOT NULL UNIQUE;
ALTER TABLE users DROP COLUMN age;
ALTER TABLE users MODIFY COLUMN name VARCHAR(150) NOT NULL;
ALTER TABLE users CHANGE COLUMN email user_email VARCHAR(255);
ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE orders DROP CONSTRAINT fk_user;

-- CREATE INDEX Examples
CREATE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_products_category ON products (category);
CREATE INDEX idx_orders_user_product ON orders (user_id, product_id);
