import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function getUserInfo(data) {
  return request({
    url: '/getUserInfo',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function getAdminInfo(data) {
  return request({
    url: '/getAdminInfo',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function getInfo(data) {
  return request({
    url: '/getInfo',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function getName(data) {
  return request({
    url: '/getName',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function logout() {
  return request({
    url: '/logout',
    method: 'post'
  })
}
