package SABDefine

var (
	PG_QUE_AD_GetDupInAD = string(`
select dn from
		XYZDBADXYZ as ad,
		(select cn from XYZDBADXYZ group by cn having count(cn)>1) as cn_mult,
		(select displayname from XYZDBADXYZ group by displayname having count(displayname)>1) as dn_mult
	where
		cn_mult.cn=ad.cn or dn_mult.displayname=ad.displayname group by dn order by dn;
    `)

	PG_QUE_AD_GetUpdateAD = string(`
select format('%s %s %s', cache.nlr, cache.nfr, cache.nmr), xad.dn from XYZDBADXYZ as xad, XYZDBPersXYZ as cache
    where
        (
            lower(format('%s %s', cache.nlr, cache.nfr))=lower(xad.displayname) or
            lower(format('%s %s', cache.nlr, cache.nfr))=lower(xad.cn) or
            lower(format('%s %s', cache.nfr, cache.nlr))=lower(xad.displayname) or
            lower(format('%s %s', cache.nfr, cache.nlr))=lower(xad.cn)
        ) and XYZSubParentCheckXYZ and 
        dn not in
            (select dn from XYZDBADXYZ as ad,
                    (select cn from XYZDBADXYZ group by cn having count(cn)>1) as cn_mult,
                    (select displayname from XYZDBADXYZ group by displayname having count(displayname)>1) as dn_mult
                where lower(cn_mult.cn)=lower(ad.cn) or lower(dn_mult.displayname)=lower(ad.displayname) group by dn order by dn)
        and lower(format('%s %s', cache.nlr, cache.nfr)) not in
            (select lower(format('%s %s', cache.nlr, cache.nfr)) from XYZDBPersXYZ as cache group by lower(format('%s %s', cache.nlr, cache.nfr)) having count(format('%s %s', cache.nlr, cache.nfr))>1)
    and xad.connected<>lower('yes')
    order by lower(format('%s %s %s', cache.nlr, cache.nfr, cache.nmr));
    `)

	PG_QUE_AD_GetNotConnected = string(`
select dn from XYZDBADXYZ where dn not in 
    (select xad.dn from XYZDBADXYZ as xad, XYZDBPersXYZ as cache
        where
            (
                lower(format('%s %s', cache.nlr, cache.nfr))=lower(xad.displayname) or
                lower(format('%s %s', cache.nlr, cache.nfr))=lower(xad.cn) or
                lower(format('%s %s', cache.nfr, cache.nlr))=lower(xad.displayname) or
                lower(format('%s %s', cache.nfr, cache.nlr))=lower(xad.cn)
            )and XYZSubParentCheckXYZ and
            dn not in
                (select dn from XYZDBADXYZ as ad,
                        (select cn from XYZDBADXYZ group by cn having count(cn)>1) as cn_mult,
                        (select displayname from XYZDBADXYZ group by displayname having count(displayname)>1) as dn_mult
                    where lower(cn_mult.cn)=lower(ad.cn) or lower(dn_mult.displayname)=lower(ad.displayname) group by dn order by dn)
            and lower(format('%s %s', cache.nlr, cache.nfr)) not in
                (select lower(format('%s %s', cache.nlr, cache.nfr)) from XYZDBPersXYZ as cache group by lower(format('%s %s', cache.nlr, cache.nfr)) having count(format('%s %s', cache.nlr, cache.nfr))>1))
    and connected<>lower('yes');
    `)

	PG_QUE_AD_SetCredentInfoToAD = string(`
select fsab.dn, fsab.mail, fsab.title, fsab.ph_mob, fsab.ph_int, fsab.ph_ext from XYZDBADXYZ as fad,
    (select
        x.dn as dn, z.mail as mail, v.bc as title,
        coalesce((select string_agg(phone, ', ') from ldapx_phones where pass=1 and pers_id=y.pers_id), '') as ph_mob,
        coalesce((select string_agg(phone, ', ') from ldapx_phones where pass=2 and pers_id=y.pers_id), '') as ph_int,
        coalesce((select string_agg(phone, ', ') from ldapx_phones where pass=3 and pers_id=y.pers_id), '') as ph_ext
            from XYZDBADXYZ as x, ldapx_ad_login as y, ldapx_mail as z, ldapx_persons as v
                where
                    z.pers_id=y.pers_id and
                    y.dlogin=x.dlogin and
                    v.uid=y.pers_id and
                    v.lang=0) as fsab
    where fad.connected='yes' and fad.dn=fsab.dn and
        (fad.mail<>fsab.mail or
        lower(fad.title)<>lower(fsab.title) or
        fad.ph_int<>fsab.ph_int or
        fad.ph_ip<>fsab.ph_ext or
        fad.ph_mob<>fsab.ph_mob);
    `)

/*	, `
        order by lower(format('%s %s %s', pers.nlr, pers.nfr, pers.nmr));
select format('%s %s', surname, name), count(format('%s %s', surname, name)) from ldapx_persons where lang=0 group by format('%s %s', surname, name) having count(format('%s %s', surname, name))>2;
`, `
select ad.dn, format('%s %s %s', cache.nlr, cache.nfr, cache.nmr) from XYZDBADXYZ as ad, XYZDBPersXYZ as cache
	where
			lower(format('%s %s', cache.nlr, cache.nfr))=lower(ad.displayname) or
			lower(format('%s %s', cache.nlr, cache.nfr))=lower(ad.cn);
`}*/
)

/*
delete from ldapx_ad_login where dlogin not in (select ad.dlogin from XYZDBADXYZ as ad, ldapx_ad_login as login where login.dlogin=ad.dlogin);

insert into ldapx_ad_login (domain,login,pers_id)
	select ad.domain, ad.dlogin, cache.uid from XYZDBADXYZ as ad, XYZDBPersXYZ as cache
		where
			(
				format('%s %s %s', cache.nlr, cache.nfr, cache.nmr)=ad.displayname or
				format('%s %s %s', cache.nlr, cache.nfr, cache.nmr)=ad.cn or
				format('%s %s', cache.nlr, cache.nfr)=ad.displayname or
				format('%s %s', cache.nlr, cache.nfr)=ad.cn
			); and ad.login not in (select ml_ch.mail from ldapx_ad_login as ml_ch where ml_ch.mail=ml.mail and cache.uid=ml_ch.pers_id);
*/
