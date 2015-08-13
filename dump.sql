--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

--
-- Name: add_phone(character varying, integer); Type: FUNCTION; Schema: public; Owner: asterisk
--

CREATE FUNCTION add_phone(character varying, integer) RETURNS integer
    LANGUAGE sql
    AS $_$
	select setval ('ldapx_contacts_id_seq', (select case when max(id) is null then 1 else max(id) end from ldapx_contacts));
	insert into ldapx_contacts (id,phone,pers_id)
		values (nextval('ldapx_contacts_id_seq'),$1,$2);
	select max(id) from ldapx_contacts
$_$;


ALTER FUNCTION public.add_phone(character varying, integer) OWNER TO asterisk;

--
-- Name: create_doc(); Type: FUNCTION; Schema: public; Owner: asterisk
--

CREATE FUNCTION create_doc() RETURNS integer
    LANGUAGE sql
    AS $$
	select setval ('documents_id_seq', (select case when max(id) is null then 1 else max(id) end from documents));
	insert into documents (id,title,abstract) 
		values ((select case when max(id) is null then 1 else nextval('documents_id_seq') end from documents),'','');
	select max(id) from documents
$$;


ALTER FUNCTION public.create_doc() OWNER TO asterisk;

--
-- Name: create_o(); Type: FUNCTION; Schema: public; Owner: asterisk
--

CREATE FUNCTION create_o() RETURNS integer
    LANGUAGE sql
    AS $$
	select setval ('ldapx_institutes_id_seq', (select case when max(id) is null then 1 else max(id) end from ldapx_institutes));
	insert into ldapx_institutes (id,name) 
		values ((select case when max(id) is null then 1 else nextval('ldapx_institutes_id_seq') end from ldapx_institutes),'');
	select max(id) from ldapx_institutes
$$;


ALTER FUNCTION public.create_o() OWNER TO asterisk;

--
-- Name: create_person(); Type: FUNCTION; Schema: public; Owner: asterisk
--

CREATE FUNCTION create_person() RETURNS integer
    LANGUAGE sql
    AS $$
	select setval ('ldapx_persons_id_seq', (select case when max(id) is null then 1 else max(id) end from ldapx_persons));
	insert into ldapx_persons (id,name,surname) 
		values ((select case when max(id) is null then 1 else nextval('ldapx_persons_id_seq') end from ldapx_persons),'','');
	select max(id) from ldapx_persons
$$;


ALTER FUNCTION public.create_person() OWNER TO asterisk;

--
-- Name: create_referral(); Type: FUNCTION; Schema: public; Owner: asterisk
--

CREATE FUNCTION create_referral() RETURNS integer
    LANGUAGE sql
    AS $$
	select setval ('referrals_id_seq', (select case when max(id) is null then 1 else max(id) end from referrals));
	insert into referrals (id,name,url) 
		values ((select case when max(id) is null then 1 else nextval('referrals_id_seq') end from referrals),'','');
	select max(id) from referrals
$$;


ALTER FUNCTION public.create_referral() OWNER TO asterisk;

--
-- Name: update_person_cn(character varying, integer); Type: FUNCTION; Schema: public; Owner: asterisk
--

CREATE FUNCTION update_person_cn(character varying, integer) RETURNS integer
    LANGUAGE sql
    AS $_$
	update ldapx_persons set name = (
		select case 
			when position(' ' in $1) = 0 then $1 
			else substr($1, 1, position(' ' in $1) - 1)
		end
	),surname = (
		select case 
			when position(' ' in $1) = 0 then ''
			else substr($1, position(' ' in $1) + 1) 
		end
	) where id = $2;
	select $2 as return
$_$;


ALTER FUNCTION public.update_person_cn(character varying, integer) OWNER TO asterisk;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: ldap_attr_mappings; Type: TABLE; Schema: public; Owner: asterisk; Tablespace: 
--

CREATE TABLE ldap_attr_mappings (
    id integer NOT NULL,
    oc_map_id integer NOT NULL,
    name character varying(255) NOT NULL,
    sel_expr character varying(255) NOT NULL,
    sel_expr_u character varying(255),
    from_tbls character varying(255) NOT NULL,
    join_where character varying(255),
    add_proc character varying(255),
    delete_proc character varying(255),
    param_order integer NOT NULL,
    expect_return integer NOT NULL
);


ALTER TABLE public.ldap_attr_mappings OWNER TO asterisk;

--
-- Name: ldap_attr_mappings_id_seq; Type: SEQUENCE; Schema: public; Owner: asterisk
--

CREATE SEQUENCE ldap_attr_mappings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ldap_attr_mappings_id_seq OWNER TO asterisk;

