import request from '@/utils/request'

export function userApproveAs(data) {
  return request({
    url: '/userApproveAs',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function approveUserAs(data) {
  return request({
    url: '/approveUserAs',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function uplink(data) {
  return request({
    url: '/uplink',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function usermodify(data) {
  return request({
    url: '/usermodify',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}


export function userGetAllOffer(data) {
  return request({
    url: '/userGetAllOffer',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function adminGetAllOffers(data) {
  return request({
    url: '/adminGetAllOffers',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function getOfferHistory(data) {
  return request({
    url: '/getOfferHistory',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function getBalanceHistory(data) {
  return request({
    url: '/getBalanceHistory',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function getUserContracts(data) {
  return request({
    url: '/getUserContracts',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function getConfig(data) {
  return request({
    url: '/getConfig',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function adminModify(data) {
  return request({
    url: '/adminModify',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function getAdminActionHistory(data) {
  return request({
    url: '/getAdminActionHistory',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function getAdminMoneyHistory(data) {
  return request({
    url: '/getAdminMoneyHistory',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

export function adminGetAllOffer(data) {
  return request({
    url: '/adminGetAllOffer',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
export function getAllContract(data) {
  return request({
    url: '/getAllContract',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
