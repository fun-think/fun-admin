import request from '@/utils/request'

// 获取咨询列表数据
export function getConsultList(params) {
  return request({
    url: '/api/v1/list/consult-list',
    method: 'post',
    data: params,
  })
}

// 创建咨询列表项
export function createConsultItem(data) {
  return request({
    url: '/api/v1/list/consult-list',
    method: 'post',
    data,
  })
}

// 更新咨询列表项
export function updateConsultItem(id, data) {
  return request({
    url: `/api/v1/list/consult-list/${id}`,
    method: 'put',
    data,
  })
}

// 删除咨询列表项
export function deleteConsultItem(id) {
  return request({
    url: `/api/v1/list/consult-list/${id}`,
    method: 'delete',
  })
}
