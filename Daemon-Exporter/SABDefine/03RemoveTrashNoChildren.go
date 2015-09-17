package SABDefine

var	(

	PG_QUE_RemoveBlackListed	=	string(`
delete from XYZDBOrgsXYZ where uid='XYZUidXYZ';
delete from XYZDBDepsXYZ where uid='XYZUidXYZ' or idparent='XYZUidXYZ' or idorg='XYZUidXYZ';
`)

	PG_QUE_RemoveNoChildren	=	[]string {`
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

