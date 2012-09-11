
CREATE TABLE users (
    login character varying(50),
    hpasswd character varying(100)
);

COPY users (login, hpasswd) FROM stdin;
pauek	$2a$10$yL6CXIESa0eGYN5xqI1tfeubfw3LuSOJshq1N.G8AXkvBGtyzN6tO
\.

ALTER TABLE public.users OWNER TO academio;

