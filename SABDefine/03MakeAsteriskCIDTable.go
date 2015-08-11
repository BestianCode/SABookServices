package SABDefine

var	(
	AsteriskCIDTable	=	string ("Asterisk_CID")
	AsteriskCIDTableTemp	=	string ("Asterisk_CIDtemp")

	PG_QUE			=	string (`
drop table if exists XYZTempTableZYX;
CREATE TEMP TABLE XYZTempTableZYX ( name character varying(255), number character varying(255));
insert into XYZTempTableZYX (name, number) select trnamelf, format('8%s', regexp_split_to_table(regexp_replace(regexp_replace(phoneint, '[^0-9\n]', '', 'g'), '\n', ',' ,'g'), ',')) from XYZOraclePersTableZYX where phoneint similar to '%[0-9]%';
insert into XYZWorkTableZYX (number) select distinct number from XYZTempTableZYX where number not in (select number from XYZWorkTableZYX) and length(number)>4 and length(number)<8 order by number;
delete from XYZWorkTableZYX where number in (select number from XYZWorkTableZYX where number not in (select number from XYZTempTableZYX where length(number)>4 and length(number)<8));
update XYZWorkTableZYX set name=subq.name from (select name,number from XYZTempTableZYX) as subq where XYZWorkTableZYX.number=subq.number and (XYZWorkTableZYX.name<>subq.name or XYZWorkTableZYX.name is NULL);
drop table if exists XYZTempTableZYX;
`)

)
