-- Database: goauth

-- DROP DATABASE goauth;

CREATE DATABASE goauth
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE goauth
    IS 'Has all tables required for goauth OAuth2 server';

CREATE TYPE tokenType AS ENUM ('Reference','Bearer');

-- Table: public.clients

-- Table: public.clients

DROP TABLE public.clients;

CREATE TABLE public.clients
(
    "clientId" character varying(255) COLLATE pg_catalog."default" NOT NULL,
    "clientSecret" character varying(4096) COLLATE pg_catalog."default" NOT NULL,
    "accessTokenLifetime" bigint NOT NULL,
    "tokenType" tokentype NOT NULL,
    "allowedScopes" text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT clients_pkey PRIMARY KEY ("clientId")
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.clients
    OWNER to postgres;
COMMENT ON TABLE public.clients
    IS 'clientSecret should be encrypted by your application.';