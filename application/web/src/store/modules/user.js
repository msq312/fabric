import { login, logout, getInfo, register } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
//import { generateRoutes } from '@/router/index'
import { resetRouter } from '@/router'

const getDefaultState = () => {
  return {
    token: getToken(),
    name: '',
    avatar: '',
    userType: ''
  }
}

const state = getDefaultState()

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_USERTYPE: (state, userType) => {
    state.userType = userType
    localStorage.setItem('userType', userType)
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    const formData = new FormData()
    formData.append('username', username.trim())
    formData.append('password', password)
    return new Promise((resolve, reject) => {
      login(formData).then(response => {
        commit('SET_TOKEN', response.jwt)
        setToken(response.jwt)
        commit('SET_USERTYPE', response.userType)
        // 获取用户信息并设置用户角色
        getInfo(response.jwt).then(infoResponse => {
          console.log('GetInfo response userType:', infoResponse.userType);
          const { userType } = infoResponse
          commit('SET_USERTYPE', userType)
          resetRouter() // 重新生成路由配置
          resolve()
        }).catch(error => {
          reject(error)
        })
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user register
  register({ commit }, userInfo) {
    const { username, password, userType } = userInfo
    const formData = new FormData()
    formData.append('username', username.trim())
    formData.append('password', password)
    formData.append('userType', userType)
    return new Promise((resolve, reject) => {
      register(formData).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo(state.token).then(response => {
        const { username } = response
        const { userType } = response

        if (!username) {
          return reject('Verification failed, please Login again.')
        }

        commit('SET_NAME', username)
        commit('SET_USERTYPE', userType)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        removeToken() // must remove  token  first
        localStorage.removeItem('userType') // 清除 localStorage 中的用户角色信息
        resetRouter()
        commit('RESET_STATE')
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      removeToken() // must remove  token  first
      commit('RESET_STATE')
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}