--
-- Name: ldap_attr_mappings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: asterisk
--

ALTER SEQUENCE ldap_attr_mappings_id_seq OWNED BY ldap_attr_mappings.id;


--
-- Name: ldap_entries; Type: TABLE; Schema: public; Owner: asterisk; Tablespace: 
--

CREATE TABLE ldap_entries (
    id integer NOT NULL,
    dn character varying(2048) NOT NULL,
    oc_map_id integer NOT NULL,
    parent integer NOT NULL,
    keyval integer NOT NULL,
    uid bytea,
    idparent bytea,
    pass integer
);


ALTER TABLE public.ldap_entries OWNER TO asterisk;

--
-- Name: ldap_entries_id_seq; Type: SEQUENCE; Schema: public; Owner: asterisk
--

CREATE SEQUENCE ldap_entries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ldap_entries_id_seq OWNER TO asterisk;

--
-- Name: ldap_entries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: asterisk
--

ALTER SEQUENCE ldap_entries_id_seq OWNED BY ldap_entries.id;


--
-- Name: ldap_entry_objclasses; Type: TABLE; Schema: public; Owner: asterisk; Tablespace: 
--

CREATE TABLE ldap_entry_objclasses (
    entry_id integer NOT NULL,
    oc_name character varying(64)
);


ALTER TABLE public.ldap_entry_objclasses OWNER TO asterisk;

--
-- Name: ldap_oc_mappings; Type: TABLE; Schema: public; Owner: asterisk; Tablespace: 
--

CREATE TABLE ldap_oc_mappings (
    id integer NOT NULL,
    name character varying(64) NOT NULL,
    keytbl character varying(64) NOT NULL,
    keycol character varying(64) NOT NULL,
    create_proc character varying(255),
    delete_proc character varying(255),
    expect_return integer NOT NULL
);


ALTER TABLE public.ldap_oc_mappings OWNER TO asterisk;

--
-- Name: ldap_oc_mappings_id_seq; Type: SEQUENCE; Schema: public; Owner: asterisk
--

CREATE SEQUENCE ldap_oc_mappings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ldap_oc_mappings_id_seq OWNER TO asterisk;

--
-- Name: ldap_oc_mappings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: asterisk
--

ALTER SEQUENCE ldap_oc_mappings_id_seq OWNED BY ldap_oc_mappings.id;


--
-- Name: ldapx_contacts; Type: TABLE; Schema: public; Owner: asterisk; Tablespace: 
--

CREATE TABLE ldapx_contacts (
    id integer NOT NULL,
    phone character varying(255),
    pers_id integer NOT NULL,
    mobile character varying(255),
    mail character varying(255)
);


ALTER TABLE public.ldapx_contacts OWNER TO asterisk;

--
-- Name: ldapx_contacts_id_seq; Type: SEQUENCE; Schema: public; Owner: asterisk
--

CREATE SEQUENCE ldapx_contacts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ldapx_contacts_id_seq OWNER TO asterisk;

--
-- Name: ldapx_contacts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: asterisk
--

ALTER SEQUENCE ldapx_contacts_id_seq OWNED BY ldapx_contacts.id;


--
-- Name: ldapx_institutes; Type: TABLE; Schema: public; Owner: asterisk; Tablespace: 
--

CREATE TABLE ldapx_institutes (
    id integer NOT NULL,
    name character varying(255),
    uid bytea,
    idparent bytea,
    pass integer
);


ALTER TABLE public.ldapx_institutes OWNER TO asterisk;

--
-- Name: ldapx_institutes_id_seq; Type: SEQUENCE; Schema: public; Owner: asterisk
--

CREATE SEQUENCE ldapx_institutes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ldapx_institutes_id_seq OWNER TO asterisk;

--
-- Name: ldapx_institutes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: asterisk
--

ALTER SEQUENCE ldapx_institutes_id_seq OWNED BY ldapx_institutes.id;


--
-- Name: ldapx_persons; Type: TABLE; Schema: public; Owner: asterisk; Tablespace: 
--

CREATE TABLE ldapx_persons (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    surname character varying(255) NOT NULL,
    password character varying(64),
    mn character varying(255),
    uid bytea,
    idparent bytea
);


ALTER TABLE public.ldapx_persons OWNER TO asterisk;

--
-- Name: ldapx_persons_id_seq; Type: SEQUENCE; Schema: public; Owner: asterisk
--

CREATE SEQUENCE ldapx_persons_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ldapx_persons_id_seq OWNER TO asterisk;

--
-- Name: ldapx_persons_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: asterisk
--

