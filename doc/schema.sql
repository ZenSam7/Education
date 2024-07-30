--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: articles; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.articles (
    id_article integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    edited_at timestamp with time zone,
    title character varying NOT NULL,
    text text NOT NULL,
    comments integer[],
    authors integer[] NOT NULL,
    evaluation integer DEFAULT 0 NOT NULL,
    CONSTRAINT articles_text_check CHECK ((text <> ''::text)),
    CONSTRAINT articles_title_check CHECK (((title)::text <> ''::text))
);


ALTER TABLE public.articles OWNER TO root;

--
-- Name: articles_id_article_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.articles_id_article_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.articles_id_article_seq OWNER TO root;

--
-- Name: articles_id_article_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.articles_id_article_seq OWNED BY public.articles.id_article;


--
-- Name: comments; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.comments (
    id_comment integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    edited_at timestamp with time zone,
    text text NOT NULL,
    author integer NOT NULL,
    evaluation integer DEFAULT 0 NOT NULL,
    CONSTRAINT comments_text_check CHECK ((text <> ''::text))
);


ALTER TABLE public.comments OWNER TO root;

--
-- Name: comments_id_comment_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.comments_id_comment_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.comments_id_comment_seq OWNER TO root;

--
-- Name: comments_id_comment_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.comments_id_comment_seq OWNED BY public.comments.id_comment;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO root;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.sessions (
    id_session uuid NOT NULL,
    issued_at timestamp with time zone DEFAULT now(),
    expired_at timestamp with time zone NOT NULL,
    refresh_token character varying NOT NULL,
    id_user integer NOT NULL,
    client_ip character varying NOT NULL,
    blocked boolean DEFAULT false NOT NULL
);


ALTER TABLE public.sessions OWNER TO root;

--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    id_user integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    name character varying NOT NULL,
    description text,
    karma integer DEFAULT 0 NOT NULL,
    email character varying NOT NULL,
    password_hash character varying NOT NULL,
    email_verified boolean DEFAULT false NOT NULL,
    CONSTRAINT email_must_be_valid CHECK (((email)::text ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$'::text)),
    CONSTRAINT users_email_check CHECK (((email)::text <> ''::text)),
    CONSTRAINT users_name_check CHECK (((name)::text <> ''::text)),
    CONSTRAINT users_password_hash_check CHECK (((password_hash)::text <> ''::text))
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: users_id_user_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.users_id_user_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_user_seq OWNER TO root;

--
-- Name: users_id_user_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.users_id_user_seq OWNED BY public.users.id_user;


--
-- Name: verify_emails; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.verify_emails (
    id_verify_email integer NOT NULL,
    id_user integer NOT NULL,
    secret_key character varying NOT NULL,
    expired_at timestamp with time zone DEFAULT (now() + '00:10:00'::interval) NOT NULL
);


ALTER TABLE public.verify_emails OWNER TO root;

--
-- Name: verify_emails_id_verify_email_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.verify_emails_id_verify_email_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.verify_emails_id_verify_email_seq OWNER TO root;

--
-- Name: verify_emails_id_verify_email_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.verify_emails_id_verify_email_seq OWNED BY public.verify_emails.id_verify_email;


--
-- Name: articles id_article; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.articles ALTER COLUMN id_article SET DEFAULT nextval('public.articles_id_article_seq'::regclass);


--
-- Name: comments id_comment; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.comments ALTER COLUMN id_comment SET DEFAULT nextval('public.comments_id_comment_seq'::regclass);


--
-- Name: users id_user; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users ALTER COLUMN id_user SET DEFAULT nextval('public.users_id_user_seq'::regclass);


--
-- Name: verify_emails id_verify_email; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.verify_emails ALTER COLUMN id_verify_email SET DEFAULT nextval('public.verify_emails_id_verify_email_seq'::regclass);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id_article);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id_comment);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id_session);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_password_hash_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_password_hash_key UNIQUE (password_hash);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id_user);


--
-- Name: verify_emails verify_emails_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.verify_emails
    ADD CONSTRAINT verify_emails_pkey PRIMARY KEY (id_verify_email);


--
-- Name: article_indx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX article_indx ON public.articles USING btree (id_article);


--
-- Name: comment_indx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX comment_indx ON public.comments USING btree (id_comment);


--
-- Name: user_indx; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX user_indx ON public.users USING btree (id_user);


--
-- Name: comments comments_author_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_author_fkey FOREIGN KEY (author) REFERENCES public.users(id_user);


--
-- Name: sessions sessions_id_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_id_user_fkey FOREIGN KEY (id_user) REFERENCES public.users(id_user);


--
-- Name: verify_emails verify_emails_id_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.verify_emails
    ADD CONSTRAINT verify_emails_id_user_fkey FOREIGN KEY (id_user) REFERENCES public.users(id_user);


--
-- PostgreSQL database dump complete
--

