CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  quantity INT NOT NULL,
  price DOUBLE PRECISION NOT NULL
);