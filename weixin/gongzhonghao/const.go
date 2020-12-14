package gongzhonghao

const (
	TokenURL                   = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	GetServerURL               = "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=%s"
	GetTicketStr               = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	weixin_template_url        = "https://api.weixin.qq.com/cgi-bin/message/template/send"
	weixin_user_token          = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	weixin_user_info           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	weixin_userinfo_url        = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
	app_pay_unifiedorder       = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	app_pay_qrcode             = "https://api.mch.weixin.qq.com/pay/micropay"
	app_pay_orderquery         = "https://api.mch.weixin.qq.com/pay/orderquery"
	app_pay_authcodetoopenid   = "https://api.mch.weixin.qq.com/tools/authcodetoopenid"
	refund_url                 = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	draw_ulr = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
)

const (
	CacheTokenName  = "cachetoken"
	CacheTicketName = "cacheticket"
	TokenTimeOut    = 7100
)

//微信错误定义
//https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1433747234
const (
	CODE_BUSY                 = -1
	CODE_SUCCESS              = 0
	CODE_INVALIDSECRET        = 40001
	CODE_INVALIDTOKEN         = 40002
	CODE_INVALIDOPENID        = 40003
	CODE_INVALIDMEDIATYPE     = 40004
	CODE_INVALIDFILETYPE      = 40005
	CODE_INVALIDFILESIZE      = 40006
	CODE_INVALIDMEDIAID       = 40007
	CODE_INVALIDMESSAGETYPE   = 40008
	CODE_INVALIDIMAGESIZE     = 40009
	CODE_INVALIDAUDIOSIZE     = 40010
	CODE_INVALIDVIDEOSIZE     = 40011
	CODE_INVALIDTHUMBNAILSIZE = 40012
	CODE_INVALIDAPPID         = 40013
	CODE_INVALIDACCESS_TOKEN  = 40014
	CODE_INVALIDMENUTYPE      = 40015
	CODE_INVALIDBUTTONNUMBERS = 40016
	//CODE_INVALIDBUTTONNUMBERS = 40017
	CODE_INVALIDBUTTONNAME_LENGTH        = 40018
	CODE_INVALIDBUTTONKEY_LENGTH         = 40019
	CODE_INVALIDBUTTONURL_LENGTH         = 40020
	CODE_INVALIDMENUVERSION              = 40021
	CODE_INVALIDSUBMENUSERIES            = 40022
	CODE_INVALIDSUBMENUBUTTONNUBMER      = 40023
	CODE_INVALIDSUBMENUBUTTONTYPE        = 40024
	CODE_INVALIDSUBMENUBUTTONNAME_LENGTH = 40025
	CODE_INVALIDSUBMENUBUTTONKEY_LENGTH  = 40026
	CODE_INVALIDSUBMENUBUTTONURL_LENGTH  = 40027
	CODE_INVALIDCUSTOMMENUUSER           = 40028
	CODE_INVALIDOAUTH_CODE               = 40029
	CODE_INVALIDREFRESH_TOKEN            = 40030
	CODE_INVALIDOPENIDLIST               = 40031
	CODE_INVALIDOPENIDLIST_LENGTH        = 40032
	CODE_INVALIDREQUESTCHARACTER         = 40033
	CODE_INVALIDREQUESTPARAM             = 40035
	CODE_INVALIDREQUESTFORMAT            = 40038
	CODE_INVALIDURL_LENGTH               = 40039
	CODE_INVALIDGROUPID                  = 40050
	CODE_INVALIDGROUPNAME                = 40051
	CODE_INVALIDARTICLEIDBYDELETEARTICLE = 40060
	//CODE_INVALIDGROUPNAME                = 40117
	CODE_INVALIDMEDIAIDSIZE = 40118
	CODE_ERROROFBUTTONTYPE  = 40119
	//CODE_ERROROFBUTTONTYPE               = 40120
	CODE_INVALIDMEDIAIDTYPE                  = 40121
	CODE_INVALIDWXID                         = 40132
	CODE_NOTSUPPORTIMAGEFORMAT               = 40137
	CODE_DONOTADDLINK                        = 40155
	CODE_LACKACCESS_TOKEN                    = 41001
	CODE_LACKAPPID                           = 41002
	CODE_LACKREFRESH_TOKEN                   = 41003
	CODE_LACKSECRET                          = 41004
	CODE_LACKMEDIAFILE                       = 41005
	CODE_LACKMEDIA_ID                        = 41006
	CODE_LACKSUBMENU                         = 41007
	CODE_LACKOAUTH_CODE                      = 41008
	CODE_LACKOPENID                          = 41009
	CODE_ACCESS_TOKENTIEOUT                  = 42001
	CODE_REFRESH_TOKENTIMEOUT                = 42002
	CODE_OAUTH_CODETIMEOUT                   = 42003
	CODE_REFRESH_TOKENANDACCESS_TOKENTIMEOUT = 42007
	CODE_NEEDGETREQUEST                      = 43001
	CODE_NEEDPOSTREQUEST                     = 43002
	CODE_NEEDHTTPSREQUEST                    = 43003
	CODE_NEEDRECEVIERCONCERN                 = 43004
	CODE_NEEDFRENDSSHIP                      = 43005
	CODE_NEEDRECOVERY                        = 43019
	CODE_MEDIAFILEISNULL                     = 44001
	CODE_POSTDATAISNULL                      = 44002
	CODE_MESSAGECONTENTISNULL                = 44003
	CODE_TEXTMESSAGECONTENTISNULL            = 44004
	CODE_MEDIAFILESIZEOVERLIMIT              = 45001
	CODE_MESSAGECONTENTOVERLIMIT             = 45002
	CODE_TITLEOVERLIMIT                      = 45003
	CODE_DESCRIPTIONOVERLIMIT                = 45004
	CODE_LINKOVERLIMIT                       = 45005
	CODE_IMAGE_LINKOVERLIMIT                 = 45006
	CODE_AUDIOTIMEOVERLIMIT                  = 45007
	CODE_MESSAGEOVERLIMIT                    = 45008
	CODE_INTERFACEINVOKEOVERLIMIT            = 45009
	CODE_SUBMENUOVERLIMIT                    = 45010
	CODE_INVOKEAPIFREQUENTLY                 = 45011
	CODE_RECEIVETIEOVERLIMIT                 = 45015
	CODE_NOTALLOWEDUPDATESYSTEMGROUP         = 45016
	CODE_GROUPNAMETOOLONG                    = 45017
	CODE_GROUPNUMBEROVERLIMIT                = 45018
	CODE_CUSTOMSERVICEINTERFACEOVERLIMIT     = 45047
	CODE_NOTEXISTMEDIADATA                   = 46001
	CODE_NOTEXISTMENUVERSION                 = 46002
	CODE_NOTEXISTMENUDATA                    = 46003
	CODE_NOTEXISTUSER                        = 46004
	CODE_ERROEOFRESOLVEJSONXML               = 47001
	CODE_APINOTAUTHORIZE                     = 48001
	CODE_USERREFUSEMESSAGE                   = 48002
	CODE_APIINTERFACEABOLISHED               = 48004
	CODE_APINOTALLOWEDDELETEONTENT           = 48005
	CODE_APINOTALLOWEDCLEARINVOKETIMES       = 48006
	CODE_NOTEXISTPERMISSIONOGSENDMESSAGE     = 48008
	CODE_USERNOTANTHORIZEAPI                 = 50001
	CODE_USERNOTALLOWED                      = 50002
	CODE_INVALIDPARAMETER                    = 61451
	CODE_INVALIDKF_ACCOUNT                   = 61452
	CODE_KF_ACCOUNTEXSITED                   = 61453
	CODE_INVALIDKF_ACCOUNTLENGTH             = 61454
	CODE_ILLEGALCHARACTERINKF_ACCOUNT        = 61455
	CODE_KF_ACCOUNTCOUNTEXCEEDED             = 61456
	CODE_INVALIDHEADIMAGEFILETYPE            = 61457
	CODE_SYSTEMERROR                         = 61450
	CODE_DATAFORMATERROE                     = 61500
	CODE_NOTEXISTMENUBYMENUID                = 65301
	CODE_NOTEXISTTHISUSER                    = 65302
	CODE_NOTEXISTDEFAULTMENU                 = 65303
	CODE_MATCHRULEISNULL                     = 65304
	CODE_MENUNUMBEROVERLIMIT                 = 65305
	CODE_NOTSUPPORTMENUACCOUNT               = 65306
	CODE_MENUINFOISNULL                      = 65307
	CODE_CONTAINNOTRESPONCEBUTTON            = 65308
	CODE_MENUISOFF                           = 65309
	CODE_COUNTRYINFONOTALLOWEDNULL           = 65310
	CODE_PROVICEINFONOTALLOWEDNULL           = 65311
	CODE_INVALIDCOUNTRYINFO                  = 65312
	CODE_INVALIDPROVICEINFO                  = 65313
	CODE_INVALIDCITYINFO                     = 65314
	CODE_TOOMANYURLLINK                      = 65316
	CODE_INVALIDURL                          = 65317
	CODE_INVALIDPOSTDATA                     = 9001001
	CODE_REMOTESERVERNOTAVAILABLE            = 9001002
	CODE_INVALIDTICKET                       = 9001003
	CODE_GETAROUDUSERERROR                   = 9001004
	CODE_GETMERCHANTSINFOERROR               = 9001005
	CODE_GETOPENIDERROR                      = 9001006
	CODE_LACKUPLOADFILE                      = 9001007
	CODE_INVALIDUPLOADFILETYPE               = 9001008
	CODE_INVALIDUPLOADFILESIZE               = 9001009
	CODE_UPLOADFILEERROR                     = 9001010
	CODE_INVALIDACCOUNT                      = 9001020
	CODE_NOTALLOWEDADDNEWDEVICE              = 9001021
	CODE_INVALIDADDDEVICENUMBER              = 9001022
	CODE_EXISTDEVICEID                       = 9001023
	CODE_INVALIDQUERYDEVICENUMBEROVERLIMIT   = 9001024
	CODE_INVALIDDEVICEID                     = 9001025
	CODE_INVALIDPAGEID                       = 9001026
	CODE_INVALIDPAGEPARAM                    = 9001027
	CODE_INVALIDDELETEPAGEIDONCEOVERLIMIT    = 9001028
	CODE_DELETEPAGEERROR                     = 9001029
	CODE_QUERYPAGEIDONCEOVERLIMIT            = 9001030
	CODE_INVALIDTIME                         = 9001031
	CODE_BINDPARAMERROR                      = 9001032
	CODE_INVALIDSTOREID                      = 9001033
	CODE_DEVICEREMARKTOOLONG                 = 9001034
	CODE_INVALIDDEVICEpARAM                  = 9001035
	CODE_INVALIDBEGIN                        = 9001036
)

