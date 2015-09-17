package SABDefine

var	(

	Oracle_QUE	=	string (`
select
    z.staff_id, y.adr_l2, y.adr_l1, y.adr_l3, y.afr_comment, y.tm, y.isvisible, x.kind_type
  from
    pbook.pb_kindtel x,
    pbook.pb_cinf y,
    pbook.v_staff z
  where
      x.kind_id=y.kind_id and y.owner_uin=z.fl_id and x.kind_type<4 and owner_uin is not null and isvisible='Y'
`)
)
