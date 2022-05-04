CREATE TABLE IF NOT EXISTS menus (
   id uuid PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   is_enabled BOOLEAN NOT NULL DEFAULT FALSE 
);

CREATE TABLE IF NOT EXISTS categories (
   id uuid PRIMARY KEY,
   name VARCHAR (50) NOT NULL
);

CREATE TABLE IF NOT EXISTS menus_categories (
   menu_id uuid NOT NULL,
   category_id uuid NOT NULL UNIQUE,
   PRIMARY KEY(menu_id, category_id),
   FOREIGN KEY(menu_id) REFERENCES menus(id),
   FOREIGN KEY(category_id) REFERENCES categories(id)
);