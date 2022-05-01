CREATE TABLE IF NOT EXISTS menus (
   id uuid PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   is_enabled BOOLEAN NOT NULL DEFAULT FALSE 
);

CREATE TABLE IF NOT EXISTS categories (
   id uuid PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   image_url VARCHAR (50) NOT NULL
);

CREATE TABLE IF NOT EXISTS menu_categories (
   menu_id uuid,
   category_id uuid,
   PRIMARY KEY(menu_id, category_id),
   FOREIGN KEY(menu_id) REFERENCES menus(id),
   FOREIGN KEY(category_id) REFERENCES categories(id)
);