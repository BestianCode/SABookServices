package SABDefine

var	(

	MSSQL_QUE	=	[]string {`
SELECT	dbo._Reference105._IDRRef as idorg, dbo._Reference105._Description AS org, dbo._Reference105._Fld1559RRef as idparent FROM dbo._Reference105;
`,`
SELECT	dbo._Reference119._IDRRef AS idpodr, dbo._Reference119._OwnerIDRRef AS idorg, dbo._Reference119._ParentIDRRef AS idparent, dbo._Reference119._Description AS podr FROM dbo._Reference119;
`,`
SELECT	dbo._Reference22602._IDRRef AS uid, dbo._Reference22602._Description AS name,
	dbo._Reference22602._Code AS tab, dbo._Reference105._IDRRef as idorg, dbo._Reference119._IDRRef AS idparent,
	dbo._Reference50._Description AS position,
	_EnumOrder AS contract
FROM	dbo._Reference22602 LEFT OUTER JOIN
	dbo._Reference105 ON dbo._Reference22602._Fld27474RRef = dbo._Reference105._IDRRef LEFT OUTER JOIN
	dbo._Reference119 ON dbo._Reference22602._Fld27475RRef = dbo._Reference119._IDRRef LEFT OUTER JOIN
	dbo._Reference50 ON dbo._Reference22602._Fld27476RRef = dbo._Reference50._IDRRef INNER JOIN
	dbo._enum576 ON dbo._Reference22602._Fld22809RRef=dbo._enum576._IDRRef
WHERE	(dbo._Reference22602._Fld27477 <> CONVERT(DATETIME, '01.01.1753', 104)) AND (dbo._Reference22602._Fld27478 = CONVERT(DATETIME, '01.01.1753', 104));
`}

)