ALTER SEQUENCE ldapx_persons_id_seq OWNED BY ldapx_persons.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldap_attr_mappings ALTER COLUMN id SET DEFAULT nextval('ldap_attr_mappings_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldap_entries ALTER COLUMN id SET DEFAULT nextval('ldap_entries_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldap_oc_mappings ALTER COLUMN id SET DEFAULT nextval('ldap_oc_mappings_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldapx_contacts ALTER COLUMN id SET DEFAULT nextval('ldapx_contacts_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldapx_institutes ALTER COLUMN id SET DEFAULT nextval('ldapx_institutes_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldapx_persons ALTER COLUMN id SET DEFAULT nextval('ldapx_persons_id_seq'::regclass);


--
-- Data for Name: ldap_attr_mappings; Type: TABLE DATA; Schema: public; Owner: asterisk
--

INSERT INTO ldap_attr_mappings VALUES (2, 1, 'telephoneNumber', 'ldapx_contacts.phone', NULL, 'ldapx_persons,ldapx_contacts', 'ldapx_contacts.pers_id=ldapx_persons.id', 'SELECT add_phone(?,?)', 'DELETE FROM ldapx_contacts WHERE phone=? AND pers_id=?', 3, 0);
INSERT INTO ldap_attr_mappings VALUES (3, 1, 'givenName', 'ldapx_persons.name', NULL, 'ldapx_persons', NULL, 'UPDATE ldapx_persons SET name=? WHERE id=?', 'UPDATE ldapx_persons SET name='''' WHERE (name=? OR name='''') AND id=?', 3, 0);
INSERT INTO ldap_attr_mappings VALUES (4, 1, 'sn', 'ldapx_persons.surname', NULL, 'ldapx_persons', NULL, 'UPDATE ldapx_persons SET surname=? WHERE id=?', 'UPDATE ldapx_persons SET surname='''' WHERE (surname=? OR surname='''') AND id=?', 3, 0);
INSERT INTO ldap_attr_mappings VALUES (5, 1, 'userPassword', 'ldapx_persons.password', NULL, 'ldapx_persons', 'ldapx_persons.password IS NOT NULL', 'UPDATE ldapx_persons SET password=? WHERE id=?', 'UPDATE ldapx_persons SET password=NULL WHERE password=? AND id=?', 3, 0);
INSERT INTO ldap_attr_mappings VALUES (11, 3, 'o', 'ldapx_institutes.name', NULL, 'ldapx_institutes', NULL, 'UPDATE ldapx_institutes SET name=? WHERE id=?', 'UPDATE ldapx_institutes SET name='''' WHERE name=? AND id=?', 3, 0);
INSERT INTO ldap_attr_mappings VALUES (12, 3, 'dc', 'lower(ldapx_institutes.name)', NULL, 'ldapx_institutes,ldap_entries AS dcObject,ldap_entry_objclasses AS auxObjectClass', 'ldapx_institutes.id=dcObject.keyval AND dcObject.oc_map_id=3 AND dcObject.id=auxObjectClass.entry_id AND auxObjectClass.oc_name=''dcObject''', NULL, 'SELECT 1 FROM ldapx_institutes WHERE lower(name)=? AND id=? and 1=0', 3, 0);
INSERT INTO ldap_attr_mappings VALUES (1, 1, 'cn', 'text(ldapx_persons.name||'' ''||ldapx_persons.mn||'' ''||ldapx_persons.surname)', NULL, 'ldapx_persons', NULL, 'SELECT update_person_cn(?,?)', 'SELECT 1 FROM ldapx_persons WHERE ldapx_persons.name=? AND ldapx_persons.id=? AND 1=0', 3, 0);
INSERT INTO ldap_attr_mappings VALUES (16, 1, 'mobile', 'ldapx_contacts.mobile', NULL, 'ldapx_persons,ldapx_contacts', 'ldapx_contacts.pers_id=ldapx_persons.id', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (17, 1, 'mail', 'ldapx_contacts.mail', NULL, 'ldapx_persons,ldapx_contacts', 'ldapx_contacts.pers_id=ldapx_persons.id', NULL, NULL, 3, 0);


--
-- Name: ldap_attr_mappings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: asterisk
--

SELECT pg_catalog.setval('ldap_attr_mappings_id_seq', 1, false);


--
-- Data for Name: ldap_entries; Type: TABLE DATA; Schema: public; Owner: asterisk
--

INSERT INTO ldap_entries VALUES (1, 'O=Enterprise', 3, 0, 1, '\x31', '\x30', 0);
INSERT INTO ldap_entries VALUES (2, 'OU=Quadra,O=Enterprise', 3, 1, 2, '\x32', '\x31', 0);


--
-- Name: ldap_entries_id_seq; Type: SEQUENCE SET; Schema: public; Owner: asterisk
--

SELECT pg_catalog.setval('ldap_entries_id_seq', 94, true);


--
-- Data for Name: ldap_entry_objclasses; Type: TABLE DATA; Schema: public; Owner: asterisk
--

INSERT INTO ldap_entry_objclasses VALUES (1, 'dcObject');


--
-- Data for Name: ldap_oc_mappings; Type: TABLE DATA; Schema: public; Owner: asterisk
--

INSERT INTO ldap_oc_mappings VALUES (1, 'inetOrgPerson', 'ldapx_persons', 'id', 'SELECT create_person()', 'DELETE FROM ldapx_persons WHERE id=?', 0);
INSERT INTO ldap_oc_mappings VALUES (3, 'organization', 'ldapx_institutes', 'id', 'SELECT create_o()', 'DELETE FROM ldapx_institutes WHERE id=?', 0);


--
-- Name: ldap_oc_mappings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: asterisk
--

SELECT pg_catalog.setval('ldap_oc_mappings_id_seq', 1, false);


--
-- Data for Name: ldapx_contacts; Type: TABLE DATA; Schema: public; Owner: asterisk
--



--
-- Name: ldapx_contacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: asterisk
--

SELECT pg_catalog.setval('ldapx_contacts_id_seq', 1, false);


--
-- Data for Name: ldapx_institutes; Type: TABLE DATA; Schema: public; Owner: asterisk
--

INSERT INTO ldapx_institutes VALUES (1, 'Вселенная', '\x31', '\x30', 0);
INSERT INTO ldapx_institutes VALUES (2, 'Квадра', '\x32', '\x31', 0);


--
-- Name: ldapx_institutes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: asterisk
--

SELECT pg_catalog.setval('ldapx_institutes_id_seq', 78, true);


--
-- Data for Name: ldapx_persons; Type: TABLE DATA; Schema: public; Owner: asterisk
--



--
-- Name: ldapx_persons_id_seq; Type: SEQUENCE SET; Schema: public; Owner: asterisk
--

SELECT pg_catalog.setval('ldapx_persons_id_seq', 1, false);


--
-- Name: ldap_attr_mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldap_attr_mappings
    ADD CONSTRAINT ldap_attr_mappings_pkey PRIMARY KEY (id);


--
-- Name: ldap_entries_dn_key; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldap_entries
    ADD CONSTRAINT ldap_entries_dn_key UNIQUE (dn);


--
-- Name: ldap_entries_oc_map_id_keyval_key; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldap_entries
    ADD CONSTRAINT ldap_entries_oc_map_id_keyval_key UNIQUE (oc_map_id, keyval);


--
-- Name: ldap_entries_pkey; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldap_entries
    ADD CONSTRAINT ldap_entries_pkey PRIMARY KEY (id);


--
-- Name: ldap_oc_mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldap_oc_mappings
    ADD CONSTRAINT ldap_oc_mappings_pkey PRIMARY KEY (id);


--
-- Name: ldapx_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldapx_contacts
    ADD CONSTRAINT ldapx_contacts_pkey PRIMARY KEY (id);


--
-- Name: ldapx_institutes_pkey; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldapx_institutes
    ADD CONSTRAINT ldapx_institutes_pkey PRIMARY KEY (id);


--
-- Name: ldapx_persons_pkey; Type: CONSTRAINT; Schema: public; Owner: asterisk; Tablespace: 
--

ALTER TABLE ONLY ldapx_persons
    ADD CONSTRAINT ldapx_persons_pkey PRIMARY KEY (id);


--
-- Name: ldap_attr_mappings_oc_map_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldap_attr_mappings
    ADD CONSTRAINT ldap_attr_mappings_oc_map_id_fkey FOREIGN KEY (oc_map_id) REFERENCES ldap_oc_mappings(id);


--
-- Name: ldap_entries_oc_map_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldap_entries
    ADD CONSTRAINT ldap_entries_oc_map_id_fkey FOREIGN KEY (oc_map_id) REFERENCES ldap_oc_mappings(id);


--
-- Name: ldap_entry_objclasses_entry_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: asterisk
--

ALTER TABLE ONLY ldap_entry_objclasses
    ADD CONSTRAINT ldap_entry_objclasses_entry_id_fkey FOREIGN KEY (entry_id) REFERENCES ldap_entries(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: pgsql
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM pgsql;
GRANT ALL ON SCHEMA public TO pgsql;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

