package SABDefine

var	(
	AsteriskCIDTable	=	string ("callerid")
	AsteriskCIDTableTemp	=	string ("calleridtemp")

	PG_QUE			=	string (`
drop table if exists XYZZYX;
CREATE TEMP TABLE XYZZYX ( name character varying(255), number character varying(255));
insert into XYZZYX (name, number) select namelf, format('8%s', regexp_split_to_table(regexp_replace(regexp_replace(phoneint, '[^0-9\n]', '', 'g'), '\n', ',' ,'g'), ',')) from dump_oracle_pers where phoneint similar to '%[0-9]%';
insert into callerid (number) select distinct number from XYZZYX where number not in (select number from callerid) and length(number)>4 and length(number)<8 order by number;
delete from callerid where number in (select number from callerid where number not in (select number from XYZZYX where length(number)>4 and length(number)<8));
update callerid set name=subq.name from (select name,number from XYZZYX) as subq where callerid.number=subq.number and callerid.name<>subq.name;
drop table if exists XYZZYX;
`)

)

//insert into callerid (name, number) select name, number from XYZZYX where length(number)>4 and length(number)<8;
//insert into callerid (name,number) select name, distinct number from XYZZYX where number not in (select number from callerid) and length(number)>4 and length(number)<8 order by number;
