CREATE TABLE public.pack_size (
	id serial NOT NULL,
	"size" int NOT NULL,
	CONSTRAINT pack_size_pk PRIMARY KEY (id),
	CONSTRAINT pack_size_unique UNIQUE ("size")
);

INSERT INTO public.pack_size(size)
VALUES
    (5000),
    (2000),
    (1000),
    (500),
    (250);


CREATE TABLE public."order" (
	id uuid NOT NULL,
	amount_items int NOT NULL,
	packs jsonb NOT NULL,
	CONSTRAINT order_pk PRIMARY KEY (id)
);
