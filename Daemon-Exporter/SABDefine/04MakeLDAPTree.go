package SABDefine

var (
	LDAP_Tables_am = int(9)

	LDAP_Scheme_create = string(`

drop table if exists ldapx_phones;
drop table if exists ldapx_persons;
drop table if exists ldapx_mail;
drop table if exists ldapx_institutes;
drop table if exists ldap_entry_objclasses;
drop table if exists ldap_entries;
drop table if exists ldap_attr_mappings;
drop table if exists ldap_oc_mappings;

CREATE TABLE IF NOT EXISTS ldap_attr_mappings (
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
    expect_return integer NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS ldap_oc_mappings (
    id integer NOT NULL,
    name character varying(64) NOT NULL,
    keytbl character varying(64) NOT NULL,
    keycol character varying(64) NOT NULL,
    create_proc character varying(255),
    delete_proc character varying(255),
    expect_return integer NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS ldap_entry_objclasses (
    entry_id integer NOT NULL,
    oc_name character varying(64)
);

CREATE TABLE IF NOT EXISTS ldap_entries (
    id integer NOT NULL,
    dn character varying(2048) NOT NULL,
    oc_map_id integer NOT NULL,
    parent integer NOT NULL,
    keyval integer NOT NULL,
    uid bytea,
    idparent bytea,
    pass integer,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS ldapx_institutes (
    id integer NOT NULL,
    name character varying(255),
    uid bytea,
    idparent bytea,
    pass integer,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS ldapx_phones (
    id integer NOT NULL,
    phone character varying(255) NOT NULL,
    pers_id bytea NOT NULL,
    pass integer,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS ldapx_persons (
    id integer NOT NULL,
    name character varying(255),
    surname character varying(255),
    password character varying(64),
    mn character varying(255),
    uid bytea,
    idparent bytea,
    login character varying(255),
    fullname character varying(255),
    lang integer NOT NULL,
    bc character varying(255),
    cid_name character varying(255),
    contract integer NULL
);

CREATE TABLE IF NOT EXISTS ldapx_mail (
    id integer NOT NULL,
    mail character varying(255) NOT NULL,
    pers_id bytea NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS ldapx_ad_login (
    id integer NOT NULL,
    domain character varying(255) NOT NULL,
    dlogin character varying(255) NOT NULL,
    login character varying(255) NOT NULL,
    pers_id bytea NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO ldap_attr_mappings VALUES (1, 1, 'cn', 'text(ldapx_persons.surname||'' ''||ldapx_persons.name)', NULL, 'ldapx_persons', 'ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (2, 3, 'o', 'ldapx_institutes.name', NULL, 'ldapx_institutes', NULL, NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (3, 1, 'givenName', 'ldapx_persons.name', NULL, 'ldapx_persons', 'ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (4, 1, 'sn', 'ldapx_persons.surname', NULL, 'ldapx_persons', 'ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (5, 1, 'userPassword', 'ldapx_persons.password', NULL, 'ldapx_persons', 'ldapx_persons.password IS NOT NULL', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (6, 1, 'displayName', 'text(ldapx_persons.fullname)', NULL, 'ldapx_persons', 'ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (7, 3, 'dc', 'lower(ldapx_institutes.name)', NULL, 'ldapx_institutes,ldap_entries AS dcObject,ldap_entry_objclasses AS auxObjectClass', 'ldapx_institutes.id=dcObject.keyval AND dcObject.oc_map_id=3 AND dcObject.id=auxObjectClass.entry_id AND auxObjectClass.oc_name=''dcObject''', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (8, 1, 'mobile', 'ldapx_phones.phone', NULL, 'ldapx_persons,ldapx_phones', 'ldapx_phones.pers_id=ldapx_persons.uid and ldapx_persons.lang=0 and ldapx_phones.pass=1', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (9, 1, 'telephoneNumber', 'ldapx_phones.phone', NULL, 'ldapx_persons,ldapx_phones', 'ldapx_phones.pers_id=ldapx_persons.uid and ldapx_persons.lang=0 and ldapx_phones.pass=2', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (10, 1, 'pager', 'ldapx_phones.phone', NULL, 'ldapx_persons,ldapx_phones', 'ldapx_phones.pers_id=ldapx_persons.uid and ldapx_persons.lang=0 and ldapx_phones.pass=3', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (11, 1, 'mail', 'ldapx_mail.mail', NULL, 'ldapx_persons,ldapx_mail', 'ldapx_mail.pers_id=ldapx_persons.uid and ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (12, 1, 'businessCategory', 'ldapx_persons.bc', NULL, 'ldapx_persons', 'ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (13, 1, 'uid', 'ldapx_persons.uid', NULL, 'ldapx_persons', 'ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (14, 1, 'adLogin', 'ldapx_ad_login.login', NULL, 'ldapx_persons,ldapx_ad_login', 'ldapx_ad_login.pers_id=ldapx_persons.uid and ldapx_persons.lang=0', NULL, NULL, 3, 0);
INSERT INTO ldap_attr_mappings VALUES (15, 1, 'adDomain', 'ldapx_ad_login.domain', NULL, 'ldapx_persons,ldapx_ad_login', 'ldapx_ad_login.pers_id=ldapx_persons.uid and ldapx_persons.lang=0', NULL, NULL, 3, 0);

INSERT INTO ldap_oc_mappings VALUES (1, 'inetOrgPerson', 'ldapx_persons', 'id', 'SELECT create_person()', 'DELETE FROM ldapx_persons WHERE id=?', 0);
INSERT INTO ldap_oc_mappings VALUES (3, 'organization', 'ldapx_institutes', 'id', 'SELECT create_o()', 'DELETE FROM ldapx_institutes WHERE id=?', 0);

INSERT INTO ldap_entry_objclasses VALUES (1, 'dcObject');




XYZInsertIntoXYZ




ALTER TABLE ONLY ldap_attr_mappings
		ADD CONSTRAINT ldap_attr_mappings_oc_map_id_fkey	FOREIGN KEY (oc_map_id)	REFERENCES ldap_oc_mappings(id);

ALTER TABLE ONLY ldap_entry_objclasses
		ADD CONSTRAINT ldap_entry_objclasses_entry_id_fkey	FOREIGN KEY (entry_id)	REFERENCES ldap_entries(id);

ALTER TABLE ONLY ldap_entries
		ADD CONSTRAINT ldap_entries_oc_map_id_fkey		FOREIGN KEY (oc_map_id)	REFERENCES ldap_oc_mappings(id);




CREATE SEQUENCE ldap_entries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE ldap_entries_id_seq OWNED BY ldap_entries.id;
ALTER TABLE ONLY ldap_entries ALTER COLUMN id SET DEFAULT nextval('ldap_entries_id_seq'::regclass);
SELECT pg_catalog.setval('ldap_entries_id_seq', 10, true);




CREATE SEQUENCE ldapx_institutes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE ldapx_institutes_id_seq OWNED BY ldapx_institutes.id;
ALTER TABLE ONLY ldapx_institutes ALTER COLUMN id SET DEFAULT nextval('ldapx_institutes_id_seq'::regclass);
SELECT pg_catalog.setval('ldapx_institutes_id_seq', 10, true);




CREATE SEQUENCE ldapx_phones_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE ldapx_phones_id_seq OWNED BY ldapx_phones.id;
ALTER TABLE ONLY ldapx_phones ALTER COLUMN id SET DEFAULT nextval('ldapx_phones_id_seq'::regclass);
SELECT pg_catalog.setval('ldapx_phones_id_seq', 10, true);




CREATE SEQUENCE ldapx_persons_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE ldapx_persons_id_seq OWNED BY ldapx_persons.id;
ALTER TABLE ONLY ldapx_persons ALTER COLUMN id SET DEFAULT nextval('ldapx_persons_id_seq'::regclass);
SELECT pg_catalog.setval('ldapx_persons_id_seq', 10, true);




CREATE SEQUENCE ldapx_mail_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE ldapx_mail_id_seq OWNED BY ldapx_mail.id;
ALTER TABLE ONLY ldapx_mail ALTER COLUMN id SET DEFAULT nextval('ldapx_mail_id_seq'::regclass);
SELECT pg_catalog.setval('ldapx_mail_id_seq', 10, true);




CREATE SEQUENCE ldapx_ad_login_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE ldapx_ad_login_id_seq OWNED BY ldapx_ad_login.id;
ALTER TABLE ONLY ldapx_ad_login ALTER COLUMN id SET DEFAULT nextval('ldapx_ad_login_id_seq'::regclass);
SELECT pg_catalog.setval('ldapx_ad_login_id_seq', 10, true);

`)

	PG_QUE_LDAP_ORGS1 = string(`
delete
	from ldapx_institutes where uid in
		(select uid from ldapx_institutes where uid not in (select uid from XYZDBOrgsXYZ where uid is not null) and uid is not null) and
		id>XYZGlbParXYZ and pass=1;
insert into ldapx_institutes (uid,pass) select distinct uid,1 from XYZDBOrgsXYZ where uid is not null and uid not in
		(select uid from ldapx_institutes where uid is not null);
update ldapx_institutes set idparent='XYZGlbParXYZ', name=subq.name
	from (select name,uid from XYZDBOrgsXYZ) as subq where ldapx_institutes.uid=subq.uid and
		(ldapx_institutes.name<>subq.name or ldapx_institutes.name is NULL or ldapx_institutes.idparent<>'XYZGlbParXYZ') and
		id>XYZGlbParXYZ and pass=1;




insert into ldap_entries (dn,oc_map_id,parent,keyval,uid,idparent,pass) select
		format('OU=%s,XYZGlbDNXYZ', regexp_replace(XYZDBOrgsXYZ.nametr,'[^A-Za-z0-9\ \_\-]', '', 'g')),3,XYZGlbParXYZ,
			ldapx_institutes.id,XYZDBOrgsXYZ.uid,'XYZGlbParXYZ',1 from XYZDBOrgsXYZ,ldapx_institutes
		where XYZDBOrgsXYZ.uid is not null and XYZDBOrgsXYZ.uid not in (select uid from ldap_entries where uid is not null) and
		ldapx_institutes.uid=XYZDBOrgsXYZ.uid;
delete from ldap_entries where uid in (select uid from ldap_entries where uid not in (select uid from XYZDBOrgsXYZ where uid is not null) and
		uid is not null) and id>XYZGlbParXYZ and pass=1;
update ldap_entries set dn=format('OU=%s,XYZGlbDNXYZ', regexp_replace(subq.nametr,'[^A-Za-z0-9\ \_\-]', '', 'g')), oc_map_id=3, parent=XYZGlbParXYZ,
			idparent='XYZGlbParXYZ'
	from (select nametr, uid from XYZDBOrgsXYZ) as subq
		where ldap_entries.uid=subq.uid and (ldap_entries.dn<>format('OU=%s,XYZGlbDNXYZ',
			regexp_replace(subq.nametr,'[^A-Za-z0-9\ \_\-]', '', 'g')) or ldap_entries.dn is NULL or ldap_entries.idparent<>'XYZGlbParXYZ') and
		id>XYZGlbParXYZ and pass=1;




delete
	from ldapx_institutes where uid in
		(select uid from ldapx_institutes where uid not in (select uid from XYZDBDepsXYZ where uid is not null) and uid is not null) and
		id>XYZGlbParXYZ and pass=2;
insert into ldapx_institutes (uid,pass) select distinct uid,2
	from XYZDBDepsXYZ where uid is not null and uid not in (select uid from ldapx_institutes where uid is not null);
update ldapx_institutes set idparent=subq.idparent, name=subq.name
	from (select name,uid,idparent from XYZDBDepsXYZ) as subq
		where ldapx_institutes.uid=subq.uid and (ldapx_institutes.name<>subq.name or ldapx_institutes.name is NULL or
			ldapx_institutes.idparent<>subq.idparent) and id>XYZGlbParXYZ and pass=2;




delete
	from ldap_entries where uid in
		(select uid from ldap_entries where uid not in (select uid from XYZDBDepsXYZ where uid is not null) and
		uid is not null) and id>XYZGlbParXYZ and pass=2;
`)

	PG_QUE_LDAP_ORGS1X_GET = string("select count(uid) from XYZDBDepsXYZ where uid not in (select uid from ldap_entries where uid is not null);")

	PG_QUE_LDAP_ORGS1X_PUT = string(`
insert into ldap_entries (dn,oc_map_id,parent,keyval,uid,idparent,pass)
		select format('OU=%s,%s', regexp_replace(XYZDBDepsXYZ.nametr, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn),
		3, ldap_entries.id, ldapx_institutes.id, XYZDBDepsXYZ.uid, XYZDBDepsXYZ.idparent, 2
			from XYZDBDepsXYZ, ldapx_institutes, ldap_entries
				where XYZDBDepsXYZ.uid is not null and XYZDBDepsXYZ.uid not in
				(select uid from ldap_entries where uid is not null) and
				ldapx_institutes.uid=XYZDBDepsXYZ.uid and
				ldap_entries.uid=XYZDBDepsXYZ.idparent;
`)

	PG_QUE_LDAP_ORGS1_END = string(`
update ldap_entries
		set dn=subq.dn, parent=subq.parent, idparent=subq.idparent
		from (select format('OU=%s,%s', regexp_replace(XYZDBDepsXYZ.nametr, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn) as dn,
				ldap_entries.id as parent, XYZDBDepsXYZ.uid as uid, XYZDBDepsXYZ.idparent as idparent, curr.uid as realuid
			from XYZDBDepsXYZ, ldapx_institutes, ldap_entries, ldap_entries as curr
			where ldapx_institutes.uid=XYZDBDepsXYZ.uid and ldap_entries.uid=XYZDBDepsXYZ.idparent and curr.uid=XYZDBDepsXYZ.uid and
				(curr.dn<>format('OU=%s,%s', regexp_replace(XYZDBDepsXYZ.nametr, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn) or
				curr.parent<>ldap_entries.id or curr.idparent<>XYZDBDepsXYZ.idparent)) as subq
		where ldap_entries.uid=subq.uid and ldap_entries.pass=2;
`)

	PG_QUE_LDAP_PERS1 = string(`
delete from ldapx_persons where lang=0 and uid not in (select uid from XYZDBPersXYZ where uid is not null) and uid is not null;

insert into ldapx_persons (uid,lang)
	select distinct uid,0 from XYZDBPersXYZ where uid is not null and uid not in (select uid from ldapx_persons where uid is not null and lang=0);

update ldapx_persons set idparent=subq.idparent, fullname=format('%s %s %s',subq.surname,subq.name,subq.mn), name=subq.name,
			surname=subq.surname, mn=subq.mn, bc=subq.pos, contract=subq.contract, cid_name=format('%s %s.%s.',subq.surname,subq.fni,subq.mni)
	from (select nfr as name, nlr as surname, nmr as mn, nmir as mni, nfir as fni, uid, idparent, pos, contract from XYZDBPersXYZ) as subq
		where ldapx_persons.lang=0 and ldapx_persons.uid=subq.uid and (ldapx_persons.fullname<>format('%s %s %s',subq.surname,subq.name,subq.mni) or ldapx_persons.fullname is NULL or
			ldapx_persons.idparent<>subq.idparent or ldapx_persons.bc<>subq.pos or ldapx_persons.contract<>subq.contract or cid_name<>format('%s %s.%s.',subq.surname,subq.fni,subq.mni));




delete from ldapx_persons where lang=1 and uid not in (select uid from XYZDBPersXYZ where uid is not null) and uid is not null;

insert into ldapx_persons (id,uid,lang)
	select subq.id,XYZDBPersXYZ.uid,1 from XYZDBPersXYZ, (select id,uid from ldapx_persons) as subq
	where XYZDBPersXYZ.uid is not null and XYZDBPersXYZ.uid not in (select uid from ldapx_persons where uid is not null and lang=1) and
	XYZDBPersXYZ.uid in (select uid from ldapx_persons where uid is not null and lang=0) and XYZDBPersXYZ.uid=subq.uid;

update ldapx_persons set idparent=subq.idparent, fullname=format('%s %s %s',subq.surname,subq.name,subq.mn), name=subq.name,
			surname=subq.surname, mn=subq.mn, bc='-', contract=subq.contract, cid_name=format('%s %s.%s.',subq.surname,subq.fni,subq.mni)
	from (select nft as name, nlt as surname, nmt as mn, nmit as mni, nfit as fni, uid, idparent, contract from XYZDBPersXYZ) as subq
		where ldapx_persons.lang=1 and ldapx_persons.uid=subq.uid and (ldapx_persons.fullname<>format('%s %s %s',subq.surname,subq.name,subq.mni) or ldapx_persons.fullname is NULL or
			ldapx_persons.idparent<>subq.idparent or ldapx_persons.contract<>subq.contract or cid_name<>format('%s %s.%s.',subq.surname,subq.fni,subq.mni));




delete
	from ldap_entries where uid not in (select uid from XYZDBPersXYZ where uid is not null) and
		uid is not null and id>XYZGlbParXYZ and pass=3;
`)

	PG_QUE_LDAP_PERS1X_GET = string("select count(uid) from XYZDBPersXYZ where uid not in (select uid from ldap_entries where uid is not null);")

	PG_QUE_LDAP_PERS1X_PUT = string(`
insert into ldap_entries (dn,oc_map_id,parent,keyval,uid,idparent,pass)
		select format('CN=%s,%s', format('%s %s %s',XYZDBPersXYZ.nlt,XYZDBPersXYZ.nft,XYZDBPersXYZ.nmt), ldap_entries.dn),
		1, ldap_entries.id, ldapx_persons.id, XYZDBPersXYZ.uid, XYZDBPersXYZ.idparent, 3
			from XYZDBPersXYZ, ldapx_persons, ldap_entries
				where ldapx_persons.lang=0 and XYZDBPersXYZ.uid is not null and XYZDBPersXYZ.uid not in
				(select uid from ldap_entries where uid is not null) and
				ldapx_persons.uid=XYZDBPersXYZ.uid and
				ldap_entries.uid=XYZDBPersXYZ.idparent;
`)

	PG_QUE_LDAP_PERS1_END = string(`
update ldap_entries
		set dn=subq.dn, parent=subq.parent, idparent=subq.idparent
		from (select format('CN=%s,%s', format('%s %s %s',XYZDBPersXYZ.nlt,XYZDBPersXYZ.nft,XYZDBPersXYZ.nmt), ldap_entries.dn) as dn,
				ldap_entries.id as parent, XYZDBPersXYZ.uid as uid, XYZDBPersXYZ.idparent as idparent, curr.uid as realuid
			from XYZDBPersXYZ, ldapx_persons, ldap_entries, ldap_entries as curr
			where ldapx_persons.lang=0 and ldapx_persons.uid=XYZDBPersXYZ.uid and ldap_entries.uid=XYZDBPersXYZ.idparent and
				curr.uid=XYZDBPersXYZ.uid and
				(curr.dn<>format('CN=%s,%s', format('%s %s %s',XYZDBPersXYZ.nlt,XYZDBPersXYZ.nft,XYZDBPersXYZ.nmt), ldap_entries.dn) or
				curr.parent<>ldap_entries.id or curr.idparent<>XYZDBPersXYZ.idparent)) as subq
		where ldap_entries.uid=subq.uid and ldap_entries.pass=3;
`)

	PG_QUE_LDAP_PHONES = [][]string{{`
delete from ldapx_phones as ph where ph.pers_id not in (select cache.uid from XYZDBPhonesXYZ as cache where cache.uid=ph.pers_id and cache.phone=ph.phone);
    `, " purge old phones"},
		{`
insert into ldapx_phones (phone,pers_id, pass)
	select cache.phone, cache.uid, cache.type from XYZDBPhonesXYZ as cache
 		where cache.uid not in (select pers_id from ldapx_phones as ph where ph.pers_id=cache.uid and ph.phone=cache.phone) and
        ((cache.type=1 and cache.tm='Y') or (cache.type<>1));
    `, "insert new phones"},
		{`
delete from ldapx_mail where mail not in (select ml.mail from XYZDBMailXYZ as ml, XYZDBPersXYZ as cache, ldapx_mail as mail where mail.mail=ml.mail and mail.pers_id=cache.uid);
    `, " purge old e-mail's"},
		{`
insert into ldapx_mail (mail,pers_id)
	select ml.mail, cache.uid from XYZDBMailXYZ as ml, XYZDBPersXYZ as cache
		where
			(format('%s %s %s', cache.nlr, cache.nfr, cache.nmir)=ml.namerus or format('%s %s %s', cache.nlr, cache.nfr, cache.nmr)=ml.namerus or
				format('%s %s', cache.nlr, cache.nfr)=ml.namerus) and
			ml.mail not in (select ml_ch.mail from ldapx_mail as ml_ch where ml_ch.mail=ml.mail);
	`, "insert new e-mail's"},
		{`
delete from ldapx_ad_login where dlogin not in (select ad.dlogin from XYZDBADXYZ as ad, ldapx_ad_login as login where login.dlogin=ad.dlogin);
    `, " purge old ad-logins"},
		{`
insert into ldapx_ad_login (domain,dlogin,login,pers_id)
	select ad.domain, ad.dlogin, ad.login, cache.uid from XYZDBADXYZ as ad, XYZDBPersXYZ as cache
		where lower(format('%s %s %s', cache.nlr, cache.nfr, cache.nmr))=lower(ad.displayname) and XYZSubParentCheckXYZ and
			ad.dlogin not in (select dlogin from ldapx_ad_login where dlogin=ad.dlogin);
    `, "insert new ad-logins"},
		{`
update XYZDBADXYZ set connected='yes' where lower(displayname) in (select lower(format('%s %s %s', nlr, nfr, nmr)) from XYZDBPersXYZ as cache, ldapx_ad_login as ad where ad.dlogin like '%@%' and XYZSubParentCheckXYZ and ad.pers_id=cache.uid);
	`, "update connected"}}
)
