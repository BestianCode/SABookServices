package SABDefine

var	(

	Oracle_QUE	=	[]string {`SELECT
  DP.ID_DEP as Id
  , DP.DEP_FULL as Name
FROM 
  V_DEP DP
WHERE
  (DP.BRANCH_ID = DP.ID_DEP
    AND (BITAND(DP.DEP_FLAGS,2)>0 OR BITAND(DP.DEP_FLAGS,3)>0 AND BITAND(DP.SYS_FLAGS,64)=0))
ORDER BY Name`,`SELECT
  DP.ID_DEP as Id
  , DP.ID_PARENT as IdParent
  , DP.DEP_FULL as Name
  , DP.DEP_CODE as Code
  , DP.BRANCH_ID as IdRoot
FROM 
  V_DEP DP
WHERE 
  DP.ID_PARENT IS NOT NULL
  AND (DP.VALID_TO IS NULL OR DP.VALID_TO > CURRENT_DATE) and DP.DEP_FULL is not null
ORDER BY DP.DEP_FULL`,`SELECT
  ST.FL_ID as Id
  , ST.ID_DEP as IdDep
  , ST.FL_FNAME as strFIO
  , CASE WHEN INSTR(ST.FL_FNAME, ' ', 1, 2) = 0 THEN ST.FL_FNAME ELSE SUBSTR(ST.FL_FNAME, 1, INSTR(ST.FL_FNAME, ' ', 1, 2) - 1) END as strFI
  , CASE WHEN INSTR(ST.FL_FNAME, ' ', 1, 2) = 0 THEN ST.FL_FNAME ELSE SUBSTR(ST.FL_FNAME, 1, INSTR(ST.FL_FNAME, ' ', 1, 2) + 1) END as strFI_
  , PS.POS_NAME as strPosit
  , PBOOK.GET_CINF(ST.FL_ID, '9F9D24D6BDAA356642184F3D1F3D951D') as strTel
  , PBOOK.GET_CINF(ST.FL_ID, '9E1380A9DE9B45324F38541444A55363') as strTelGTS
  , PBOOK.GET_CINF(ST.FL_ID, 'B5F29F40094CFD744A149000CC9C4191', 'Y') as strTelGSM
  , PBOOK.GET_CINF(ST.FL_ID, 'AECE5F7A7E1209E4462BD733E2C8493E') as strEMail_
  , ST.BRANCH_ID as IdRoot
  , PS.POS_ID
  , CASE WHEN PS.RANG = 0 THEN 9999999 ELSE NVL(PS.RANG, 9999999) END as intNo
FROM
  (
    SELECT 
      ST.FL_ID
      , ST.ID_DEP
      , MAX(ST.FL_FNAME) AS FL_FNAME 
      , MAX(ST.POS_ID) AS POS_ID 
      , MAX(ST.BRANCH_ID) AS BRANCH_ID
    FROM PBOOK.V_STAFF ST
    WHERE 
      ST.OUT_DATE IS NULL OR ST.OUT_DATE > CURRENT_DATE
    GROUP BY FL_ID, ID_DEP
  ) ST
  , PBOOK.V_EMP PS
  , 
  (
    SELECT DP.ID_DEP 
    FROM PBOOK.V_DEP DP
  ) DP
WHERE
  ST.POS_ID = PS.POS_ID(+)
  AND ST.BRANCH_ID = DP.ID_DEP AND (
    REGEXP_LIKE (PBOOK.GET_CINF(ST.FL_ID, '9F9D24D6BDAA356642184F3D1F3D951D'), '([0-9])') OR
    REGEXP_LIKE (PBOOK.GET_CINF(ST.FL_ID, '9E1380A9DE9B45324F38541444A55363'), '([0-9])') OR
    REGEXP_LIKE (PBOOK.GET_CINF(ST.FL_ID, 'B5F29F40094CFD744A149000CC9C4191'), '([0-9])')
  )
ORDER BY Id, IdDep, intNo, strFIO`}

	MSSQL_QUE	=	[]string {`
SELECT
		dbo._Reference22602._IDRRef AS uid,
		dbo._Reference22602._Description AS fio,
		dbo._Reference22602._Code AS kod_fio,
		dbo._Reference105._IDRRef as idorg,
		dbo._Reference105._Description AS org,
		dbo._Reference119._IDRRef AS idpodr,
		dbo._Reference119._Description AS podr,
		dbo._Reference50._IDRRef AS iddolg,
		dbo._Reference50._Description AS dolg,
		dbo._Reference22602._Fld27478 AS uloln, 
		dbo._Reference22602._Fld27477 AS priem,
			CASE WHEN dbo._Reference22602._Fld27478 = CONVERT(DATETIME, '01.01.1753', 104) THEN 'false' ELSE 'true' END AS prizn_uvoln, dbo._Reference22602._Folder
	FROM dbo._Reference22602 LEFT OUTER JOIN
		dbo._Reference105 ON dbo._Reference22602._Fld27474RRef = dbo._Reference105._IDRRef LEFT OUTER JOIN
		dbo._Reference119 ON dbo._Reference22602._Fld27475RRef = dbo._Reference119._IDRRef LEFT OUTER JOIN
		dbo._Reference50 ON dbo._Reference22602._Fld27476RRef = dbo._Reference50._IDRRef
	WHERE (dbo._Reference22602._Folder = 1) AND (dbo._Reference22602._Fld27477 <> CONVERT(DATETIME, '01.01.1753', 104))
`}

)
