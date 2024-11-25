-- public.summaries definition

-- Drop table

-- DROP TABLE public.summaries;

CREATE TABLE public.summaries (
	id serial4 NOT NULL,
	category int4 NOT NULL,
	budget_amount numeric NOT NULL,
	"month" int2 NOT NULL,
	"year" int2 NOT NULL,
	CONSTRAINT summaries_pk PRIMARY KEY (id)
);
COMMENT ON TABLE public.summaries IS 'This table holds the home screen''s summary "budgetting" data. Each summary has a date, so that each month a new summary can be shown. Each summary is attached to a category.';


-- public.summaries foreign keys

ALTER TABLE public.summaries ADD CONSTRAINT summaries_categories_fk FOREIGN KEY (category) REFERENCES public.categories(id);
