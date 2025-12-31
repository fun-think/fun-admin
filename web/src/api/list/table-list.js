import request from '@/utils/request'

// 获取表格列表数据
export function getTableList(params) {
  return request({
    url: '/api/v1/list/table-list',
    method: 'post',
    data: params,
  })
}

// 创建表格列表项
export function createTableItem(data) {
  return request({
    url: '/api/v1/list/table-list',
    method: 'post',
    data,
  })
}

// 更新表格列表项
export function updateTableItem(id, data) {
  return request({
    url: `/api/v1/list/table-list/${id}`,
    method: 'put',
    data,
  })
}

// 删除表格列表项
export function deleteTableItem(id) {
  return request({
    url: `/api/v1/list/table-list/${id}`,
    method: 'delete',
  })
}
