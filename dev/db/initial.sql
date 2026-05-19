CREATE TABLE public.pack_size (
	id serial NOT NULL,
	"size" int NOT NULL,
	CONSTRAINT pack_size_pk PRIMARY KEY (id),
	CONSTRAINT pack_size_unique UNIQUE ("size")
);
