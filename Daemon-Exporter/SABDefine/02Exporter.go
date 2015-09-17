package SABDefine

var	(

	LDAP_attr				=	[]string{"altfullname", "cn", "mail"}

	PG_Table_Oracle			=	string("Z_Oracle_X_Cache")
	PG_Table_Oracle_Status	=	string("Z_Oracle_A_status")

	PG_Table_MSSQL			=	[]string{"Z_MSSQL_ORGS_X_Cache", "Z_MSSQL_DEPS_X_Cache", "Z_MSSQL_PERS_X_Cache"}
	PG_Table_MSSQL_Status	=	string("Z_MSSQL_A_status")

	PG_Table_Domino			=	string("Z_Domino_X_Cache")
	PG_Table_Domino_Status	=	string("Z_Domino_A_status")

)
