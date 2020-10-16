// eslint-disable-next-line
import { UserLayout, BasicLayout, BlankLayout} from '@/layouts'
// import { bxAnaalyse } from '@/core/icons'

const RouteView = {
  name: 'RouteView',
  render: (h) => h('router-view')
}

export const asyncRouterMap = [

  {
    path: '/',
    name: 'menu.home',
    component: BasicLayout,
    meta: { title: 'menu.home' },
    redirect: '/index',
    children: [
      {
        path: '/index',
        name: 'menu.index',
        component: () => import('@/views/index'),
        meta: { title: 'menu.index', keepAlive: true, permission: [ 'index', 'admin' ] }
      },
      // account
      {
        path: '/account',
        component: RouteView,
        redirect: '/account/center',
        name: 'menu.admin',
        meta: { title: 'menu.admin', icon: 'user', keepAlive: true, permission: [ 'admin' ] },
        children: [
          {
            path: '/account/center',
            name: 'menu.admin.center',
            component: () => import('@/views/account/center'),
            meta: { title: 'menu.admin', keepAlive: true, permission: [ 'admin' ] }
          },
          {
            path: '/account/settings',
            name: 'menu.admin.set',
            component: () => import('@/views/account/settings/Index'),
            meta: { title: 'menu.admin.set', hideHeader: true, permission: [ 'admin' ] },
            redirect: '/account/settings/base',
            hideChildrenInMenu: true,
            children: [
              {
                path: '/account/settings/base',
                name: 'BaseSettings',
                component: () => import('@/views/account/settings/BaseSetting'),
                meta: { title: '基本设置', hidden: true, permission: [ 'admin' ] }
              },
              {
                path: '/account/settings/security',
                name: 'SecuritySettings',
                component: () => import('@/views/account/settings/Security'),
                meta: { title: '安全设置', hidden: true, keepAlive: true, permission: [ 'admin' ] }
              },
              {
                path: '/account/settings/custom',
                name: 'CustomSettings',
                component: () => import('@/views/account/settings/Custom'),
                meta: { title: '个性化设置', hidden: true, keepAlive: true, permission: [ 'admin' ] }
              },
              {
                path: '/account/settings/binding',
                name: 'BindingSettings',
                component: () => import('@/views/account/settings/Binding'),
                meta: { title: '账户绑定', hidden: true, keepAlive: true, permission: [ 'admin' ] }
              },
              {
                path: '/account/settings/notification',
                name: 'NotificationSettings',
                component: () => import('@/views/account/settings/Notification'),
                meta: { title: '新消息通知', hidden: true, keepAlive: true, permission: [ 'admin' ] }
              }
            ]
          }
        ]
      },
      // other
      {
        path: '/power',
        name: 'power',
        component: RouteView,
        meta: { title: '权限管理', icon: 'slack', permission: [ 'admin' ] },
        redirect: '/other/icon-selector',
        children: [
          // {
          //   path: '/other/icon-selector',
          //   name: 'TestIconSelect',
          //   component: () => import('@/views/other/IconSelectorView'),
          //   meta: { title: 'IconSelector', icon: 'tool', keepAlive: true, permission: [ 'admin' ] }
          // },
          // {
          //   path: '/other/list',
          //   component: RouteView,
          //   meta: { title: '业务布局', icon: 'layout', permission: [ 'admin' ] },
          //   redirect: '/other/list/tree-list',
          //   children: [
              // {
              //   path: '/other/list/tree-list',
              //   name: 'TreeList',
              //   component: () => import('@/views/other/TreeList'),
              //   meta: { title: '树目录表格', keepAlive: true }
              // },
              // {
              //   path: '/other/list/edit-table',
              //   name: 'EditList',
              //   component: () => import('@/views/other/TableInnerEditList'),
              //   meta: { title: '内联编辑表格', keepAlive: true }
              // },
              {
                path: '/other/list/user-list',
                name: 'UserList',
                component: () => import('@/views/other/UserList'),
                meta: { title: '用户列表', keepAlive: true }
              },
              {
                path: '/other/list/role-list',
                name: 'RoleList',
                component: () => import('@/views/other/RoleList'),
                meta: { title: '角色列表', keepAlive: true }
              },
              // {
              //   path: '/other/list/system-role',
              //   name: 'SystemRole',
              //   component: () => import('@/views/role/RoleList'),
              //   meta: { title: '角色列表2', keepAlive: true }
              // },
              {
                path: '/other/list/permission-list',
                name: 'PermissionList',
                component: () => import('@/views/other/PermissionList'),
                meta: { title: '权限列表', keepAlive: true }
              }
            // ]
         // }
        ]
      }
    ]
  },
  {
    path: '*', redirect: '/404', hidden: true
  }
]

/**
 * 基础路由
 * @type { *[] }
 */
export const constantRouterMap = [
  {
    path: '/admin',
    component: UserLayout,
    redirect: '/admin/login',
    hidden: true,
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import('@/views/admin/Login')
      },
      {
        path: 'recover',
        name: 'recover',
        component: undefined
      }
    ]
  },

  {
    path: '/404',
    component: () => import('@/views/exception/404')
  }

]
