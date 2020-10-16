import antd from 'ant-design-vue/es/locale-provider/zh_CN'
import momentCN from 'moment/locale/zh-cn'

const components = {
  antLocale: antd,
  momentName: 'zh-cn',
  momentLocale: momentCN
}

const locale = {
  'message': '-',
  'menu.home': '主页',
  'menu.index': '首页',

  'menu.admin': '个人中心',
  'menu.admin.center': '个人中心',
  'menu.admin.logout': '退出',
  'menu.admin.set': '个人设置',

  'layouts.usermenu.dialog.title': '提示',
  'layouts.usermenu.dialog.content': '是否要退出登录'
}

export default {
  ...components,
  ...locale
}
