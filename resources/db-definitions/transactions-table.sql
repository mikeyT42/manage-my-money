-- DROP TYPE public."transaction_type";

CREATE TYPE public."transaction_type" AS ENUM (
	'withdrawal',
	'deposit');

-- public.transactions definition

-- Drop table

-- DROP TABLE public.transactions;

CREATE TABLE public.transactions (
	id serial4 NOT NULL,
	bank_id int4 NOT NULL, -- The unique id assigned to a transaction by the bank.
	"transaction_type" public."transaction_type" NOT NULL, -- The transaction type. It will only be withdrawls or deposits.
	amount numeric NOT NULL, -- If it's positive the transaction is a deposit; a negative amount will be a withdrawal.
	"date" date NOT NULL, -- The date of the transaction given by the bank.
	bank_description text DEFAULT ''::text NOT NULL, -- The description given by the bank.
	user_description text DEFAULT ''::text NOT NULL, -- The description given by the user.
	category int4 DEFAULT 1 NOT NULL, -- The foreign key to the categories table.
	CONSTRAINT transactions_pkey PRIMARY KEY (id),
	CONSTRAINT transactions_unique UNIQUE (bank_id)
);
COMMENT ON TABLE public.transactions IS 'This table holds all of the transactions from the imported bank statement.';

-- Column comments

COMMENT ON COLUMN public.transactions.bank_id IS 'The unique id assigned to a transaction by the bank.';
COMMENT ON COLUMN public.transactions."transaction_type" IS 'The transaction type. It will only be withdrawls or deposits.';
COMMENT ON COLUMN public.transactions.amount IS 'If it''s positive the transaction is a deposit; a negative amount will be a withdrawal.';
COMMENT ON COLUMN public.transactions."date" IS 'The date of the transaction given by the bank.';
COMMENT ON COLUMN public.transactions.bank_description IS 'The description given by the bank.';
COMMENT ON COLUMN public.transactions.user_description IS 'The description given by the user.';
COMMENT ON COLUMN public.transactions.category IS 'The foreign key to the categories table.';


-- public.transactions foreign keys

ALTER TABLE public.transactions ADD CONSTRAINT transactions_categories_fk FOREIGN KEY (category) REFERENCES public.categories(id);
