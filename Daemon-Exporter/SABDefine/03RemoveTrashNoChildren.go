package SABDefine

var	(

	PG_QUE_RemoveBlackListed	=	string(`
delete from XYZDBOrgsXYZ where uid='XYZUidXYZ';
delete from XYZDBDepsXYZ where uid='XYZUidXYZ' or idparent='XYZUidXYZ' or idorg='XYZUidXYZ';
`)

	PG_QUE_RemoveNoChildren	=	[]string {`
select count(x.uid) from XYZDBPersXYZ as x, XYZDBPhonesXYZ as y
	where x.uid not in (select uid from XYZDBPhonesXYZ where x.uid=uid)
		and y.uid not in (select uid from XYZDBPersXYZ where y.uid=uid)
		and upper(format('%s %s %s', x.nlr, x.nfr, x.nmr))=upper(y.fname);
`,`
update XYZDBPhonesXYZ set uid=subq.uidnew
	from (select x.uid as uidnew, y.uid as uidold from XYZDBPersXYZ as x, XYZDBPhonesXYZ as y
		where x.uid not in (select uid from XYZDBPhonesXYZ where x.uid=uid)
			and y.uid not in (select uid from XYZDBPersXYZ where y.uid=uid)
			and upper(format('%s %s %s', x.nlr, x.nfr, x.nmr))=upper(y.fname)) as subq where uid=subq.uidold;
`,`
select count(fname) from XYZDBPhonesXYZ group by uid, phone, type, fname having count(uid)>1;
`,`
delete from XYZDBPhonesXYZ where lower(server)=lower('XSortPhones');
create temp table tmp as select uid, phone, type, fname from XYZDBPhonesXYZ group by uid, phone, type, fname having count(uid)>1;
delete from XYZDBPhonesXYZ as y
	using (select uid, phone, type, fname from tmp group by uid, phone, type, fname) as subq
		where y.fname=subq.fname and y.phone=subq.phone and y.uid=subq.uid and y.type=subq.type;
insert into XYZDBPhonesXYZ (server,uid,phone,comment,tm,visible,type,fname) select 'XSortPhones',uid,phone,'','Z','Y',type,fname from tmp;
drop table tmp;
`,`
select count(uid) from XYZDBOrgsXYZ where
		uid not in (select idparent from XYZDBDepsXYZ) and
		uid not in (select idparent from XYZDBPersXYZ) and
		uid not in (select idorg from XYZDBPersXYZ);
`,`
delete from XYZDBOrgsXYZ where uid in (select uid from XYZDBOrgsXYZ where
		uid not in (select idparent from XYZDBDepsXYZ) and
		uid not in (select idparent from XYZDBPersXYZ) and
		uid not in (select idorg from XYZDBPersXYZ));
`,`
select count(uid) from XYZDBDepsXYZ where
		uid not in (select idparent from XYZDBDepsXYZ) and
		uid not in (select idparent from XYZDBPersXYZ) and
		uid not in (select idorg from XYZDBPersXYZ);
`,`
delete from XYZDBDepsXYZ where uid in (select uid from XYZDBDepsXYZ where
		uid not in (select idparent from XYZDBDepsXYZ) and
		uid not in (select idparent from XYZDBPersXYZ) and
		uid not in (select idorg from XYZDBPersXYZ));
`,`
select count(uid) from XYZDBDepsXYZ where
		idparent not in (select uid from XYZDBDepsXYZ) and
		idparent not in (select uid from XYZDBOrgsXYZ) and
		idorg not in (select uid from XYZDBOrgsXYZ);
`,`
delete from XYZDBDepsXYZ where
		idparent not in (select uid from XYZDBDepsXYZ) and
		idparent not in (select uid from XYZDBOrgsXYZ) and
		idorg not in (select uid from XYZDBOrgsXYZ);
`,`
select count(uid) from XYZDBPersXYZ where
		idparent not in (select uid from XYZDBDepsXYZ) and
		idparent not in (select uid from XYZDBOrgsXYZ);
`,`
delete from XYZDBPersXYZ where
		idparent not in (select uid from XYZDBDepsXYZ) and
		idparent not in (select uid from XYZDBOrgsXYZ);
`}

)