var errdef = map[int]string{
	CODE_BUSY:                 "系统繁忙，此时请开发者稍候再试",
	CODE_SUCCESS:              "请求成功",
	CODE_INVALIDSECRET:        "获取 access_token 时 AppSecret 错误，或者 access_token 无效。请开发者认真比对 AppSecret 的正确性，或查看是否正在为恰当的公众号调用接口",
	CODE_INVALIDTOKEN:         "不合法的凭证类型",
	CODE_INVALIDOPENID:        "不合法的 OpenID ，请开发者确认 OpenID （该用户）是否已关注公众号，或是否是其他公众号的 OpenID",
	CODE_INVALIDMEDIATYPE:     "不合法的媒体文件类型",
	CODE_INVALIDFILETYPE:      "不合法的文件类型",
	CODE_INVALIDFILESIZE:      "不合法的文件大小",
	CODE_INVALIDMEDIAID:       "不合法的媒体文件 id",
	CODE_INVALIDMESSAGETYPE:   "不合法的消息类型",
	CODE_INVALIDIMAGESIZE:     "不合法的图片文件大小",
	CODE_INVALIDAUDIOSIZE:     "不合法的语音文件大小",
	CODE_INVALIDVIDEOSIZE:     "不合法的视频文件大小",
	CODE_INVALIDTHUMBNAILSIZE: "不合法的缩略图文件大小",
	CODE_INVALIDAPPID:         "不合法的 AppID ，请开发者检查 AppID 的正确性，避免异常字符，注意大小写",
	CODE_INVALIDACCESS_TOKEN:  "不合法的 access_token ，请开发者认真比对 access_token 的有效性（如是否过期），或查看是否正在为恰当的公众号调用接口",
	CODE_INVALIDMENUTYPE:      "不合法的菜单类型",
	CODE_INVALIDBUTTONNUMBERS: "不合法的按钮个数",
	//CODE_INVALIDBUTTONNUMBERS: 					"不合法的按钮个数",
	CODE_INVALIDBUTTONNAME_LENGTH:        "不合法的按钮名字长度",
	CODE_INVALIDBUTTONKEY_LENGTH:         "不合法的按钮 KEY 长度",
	CODE_INVALIDBUTTONURL_LENGTH:         "不合法的按钮 URL 长度",
	CODE_INVALIDMENUVERSION:              "不合法的菜单版本号",
	CODE_INVALIDSUBMENUSERIES:            "不合法的子菜单级数",
	CODE_INVALIDSUBMENUBUTTONNUBMER:      "不合法的子菜单按钮个数",
	CODE_INVALIDSUBMENUBUTTONTYPE:        "不合法的子菜单按钮类型",
	CODE_INVALIDSUBMENUBUTTONNAME_LENGTH: "不合法的子菜单按钮名字长度",
	CODE_INVALIDSUBMENUBUTTONKEY_LENGTH:  "不合法的子菜单按钮 KEY 长度",
	CODE_INVALIDSUBMENUBUTTONURL_LENGTH:  "不合法的子菜单按钮 URL 长度",
	CODE_INVALIDCUSTOMMENUUSER:           "不合法的自定义菜单使用用户",
	CODE_INVALIDOAUTH_CODE:               "不合法的 oauth_code",
	CODE_INVALIDREFRESH_TOKEN:            "不合法的 refresh_token",
	CODE_INVALIDOPENIDLIST:               "不合法的 openid 列表",
	CODE_INVALIDOPENIDLIST_LENGTH:        "不合法的 openid 列表长度",
	//CODE_INVALIDREQUESTCHARACTER:					"不合法的请求字符，不能包含 \uxxxx 格式的字符",
	CODE_INVALIDREQUESTPARAM:             "不合法的参数",
	CODE_INVALIDREQUESTFORMAT:            "不合法的请求格式",
	CODE_INVALIDURL_LENGTH:               "不合法的 URL 长度",
	CODE_INVALIDGROUPID:                  "不合法的分组 id",
	CODE_INVALIDGROUPNAME:                "分组名字不合法",
	CODE_INVALIDARTICLEIDBYDELETEARTICLE: "删除单篇图文时，指定的 article_idx 不合法",

	//CODE_INVALIDGROUPNAME                = 40117
	CODE_INVALIDMEDIAIDSIZE: "media_id 大小不合法",
	CODE_ERROROFBUTTONTYPE:  "button 类型错误",
	//CODE_ERROROFBUTTONTYPE               = 40120
	CODE_INVALIDMEDIAIDTYPE:                  "不合法的 media_id 类型",
	CODE_INVALIDWXID:                         "不合法的微信号",
	CODE_NOTSUPPORTIMAGEFORMAT:               "不支持的图片格式",
	CODE_DONOTADDLINK:                        "请勿添加其他公众号的主页链接",
	CODE_LACKACCESS_TOKEN:                    "缺少 access_token 参数",
	CODE_LACKAPPID:                           "缺少 appid 参数",
	CODE_LACKREFRESH_TOKEN:                   "缺少 refresh_token 参数",
	CODE_LACKSECRET:                          "缺少 secret 参数",
	CODE_LACKMEDIAFILE:                       "缺少多媒体文件数据",
	CODE_LACKMEDIA_ID:                        "缺少 media_id 参数",
	CODE_LACKSUBMENU:                         "缺少子菜单数据",
	CODE_LACKOAUTH_CODE:                      "缺少 oauth code",
	CODE_LACKOPENID:                          "缺少 openid",
	CODE_ACCESS_TOKENTIEOUT:                  "access_token 超时，请检查 access_token 的有效期，请参考基础支持 - 获取 access_token 中，对 access_token 的详细机制说明",
	CODE_REFRESH_TOKENTIMEOUT:                "refresh_token 超时",
	CODE_OAUTH_CODETIMEOUT:                   "oauth_code 超时",
	CODE_REFRESH_TOKENANDACCESS_TOKENTIMEOUT: "用户修改微信密码， accesstoken 和 refreshtoken 失效，需要重新授权",
	CODE_NEEDGETREQUEST:                      "需要 GET 请求",
	CODE_NEEDPOSTREQUEST:                     "需要 POST 请求",
	CODE_NEEDHTTPSREQUEST:                    "需要 HTTPS 请求",
	CODE_NEEDRECEVIERCONCERN:                 "需要接收者关注",
	CODE_NEEDFRENDSSHIP:                      "需要好友关系",
	CODE_NEEDRECOVERY:                        "需要将接收者从黑名单中移除",
	CODE_MEDIAFILEISNULL:                     "多媒体文件为空",
	CODE_POSTDATAISNULL:                      "POST 的数据包为空",
	CODE_MESSAGECONTENTISNULL:                "图文消息内容为空",
	CODE_TEXTMESSAGECONTENTISNULL:            "文本消息内容为空",
	CODE_MEDIAFILESIZEOVERLIMIT:              "多媒体文件大小超过限制",
	CODE_MESSAGECONTENTOVERLIMIT:             "消息内容超过限制",
	CODE_TITLEOVERLIMIT:                      "标题字段超过限制",
	CODE_DESCRIPTIONOVERLIMIT:                "描述字段超过限制",
	CODE_LINKOVERLIMIT:                       "链接字段超过限制",
	CODE_IMAGE_LINKOVERLIMIT:                 "图片链接字段超过限制",
	CODE_AUDIOTIMEOVERLIMIT:                  "语音播放时间超过限制",
	CODE_MESSAGEOVERLIMIT:                    "图文消息超过限制",
	CODE_INTERFACEINVOKEOVERLIMIT:            "接口调用超过限制",
	CODE_SUBMENUOVERLIMIT:                    "创建菜单个数超过限制",
	CODE_INVOKEAPIFREQUENTLY:                 "API 调用太频繁，请稍候再试",
	CODE_RECEIVETIEOVERLIMIT:                 "回复时间超过限制",
	CODE_NOTALLOWEDUPDATESYSTEMGROUP:         "系统分组，不允许修改",
	CODE_GROUPNAMETOOLONG:                    "分组名字过长",
	CODE_GROUPNUMBEROVERLIMIT:                "分组数量超过上限",
	CODE_CUSTOMSERVICEINTERFACEOVERLIMIT:     "客服接口下行条数超过上限",
	CODE_NOTEXISTMEDIADATA:                   "不存在媒体数据",
	CODE_NOTEXISTMENUVERSION:                 "不存在的菜单版本",
	CODE_NOTEXISTMENUDATA:                    "不存在的菜单数据",
	CODE_NOTEXISTUSER:                        "不存在的用户",
	CODE_ERROEOFRESOLVEJSONXML:               "解析 JSON/XML 内容错误",
	CODE_APINOTAUTHORIZE:                     "api 功能未授权，请确认公众号已获得该接口，可以在公众平台官网 - 开发者中心页中查看接口权限",
	CODE_USERREFUSEMESSAGE:                   "粉丝拒收消息（粉丝在公众号选项中，关闭了 “ 接收消息 ” ）",
	CODE_APIINTERFACEABOLISHED:               "api 接口被封禁，请登录 mp.weixin.qq.com 查看详情",
	CODE_APINOTALLOWEDDELETEONTENT:           "api 禁止删除被自动回复和自定义菜单引用的素材",
	CODE_APINOTALLOWEDCLEARINVOKETIMES:       "api 禁止清零调用次数，因为清零次数达到上限",
	CODE_NOTEXISTPERMISSIONOGSENDMESSAGE:     "没有该类型消息的发送权限",
	CODE_USERNOTANTHORIZEAPI:                 "用户未授权该 api",
	CODE_USERNOTALLOWED:                      "用户受限，可能是违规后接口被封禁",
	CODE_INVALIDPARAMETER:                    "参数错误 (invalid parameter)",
	CODE_INVALIDKF_ACCOUNT:                   "无效客服账号 (invalid kf_account)",
	CODE_KF_ACCOUNTEXSITED:                   "客服帐号已存在 (kf_account exsited)",
	CODE_INVALIDKF_ACCOUNTLENGTH:             "客服帐号名长度超过限制 ( 仅允许 10 个英文字符，不包括 @ 及 @ 后的公众号的微信号 )(invalid kf_acount length)",
	CODE_ILLEGALCHARACTERINKF_ACCOUNT:        "客服帐号名包含非法字符 ( 仅允许英文 + 数字 )(illegal character in kf_account)",
	CODE_KF_ACCOUNTCOUNTEXCEEDED:             "客服帐号个数超过限制 (10 个客服账号 )(kf_account count exceeded)",
	CODE_INVALIDHEADIMAGEFILETYPE:            "无效头像文件类型 (invalid file type)",
	CODE_SYSTEMERROR:                         "系统错误 (system error)",
	CODE_DATAFORMATERROE:                     "日期格式错误",
	CODE_NOTEXISTMENUBYMENUID:                "不存在此 menuid 对应的个性化菜单",
	CODE_NOTEXISTTHISUSER:                    "没有相应的用户",
	CODE_NOTEXISTDEFAULTMENU:                 "没有默认菜单，不能创建个性化菜单",
	CODE_MATCHRULEISNULL:                     "MatchRule 信息为空",
	CODE_MENUNUMBEROVERLIMIT:                 "个性化菜单数量受限",
	CODE_NOTSUPPORTMENUACCOUNT:               "不支持个性化菜单的帐号",
	CODE_MENUINFOISNULL:                      "个性化菜单信息为空",
	CODE_CONTAINNOTRESPONCEBUTTON:            "包含没有响应类型的 button",
	CODE_MENUISOFF:                           "个性化菜单开关处于关闭状态",
	CODE_COUNTRYINFONOTALLOWEDNULL:           "填写了省份或城市信息，国家信息不能为空",
	CODE_PROVICEINFONOTALLOWEDNULL:           "填写了城市信息，省份信息不能为空",
	CODE_INVALIDCOUNTRYINFO:                  "不合法的国家信息",
	CODE_INVALIDPROVICEINFO:                  "不合法的省份信息",
	CODE_INVALIDCITYINFO:                     "不合法的城市信息",
	CODE_TOOMANYURLLINK:                      "该公众号的菜单设置了过多的域名外跳（最多跳转到 3 个域名的链接）",
	CODE_INVALIDURL:                          "不合法的 URL",
	CODE_INVALIDPOSTDATA:                     "POST 数据参数不合法",
	CODE_REMOTESERVERNOTAVAILABLE:            "远端服务不可用",
	CODE_INVALIDTICKET:                       "Ticket 不合法",
	CODE_GETAROUDUSERERROR:                   "获取摇周边用户信息失败",
	CODE_GETMERCHANTSINFOERROR:               "获取商户信息失败",
	CODE_GETOPENIDERROR:                      "获取 OpenID 失败",
	CODE_LACKUPLOADFILE:                      "上传文件缺失",
	CODE_INVALIDUPLOADFILETYPE:               "上传素材的文件类型不合法",
	CODE_INVALIDUPLOADFILESIZE:               "上传素材的文件尺寸不合法",
	CODE_UPLOADFILEERROR:                     "上传失败",
	CODE_INVALIDACCOUNT:                      "帐号不合法",
	CODE_NOTALLOWEDADDNEWDEVICE:              "已有设备激活率低于 50% ，不能新增设备",
	CODE_INVALIDADDDEVICENUMBER:              "设备申请数不合法，必须为大于 0 的数字",
	CODE_EXISTDEVICEID:                       "已存在审核中的设备 ID 申请",
	CODE_INVALIDQUERYDEVICENUMBEROVERLIMIT:   "一次查询设备 ID 数量不能超过 50",
	CODE_INVALIDDEVICEID:                     "设备 ID 不合法",
	CODE_INVALIDPAGEID:                       "页面 ID 不合法",
	CODE_INVALIDPAGEPARAM:                    "页面参数不合法",
	CODE_INVALIDDELETEPAGEIDONCEOVERLIMIT:    "一次删除页面 ID 数量不能超过 10",
	CODE_DELETEPAGEERROR:                     "页面已应用在设备中，请先解除应用关系再删除",
	CODE_QUERYPAGEIDONCEOVERLIMIT:            "一次查询页面 ID 数量不能超过 50",
	CODE_INVALIDTIME:                         "时间区间不合法",
	CODE_BINDPARAMERROR:                      "保存设备与页面的绑定关系参数错误",
	CODE_INVALIDSTOREID:                      "门店 ID 不合法",
	CODE_DEVICEREMARKTOOLONG:                 "设备备注信息过长",
	CODE_INVALIDDEVICEpARAM:                  "设备申请参数不合法",
	CODE_INVALIDBEGIN:                        "查询起始值 begin 不合法",
}

func getErrMessage(id int) string {
	if v, ok := errdef[id]; ok {
		return v
	}
	return "未定义错误"
}
