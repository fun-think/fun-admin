import request from '@/utils/request'

// 获取CRUD表格数据
export function getCrudTable(params) {
  return request({
    url: '/api/v1/list/crud-table',
    method: 'post',
    data: params,
  })
}

// 创建CRUD表格项
export function createCrudItem(data) {
  return request({
    url: '/api/v1/list/crud-table',
    method: 'post',
    data,
  })
}

// 更新CRUD表格项
export function updateCrudItem(id, data) {
  return request({
    url: `/api/v1/list/crud-table/${id}`,
    method: 'put',
    data,
  })
}

// 删除CRUD表格项
export function deleteCrudItem(id) {
  return request({
    url: `/api/v1/list/crud-table/${id}`,
    method: 'delete',
  })
}
