--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2
-- Dumped by pg_dump version 12.2

-- Started on 2020-06-11 20:49:18

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

--
-- TOC entry 3 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 2822 (class 0 OID 0)
-- Dependencies: 3
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 202 (class 1259 OID 17526)
-- Name: Notes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Notes" (
    "Views" bigint,
    "Content" text,
    "Id" text,
    "Lang" text,
    "LastView" bigint,
    "OwnerId" text
);


ALTER TABLE public."Notes" OWNER TO postgres;

--
-- TOC entry 203 (class 1259 OID 17551)
-- Name: Views; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Views" (
    "PostId" text,
    "Ip" text
);


ALTER TABLE public."Views" OWNER TO postgres;

--
-- TOC entry 2823 (class 0 OID 0)
-- Dependencies: 202
-- Name: TABLE "Notes"; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public."Notes" TO gonotes;


--
-- TOC entry 2824 (class 0 OID 0)
-- Dependencies: 203
-- Name: TABLE "Views"; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public."Views" TO gonotes;


-- Completed on 2020-06-11 20:49:18

--
-- PostgreSQL database dump complete
--

