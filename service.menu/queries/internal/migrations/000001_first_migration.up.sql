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

CREATE TABLE IF NOT EXISTS subcategories (
   id uuid PRIMARY KEY,
   name VARCHAR (50) NOT NULL
);

CREATE TABLE IF NOT EXISTS category_subcategories (
   category_id uuid NOT NULL,
   subcategory_id uuid NOT NULL UNIQUE,
   PRIMARY KEY(category_id, subcategory_id),
   FOREIGN KEY(category_id) REFERENCES categories(id),
   FOREIGN KEY(subcategory_id) REFERENCES subcategories(id)
);


CREATE TABLE IF NOT EXISTS menuitems (
   id uuid PRIMARY KEY,
   name VARCHAR (50) NOT NULL
);

CREATE TABLE IF NOT EXISTS subcategory_menuitems (
   subcategory_id uuid NOT NULL,
   menuitem_id uuid NOT NULL UNIQUE,
   PRIMARY KEY(subcategory_id, menuitem_id),
   FOREIGN KEY(subcategory_id) REFERENCES subcategories(id),
   FOREIGN KEY(menuitem_id) REFERENCES menuitems(id)
);