import storage from 'store'
import { login, getInfo } from '@/api/login'
import { ACCESS_TOKEN } from '@/store/mutation-types'
import { welcome } from '@/utils/util'

const user = {
  state: {
    token: '',
    name: '',
    welcome: '',
    avatar: '',
    roles: [],
    info: {}
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_NAME: (state, { name, welcome }) => {
      state.name = name
      state.welcome = welcome
    },
    SET_AVATAR: (state, avatar) => {
      state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    },
    SET_INFO: (state, info) => {
      state.info = info
    }
  },

  actions: {
    // 登录
    Login ({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        login(userInfo).then(response => {
          if (response.code === 10000) {
            const result = response.body
            storage.set(ACCESS_TOKEN, result.token, 60 * 60)
            commit('SET_TOKEN', result.token)
            resolve(response)
          }
          reject(response)
        }).catch(error => {
          reject(error)
        })
      })
    },

    // 获取用户信息
    GetInfo ({ commit }) {
      return new Promise((resolve, reject) => {
        getInfo().then(response => {
          console.log('请求用户信息返回', response)
          if (response.code === 10000) {
            const result = response.body
            if (result.role && result.role.permissions.length > 0) {
                const role = result.role
                role.permissions = result.role.permissions
                role.permissions.map(per => {
                  if (per.actionEntitySet != null && per.actionEntitySet.length > 0) {
                    const action = per.actionEntitySet.map(action => { return action.action })
                    per.actionList = action
                  }
                })
                role.permissionList = role.permissions.map(permission => { return permission.permissionId })
                commit('SET_ROLES', result.role)
                commit('SET_INFO', result)
            }
            commit('SET_NAME', { name: result.name, welcome: welcome() })
            commit('SET_AVATAR', result.avatar)
            resolve(response)
          }
          reject(response)
        })
        .catch(error => {
          reject(error)
        })
      })
    },

    // 登出
    Logout ({ commit, state }) {
      return new Promise((resolve) => {
        // logout(state.token).then(() => {
          commit('SET_TOKEN', '')
          commit('SET_ROLES', [])
          storage.remove(ACCESS_TOKEN)
          resolve()
        // }).catch(() => {
        //   resolve()
        // }).finally(() => {
        // })
      })
    }

  }
}

export default user
