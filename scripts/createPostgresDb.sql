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


CREATE TABLE public.tokens
(
    "tokenID" uuid NOT NULL,
    "clientID" character varying(255) COLLATE pg_catalog."default" NOT NULL,
    expiry bigint NOT NULL,
    "accessToken" text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT tokens_pkey PRIMARY KEY ("tokenID")
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.tokens
    OWNER to postgres;
COMMENT ON TABLE public.tokens
    IS 'Store the JWT tokens behind the reference tokens';