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


)
/*
drop table if exists ldap_entries_temp;
CREATE TEMP TABLE ldap_entries_temp (id integer, dn character varying(255), oc_map_id integer, parent integer, keyval integer, uid bytea, idparent bytea, pass integer);
insert into ldap_entries_temp (dn,oc_map_id,parent,keyval,uid,idparent,pass) select
 format('OU=%s,XYZParentXYZ', regexp_replace(zzdmp_ora_deps.trname, '[^A-Za-z0-9\ \_\-]', '', 'g')),3,0,ldapx_institutes.id,zzdmp_ora_deps.uid,zzdmp_ora_deps.idparent,2
 from zzdmp_ora_deps,ldapx_institutes where zzdmp_ora_deps.uid is not null and zzdmp_ora_deps.uid not in (select uid from ldap_entries where uid is not null) and
 ldapx_institutes.uid=zzdmp_ora_deps.uid;

update ldap_entries_temp set parent=subq.id, dn=regexp_replace(subq.olddn, 'XYZParentXYZ', subq.dn, 'g') from (select x1.id as id,x1.dn as dn,x2.dn as olddn,x2.uid as uid from
 ldap_entries as x1, ldap_entries_temp as x2 where x1.uid=x2.idparent) as subq where idparent=subq.uid;
insert into ldap_entries (dn,oc_map_id,parent,keyval,uid,idparent,pass) select dn,oc_map_id,parent,keyval,uid,idparent,pass from ldap_entries_temp where parent<>0;
*/
