-- DROP TYPE public."color";

CREATE TYPE public."color" AS ENUM (
	'red',
	'green',
	'blue',
	'purple',
	'yellow',
	'none',
	'grey');

-- public.category_colors definition

-- Drop table

-- DROP TABLE public.category_colors;

CREATE TABLE public.category_colors (
	id serial4 NOT NULL,
	"color" public."color" DEFAULT 'none'::color NOT NULL, -- A defined color for use in categories.
	CONSTRAINT category_colors_pkey PRIMARY KEY (id),
	CONSTRAINT category_colors_unique UNIQUE (color)
);

-- Column comments

COMMENT ON COLUMN public.category_colors."color" IS 'A defined color for use in categories.';
