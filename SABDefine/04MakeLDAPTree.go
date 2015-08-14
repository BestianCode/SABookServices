package SABDefine

var	(

	GlobaParentId	=	string("2")

	PG_QUE_LDAP_ORGS1	=	string (`
delete
	from ldapx_institutes where uid in
		(select uid from ldapx_institutes where uid not in (select uid from zzdmp_ora_orgs where uid is not null) and uid is not null) and
		id>XYZGlbParXYZ and pass=1;
insert into ldapx_institutes (uid,pass) select distinct uid,1 from zzdmp_ora_orgs where uid is not null and uid not in
		(select uid from ldapx_institutes where uid is not null);
update ldapx_institutes set idparent='XYZGlbParXYZ', name=subq.name
	from (select name,uid from zzdmp_ora_orgs) as subq where ldapx_institutes.uid=subq.uid and
		(ldapx_institutes.name<>subq.name or ldapx_institutes.name is NULL or ldapx_institutes.idparent<>'XYZGlbParXYZ') and
		id>XYZGlbParXYZ and pass=1;
`)

	PG_QUE_LDAP_ORGS2	=	[]string {`
select dn from ldap_entries where id=XYZGlbParXYZ;
`,`
insert into ldap_entries (dn,oc_map_id,parent,keyval,uid,idparent,pass) select
		format('OU=%s,XYZGlbDNXYZ', regexp_replace(zzdmp_ora_orgs.trname,'[^A-Za-z0-9\ \_\-]', '', 'g')),3,XYZGlbParXYZ,
			ldapx_institutes.id,zzdmp_ora_orgs.uid,'XYZGlbParXYZ',1 from zzdmp_ora_orgs,ldapx_institutes
		where zzdmp_ora_orgs.uid is not null and zzdmp_ora_orgs.uid not in (select uid from ldap_entries where uid is not null) and
		ldapx_institutes.uid=zzdmp_ora_orgs.uid;
`,`
delete from ldap_entries where uid in (select uid from ldap_entries where uid not in (select uid from zzdmp_ora_orgs where uid is not null) and
		uid is not null) and id>XYZGlbParXYZ and pass=1;
`,`
update ldap_entries set dn=format('OU=%s,XYZGlbDNXYZ', regexp_replace(subq.trname,'[^A-Za-z0-9\ \_\-]', '', 'g')), oc_map_id=3, parent=XYZGlbParXYZ,
			idparent='XYZGlbParXYZ'
	from (select trname, uid from zzdmp_ora_orgs) as subq
		where ldap_entries.uid=subq.uid and (ldap_entries.dn<>format('OU=%s,ou=Quadra,o=Enterprise',
			regexp_replace(subq.trname,'[^A-Za-z0-9\ \_\-]', '', 'g')) or ldap_entries.dn is NULL or ldap_entries.idparent<>'XYZGlbParXYZ') and
		id>XYZGlbParXYZ and pass=1;
`}

	PG_QUE_LDAP_ORGS11	=	string (`
delete
	from ldapx_institutes where uid in
		(select uid from ldapx_institutes where uid not in (select uid from zzdmp_ora_deps where uid is not null) and uid is not null) and
		id>XYZGlbParXYZ and pass=2;
insert into ldapx_institutes (uid,pass) select distinct uid,2
	from zzdmp_ora_deps where uid is not null and uid not in (select uid from ldapx_institutes where uid is not null);
update ldapx_institutes set idparent=subq.idparent, name=subq.name
	from (select name,uid,idparent from zzdmp_ora_deps) as subq
		where ldapx_institutes.uid=subq.uid and (ldapx_institutes.name<>subq.name or ldapx_institutes.name is NULL or
			ldapx_institutes.idparent<>subq.idparent) and id>XYZGlbParXYZ and pass=2;
`)

	PG_QUE_LDAP_ORGS12	=	[]string {`
delete
	from ldap_entries where uid in
		(select uid from ldap_entries where uid not in (select uid from zzdmp_ora_deps where uid is not null) and
		uid is not null) and id>XYZGlbParXYZ and pass=2;
`,`
select count(uid) from zzdmp_ora_deps where uid not in (select uid from ldap_entries where uid is not null);
`,`
insert into ldap_entries (dn,oc_map_id,parent,keyval,uid,idparent,pass)
		select format('OU=%s,%s', regexp_replace(zzdmp_ora_deps.trname, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn),
		3, ldap_entries.id, ldapx_institutes.id, zzdmp_ora_deps.uid, zzdmp_ora_deps.idparent, 2
			from zzdmp_ora_deps, ldapx_institutes, ldap_entries
				where zzdmp_ora_deps.uid is not null and zzdmp_ora_deps.uid not in
				(select uid from ldap_entries where uid is not null) and
				ldapx_institutes.uid=zzdmp_ora_deps.uid and
				ldap_entries.uid=zzdmp_ora_deps.idparent;
`,`
update ldap_entries
		set dn=subq.dn, parent=subq.parent, idparent=subq.idparent
		from (select format('OU=%s,%s', regexp_replace(zzdmp_ora_deps.trname, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn) as dn,
				ldap_entries.id as parent, zzdmp_ora_deps.uid as uid, zzdmp_ora_deps.idparent as idparent, curr.uid as realuid
			from zzdmp_ora_deps, ldapx_institutes, ldap_entries, ldap_entries as curr
			where ldapx_institutes.uid=zzdmp_ora_deps.uid and ldap_entries.uid=zzdmp_ora_deps.idparent and curr.uid=zzdmp_ora_deps.uid and
				(curr.dn<>format('OU=%s,%s', regexp_replace(zzdmp_ora_deps.trname, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn) or
				curr.parent<>ldap_entries.id or curr.idparent<>zzdmp_ora_deps.idparent)) as subq
		where ldap_entries.uid=subq.uid and ldap_entries.pass=2;
`}

	PG_QUE_LDAP_ORGS21	=	string (`
delete
	from ldapx_persons where uid in
		(select uid from ldapx_persons where uid not in (select uid from zzdmp_ora_pers where uid is not null) and uid is not null);
insert into ldapx_persons (uid)
	select distinct uid from zzdmp_ora_pers where uid is not null and uid not in (select uid from ldapx_persons where uid is not null);
update ldapx_persons set idparent=subq.idparent, fullname=subq.fullname, name=regexp_replace(subq.name, '.*\ ', '', 'g'), surname=regexp_replace(subq.name, '\ .*', '', 'g')
	from (select namefull as fullname, namelf as name, uid, iddep as idparent from zzdmp_ora_pers) as subq
		where ldapx_persons.uid=subq.uid and (ldapx_persons.fullname<>subq.fullname or ldapx_persons.fullname is NULL or
			ldapx_persons.idparent<>subq.idparent);
`)

	PG_QUE_LDAP_ORGS22	=	[]string {`
delete
	from ldap_entries where uid in
		(select uid from ldap_entries where uid not in (select uid from zzdmp_ora_pers where uid is not null) and
		uid is not null) and id>XYZGlbParXYZ and pass=3;
`,`
select count(uid) from zzdmp_ora_pers where uid not in (select uid from ldap_entries where uid is not null);
`,`
insert into ldap_entries (dn,oc_map_id,parent,keyval,uid,idparent,pass)
		select format('CN=%s,%s', regexp_replace(zzdmp_ora_pers.trnamefull, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn),
		1, ldap_entries.id, ldapx_persons.id, zzdmp_ora_pers.uid, zzdmp_ora_pers.iddep, 3
			from zzdmp_ora_pers, ldapx_persons, ldap_entries
				where zzdmp_ora_pers.uid is not null and zzdmp_ora_pers.uid not in
				(select uid from ldap_entries where uid is not null) and
				ldapx_persons.uid=zzdmp_ora_pers.uid and
				ldap_entries.uid=zzdmp_ora_pers.iddep;
`,`
update ldap_entries
		set dn=subq.dn, parent=subq.parent, idparent=subq.idparent
		from (select format('CN=%s,%s', regexp_replace(zzdmp_ora_pers.trnamefull, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn) as dn,
				ldap_entries.id as parent, zzdmp_ora_pers.uid as uid, zzdmp_ora_pers.iddep as idparent, curr.uid as realuid
			from zzdmp_ora_pers, ldapx_persons, ldap_entries, ldap_entries as curr
			where ldapx_persons.uid=zzdmp_ora_pers.uid and ldap_entries.uid=zzdmp_ora_pers.iddep and curr.uid=zzdmp_ora_pers.uid and
				(curr.dn<>format('CN=%s,%s', regexp_replace(zzdmp_ora_pers.trnamefull, '[^A-Za-z0-9\ \_\-]', '', 'g'), ldap_entries.dn) or
				curr.parent<>ldap_entries.id or curr.idparent<>zzdmp_ora_pers.iddep)) as subq
		where ldap_entries.uid=subq.uid and ldap_entries.pass=3;
`}

	PG_QUE_LDAP_ORGS31	=	string (`
drop table if exists XYZTempTableZYX;

CREATE TEMP TABLE XYZTempTableZYX ( phone character varying(255), pers_id integer);

insert into XYZTempTableZYX (phone, pers_id)
        select format('8%s', regexp_split_to_table(regexp_replace(regexp_replace(ora.phoneint, '[^0-9\n]', '', 'g'), '\n', ',' ,'g'), ',')), pers.id
                from zzdmp_ora_pers as ora, ldapx_persons as pers where pers.uid=ora.uid and ora.phoneint similar to '%[0-9]%';
insert into ldapx_phones_work (phone,pers_id, pass)
	select phone, pers_id,1 from XYZTempTableZYX as tmp where length(phone)>4 and phone not in
		(select phone from ldapx_phones_work as tst where tst.phone=tmp.phone and tst.pers_id=tmp.pers_id and pass=1);

drop table if exists XYZTempTableZYX;

CREATE TEMP TABLE XYZTempTableZYX ( phone character varying(255), pers_id integer);

insert into XYZTempTableZYX (phone, pers_id)
        select regexp_split_to_table(regexp_replace(regexp_replace(ora.phonetown, '[^0-9доб.,\(\)\-\+\n]', '', 'g'), '\n', ',' ,'g'), ','), pers.id
                from zzdmp_ora_pers as ora, ldapx_persons as pers where pers.uid=ora.uid and ora.phonetown similar to '%[0-9]%';
insert into ldapx_phones_work (phone,pers_id, pass)
	select phone, pers_id,2 from XYZTempTableZYX as tmp where length(phone)>4 and phone not in
		(select phone from ldapx_phones_work as tst where tst.phone=tmp.phone and tst.pers_id=tmp.pers_id and pass=2);

drop table if exists XYZTempTableZYX;

CREATE TEMP TABLE XYZTempTableZYX ( phone character varying(255), pers_id integer);

insert into XYZTempTableZYX (phone,pers_id)
        select regexp_split_to_table(regexp_replace(regexp_replace(ora.phonecell, '[^+0-9\n]', '', 'g'), '\n', ',' ,'g'), ','), pers.id
                from zzdmp_ora_pers as ora, ldapx_persons as pers where pers.uid=ora.uid and ora.phonecell similar to '%[0-9]%';
insert into ldapx_phones_mobile (phone,pers_id) 
        select phone, pers_id from XYZTempTableZYX as tmp where length(phone)>4 and phone not in
                (select phone from ldapx_phones_mobile as tst where tst.phone=tmp.phone and tst.pers_id=tmp.pers_id);

drop table if exists XYZTempTableZYX;
`)

)

