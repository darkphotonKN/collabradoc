--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1 (Debian 16.1-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Homebrew)

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
-- Name: comments; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.comments (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    comment text NOT NULL,
    author text NOT NULL,
    document_id bigint NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE public.comments OWNER TO myuser;

--
-- Name: comments_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.comments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.comments_id_seq OWNER TO myuser;

--
-- Name: comments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.comments_id_seq OWNED BY public.comments.id;


--
-- Name: documents; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.documents (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text NOT NULL,
    content text,
    user_id bigint NOT NULL,
    live_session_id bigint
);


ALTER TABLE public.documents OWNER TO myuser;

--
-- Name: documents_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.documents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.documents_id_seq OWNER TO myuser;

--
-- Name: documents_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.documents_id_seq OWNED BY public.documents.id;


--
-- Name: live_session_users; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.live_session_users (
    live_session_id bigint NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE public.live_session_users OWNER TO myuser;

--
-- Name: live_sessions; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.live_sessions (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    session_id text NOT NULL,
    document_id bigint,
    is_active boolean
);


ALTER TABLE public.live_sessions OWNER TO myuser;

--
-- Name: live_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.live_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.live_sessions_id_seq OWNER TO myuser;

--
-- Name: live_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.live_sessions_id_seq OWNED BY public.live_sessions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL
);


ALTER TABLE public.users OWNER TO myuser;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO myuser;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: comments id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.comments ALTER COLUMN id SET DEFAULT nextval('public.comments_id_seq'::regclass);


--
-- Name: documents id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.documents ALTER COLUMN id SET DEFAULT nextval('public.documents_id_seq'::regclass);


--
-- Name: live_sessions id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.live_sessions ALTER COLUMN id SET DEFAULT nextval('public.live_sessions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: documents documents_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_pkey PRIMARY KEY (id);


--
-- Name: live_session_users live_session_users_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.live_session_users
    ADD CONSTRAINT live_session_users_pkey PRIMARY KEY (live_session_id, user_id);


--
-- Name: live_sessions live_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.live_sessions
    ADD CONSTRAINT live_sessions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_comments_deleted_at; Type: INDEX; Schema: public; Owner: myuser
--

CREATE INDEX idx_comments_deleted_at ON public.comments USING btree (deleted_at);


--
-- Name: idx_documents_deleted_at; Type: INDEX; Schema: public; Owner: myuser
--

CREATE INDEX idx_documents_deleted_at ON public.documents USING btree (deleted_at);


--
-- Name: idx_live_sessions_deleted_at; Type: INDEX; Schema: public; Owner: myuser
--

CREATE INDEX idx_live_sessions_deleted_at ON public.live_sessions USING btree (deleted_at);


--
-- Name: idx_live_sessions_session_id; Type: INDEX; Schema: public; Owner: myuser
--

CREATE UNIQUE INDEX idx_live_sessions_session_id ON public.live_sessions USING btree (session_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: myuser
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: comments fk_documents_comment; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT fk_documents_comment FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: live_sessions fk_documents_live_session; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.live_sessions
    ADD CONSTRAINT fk_documents_live_session FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: live_session_users fk_live_session_users_live_session; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.live_session_users
    ADD CONSTRAINT fk_live_session_users_live_session FOREIGN KEY (live_session_id) REFERENCES public.live_sessions(id);


--
-- Name: live_session_users fk_live_session_users_user; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.live_session_users
    ADD CONSTRAINT fk_live_session_users_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: comments fk_users_comment; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT fk_users_comment FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: documents fk_users_doc; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT fk_users_doc FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

