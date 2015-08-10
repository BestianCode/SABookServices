package SABDefine

type	Config_STR	struct {
	Oracle_user	string
	Oracle_pass	string
	Oracle_sid	string

	PG_DSN		string

	PG_Table_LDAP	string

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

	)