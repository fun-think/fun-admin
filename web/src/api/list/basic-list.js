import request from '@/utils/request'

// 获取基础列表数据
export function getBasicList(params) {
  return request({
    url: '/api/v1/list/basic-list',
    method: 'post',
    data: params,
  })
}

// 创建基础列表项
export function createBasicItem(data) {
  return request({
    url: '/api/v1/list/basic-list',
    method: 'post',
    data,
  })
}

// 更新基础列表项
export function updateBasicItem(id, data) {
  return request({
    url: `/api/v1/list/basic-list/${id}`,
    method: 'put',
    data,
  })
}

// 删除基础列表项
export function deleteBasicItem(id) {
  return request({
    url: `/api/v1/list/basic-list/${id}`,
    method: 'delete',
  })
}
