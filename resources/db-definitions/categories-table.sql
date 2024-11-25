-- public.categories definition

-- Drop table

-- DROP TABLE public.categories;

CREATE TABLE public.categories (
	id serial4 NOT NULL,
	description text DEFAULT ''::text NOT NULL, -- The description given to the category.
	"name" text DEFAULT 'None'::text NOT NULL, -- The defined name of the category.
	"color" int4 DEFAULT 1 NOT NULL, -- The foreign key to the category_colors table.
	is_user_defined bool DEFAULT true NOT NULL,
	CONSTRAINT categories_pkey PRIMARY KEY (id),
	CONSTRAINT categories_unique UNIQUE (name)
);
COMMENT ON TABLE public.categories IS 'This table holds the system default and user created categories.';

-- Column comments

COMMENT ON COLUMN public.categories.description IS 'The description given to the category.';
COMMENT ON COLUMN public.categories."name" IS 'The defined name of the category.';
COMMENT ON COLUMN public.categories."color" IS 'The foreign key to the category_colors table.';


-- public.categories foreign keys

ALTER TABLE public.categories ADD CONSTRAINT categories_category_colors_fk FOREIGN KEY ("color") REFERENCES public.category_colors(id);
