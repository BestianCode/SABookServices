package main

var (
	mySQL_InitDB = string(`
CREATE TABLE IF NOT EXISTS addressbooks (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  principaluri varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  description text COLLATE utf8_unicode_ci,
  ctag int(11) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (id),
  UNIQUE KEY principaluri (principaluri,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS calendarobjects (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  calendardata mediumblob,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  calendarid int(10) unsigned NOT NULL,
  lastmodified int(11) unsigned DEFAULT NULL,
  etag varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  size int(11) unsigned NOT NULL,
  componenttype varchar(8) COLLATE utf8_unicode_ci DEFAULT NULL,
  firstoccurence int(11) unsigned DEFAULT NULL,
  lastoccurence int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY calendarid (calendarid,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS calendars (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  principaluri varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  ctag int(10) unsigned NOT NULL DEFAULT '0',
  description text COLLATE utf8_unicode_ci,
  calendarorder int(10) unsigned NOT NULL DEFAULT '0',
  calendarcolor varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  timezone text COLLATE utf8_unicode_ci,
  components varchar(21) COLLATE utf8_unicode_ci DEFAULT NULL,
  transparent tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (id),
  UNIQUE KEY principaluri (principaluri,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS cards (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  addressbookid int(11) unsigned NOT NULL,
  carddata mediumblob,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  lastmodified int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS groupmembers (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  principal_id int(10) unsigned NOT NULL,
  member_id int(10) unsigned NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY principal_id (principal_id,member_id)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS locks (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  owner varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  timeout int(10) unsigned DEFAULT NULL,
  created int(11) DEFAULT NULL,
  token varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  scope tinyint(4) DEFAULT NULL,
  depth tinyint(4) DEFAULT NULL,
  uri varchar(1000) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (id),
  KEY token (token),
  KEY uri (uri(333))
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS principals (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  uri varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  email varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  vcardurl varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uri (uri)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS users (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  username varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  digesta1 varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY username (username)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;





CREATE TABLE IF NOT EXISTS z_cache_users (
  id int(10) unsigned NOT NULL,
  username varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  digesta1 varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  UNIQUE KEY username (username)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS z_cache_principals (
  id  int(10) unsigned NOT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  email varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  vcardurl varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  UNIQUE KEY uri (uri)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS z_cache_cards (
  id int(11) unsigned NOT NULL,
  addressbookid int(11) unsigned NOT NULL,
  carddata mediumblob,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  lastmodified int(11) unsigned DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS z_cache_addressbooks (
  id int(11) unsigned NOT NULL,
  principaluri varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  description text COLLATE utf8_unicode_ci,
  ctag int(11) unsigned NOT NULL DEFAULT '1',
  UNIQUE KEY principaluri (principaluri,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

truncate table z_cache_users;
truncate table z_cache_principals;
truncate table z_cache_cards;
truncate table z_cache_addressbooks;
			`)

	mySQL_Update_full1 = []string{`
delete from users where username in 
	(select * from
		(select a.username from users as a where a.username not in
			(select b.username from z_cache_users as b where b.username=a.username and b.digesta1=a.digesta1)) as c);
`, `
delete from principals where uri in
	(select * from
		(select a.uri from principals as a where a.uri not in
			(select b.uri from z_cache_principals as b where b.uri=a.uri)) as c);
`, `
delete from addressbooks where principaluri in
	(select * from
		(select a.principaluri from addressbooks as a where a.principaluri not in
			(select b.principaluri from z_cache_addressbooks as b where b.principaluri=a.principaluri)) as c);
`, `
delete x,y from
	principals as x join addressbooks y on x.id=y.id
			where x.id not in
				(select * from (select a.id from principals as a,
					(select id, username from users) as subq
						where a.id=subq.id and subq.username=REPLACE(a.uri, 'principals/', '')) as c);
			`}

	mySQL_Update_full2 = []string{`
delete from cards where uri in
	(select * from
		(select a.uri from cards as a where a.uri not in
			(select b.uri from z_cache_cards as b where b.uri=a.uri and a.addressbookid=b.addressbookid)) as c) or
	addressbookid not in (select id from users);
			`}

	mySQL_Update_part1 = []string{`
delete from users where id=XYZIDXYZ and username in 
	(select * from
		(select a.username from users as a where a.id=XYZIDXYZ and a.username not in
			(select b.username from z_cache_users as b where b.id=XYZIDXYZ and b.username=a.username and b.digesta1=a.digesta1)) as c);
`, `
delete from principals where id=XYZIDXYZ and uri in
	(select * from
		(select a.uri from principals as a where a.id=XYZIDXYZ and a.uri not in
			(select b.uri from z_cache_principals as b where a.id=XYZIDXYZ and b.uri=a.uri)) as c);
`, `
delete from addressbooks where id=XYZIDXYZ and principaluri in
	(select * from
		(select a.principaluri from addressbooks as a where a.id=XYZIDXYZ and a.principaluri not in
			(select b.principaluri from z_cache_addressbooks as b where b.id=XYZIDXYZ and b.principaluri=a.principaluri)) as c);
`, `
delete x,y from
	principals as x join addressbooks y on x.id=y.id
			where x.id=XYZIDXYZ and x.id not in
				(select * from (select a.id from principals as a,
					(select id, username from users) as subq
						where a.id=XYZIDXYZ and a.id=subq.id and subq.username=REPLACE(a.uri, 'principals/', '')) as c);
			`}

	mySQL_Update_part2 = []string{`
delete from cards where addressbookid=XYZIDXYZ and uri in
	(select * from
		(select a.uri from cards as a where a.addressbookid=XYZIDXYZ and a.uri not in
			(select b.uri from z_cache_cards as b where b.uri=a.uri and a.addressbookid=b.addressbookid)) as c) or
	addressbookid not in (select id from users);
			`}

	mySQL_Update1 = []string{`
insert into users (id,username,digesta1)
	select a.id, a.username, a.digesta1 from z_cache_users as a
		where a.username not in (select b.username from users as b where b.username=a.username and
			b.digesta1=a.digesta1);
`, `
insert into principals (id,uri)
	select a.id, a.uri from (select id,username from users) as subq, z_cache_principals as a
		where a.uri not in (select b.uri from principals as b where b.uri=a.uri) and
			subq.username=REPLACE(a.uri, 'principals/', '');
`, `
insert into addressbooks (id,principaluri,uri,ctag)
	select subq.id, a.principaluri,'default',1 from (select id,username from users) as subq, z_cache_addressbooks as a
		where a.principaluri not in
			(select b.principaluri from addressbooks as b where b.principaluri=a.principaluri) and
				subq.username=REPLACE(a.principaluri, 'principals/', '');
			`}

	mySQL_Update2 = []string{`
insert into cards (addressbookid, carddata, uri)
	select a.addressbookid, a.carddata, a.uri from z_cache_cards as a
		where a.uri not in
			(select b.uri from cards as b where a.uri=b.uri and a.addressbookid=b.addressbookid);
  			`}
)
