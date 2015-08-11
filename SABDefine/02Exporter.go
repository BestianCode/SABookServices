package SABDefine

type	Config_STR	struct {
	Oracle_user	string
	Oracle_pass	string
	Oracle_sid	string

	PG_DSN		string

	LDAP_URL	string
	LDAP_User	string
	LDAP_Pass	string
	LDAP_BASE	string
	LDAP_Filter	string

	ROOT_OU		string

	LOG_File	string

}

var	(

	Conf			Config_STR
	LDAP_attr	=	[]string{"altfullname", "cn", "mail"}

	PG_Table_Oracle	=	[]string{"ZZDMP_Ora_ORGS", "ZZDMP_Ora_DEPS", "ZZDMP_Ora_PERS"}
	PG_Table_Domino	=	string("ZZDMP_Domino")

)