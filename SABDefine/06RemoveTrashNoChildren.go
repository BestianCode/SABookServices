package SABDefine

var	(

	PG_QUE_RemoveNoChildren	=	[]string {`
select count(uid) from zzdmp_ora_orgs where
		uid not in (select idparent from zzdmp_ora_deps) and
		uid not in (select iddep from zzdmp_ora_pers) and
		uid not in (select idorg from zzdmp_ora_pers);
`,`
delete from zzdmp_ora_orgs where uid in (select uid from zzdmp_ora_orgs where
                uid not in (select idparent from zzdmp_ora_deps) and
                uid not in (select iddep from zzdmp_ora_pers) and
                uid not in (select idorg from zzdmp_ora_pers));
`,`
select count(uid) from zzdmp_ora_deps where
                uid not in (select idparent from zzdmp_ora_deps) and
                uid not in (select iddep from zzdmp_ora_pers) and
                uid not in (select idorg from zzdmp_ora_pers);
`,`
delete from zzdmp_ora_deps where uid in (select uid from zzdmp_ora_deps where
                uid not in (select idparent from zzdmp_ora_deps) and
                uid not in (select iddep from zzdmp_ora_pers) and
                uid not in (select idorg from zzdmp_ora_pers));
`,`
select count(uid) from ldapx_institutes where uid not in (select idparent from ldapx_institutes) and pass>0 and pass<3;
`,`
delete from ldapx_institutes where uid in (select uid from ldapx_institutes where uid not in (select idparent from ldap_entries) and pass>0 and pass<3);
`,`
select count(uid) from ldap_entries where uid not in (select idparent from ldap_entries) and pass>0 and pass<3;
`,`
delete from ldap_entries where uid in (select uid from ldap_entries where uid not in (select idparent from ldap_entries) and pass>0 and pass<3);
`}

)

