<script setup>
import { ref, onMounted, h } from 'vue'
import { useRequest } from 'vue-request'
import axios from 'axios'
import {
  Card,
  Row,
  Col,
  Table,
  Button,
  Tag,
  message,
  Badge,
  Popconfirm,
  Empty,
  Pagination
} from 'ant-design-vue'

// 通知数据
const notifications = ref([])
const loading = ref(false)
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0
})

// 未读通知数量
const unreadCount = ref(0)

// 通知类型映射
const notificationTypes = {
  info: { color: 'blue', text: '信息' },
  success: { color: 'green', text: '成功' },
  warning: { color: 'orange', text: '警告' },
  error: { color: 'red', text: '错误' }
}

// 获取通知列表
const fetchNotifications = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/admin/notifications', {
      params: {
        page: pagination.value.current,
        page_size: pagination.value.pageSize
      }
    })
    
    if (res.data.code === 0) {
      notifications.value = res.data.data.items || []
      pagination.value.total = res.data.data.total
    } else {
      message.error(res.data.msg || '获取通知列表失败')
    }
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
    message.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

// 获取未读通知数量
const fetchUnreadCount = async () => {
  try {
    const res = await axios.get('/api/admin/notifications/unread-count')
    
    if (res.data.code === 0) {
      unreadCount.value = res.data.data.count
    } else {
      message.error(res.data.msg || '获取未读通知数量失败')
    }
  } catch (error) {
    console.error('Failed to fetch unread count:', error)
    message.error('获取未读通知数量失败')
  }
}

// 标记为已读
const markAsRead = async (id) => {
  try {
    const res = await axios.put(`/api/admin/notifications/${id}/read`)
    
    if (res.data.code === 0) {
      message.success('标记成功')
      fetchNotifications()
      fetchUnreadCount()
    } else {
      message.error(res.data.msg || '标记失败')
    }
  } catch (error) {
    console.error('Failed to mark as read:', error)
    message.error('标记失败')
  }
}

// 标记所有为已读
const markAllAsRead = async () => {
  try {
    const res = await axios.put('/api/admin/notifications/read-all')
    
    if (res.data.code === 0) {
      message.success('全部标记为已读')
      fetchNotifications()
      fetchUnreadCount()
    } else {
      message.error(res.data.msg || '标记失败')
    }
  } catch (error) {
    console.error('Failed to mark all as read:', error)
    message.error('标记失败')
  }
}

// 删除通知
const deleteNotification = async (id) => {
  try {
    const res = await axios.delete(`/api/admin/notifications/${id}`)
    
    if (res.data.code === 0) {
      message.success('删除成功')
      fetchNotifications()
      fetchUnreadCount()
    } else {
      message.error(res.data.msg || '删除失败')
    }
  } catch (error) {
    console.error('Failed to delete notification:', error)
    message.error('删除失败')
  }
}

// 处理分页变化
const handlePageChange = (page, pageSize) => {
  pagination.value.current = page
  pagination.value.pageSize = pageSize
  fetchNotifications()
}

// 通知表格列定义
const columns = [
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    customRender: ({ record }) => {
      return h(
        'div',
        {
          style: {
            fontWeight: record.is_read ? 'normal' : 'bold'
          }
        },
        record.title
      )
    }
  },
  {
    title: '内容',
    dataIndex: 'content',
    key: 'content',
    customRender: ({ record }) => {
      return h(
        'div',
        {
          style: {
            maxWidth: '300px',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap'
          }
        },
        record.content
      )
    }
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    customRender: ({ record }) => {
      const type = notificationTypes[record.type] || { color: 'default', text: record.type }
      return h(
        Tag,
        {
          color: type.color
        },
        type.text
      )
    }
  },
  {
    title: '状态',
    dataIndex: 'is_read',
    key: 'is_read',
    customRender: ({ record }) => {
      return h(
        Badge,
        {
          status: record.is_read ? 'default' : 'processing',
          text: record.is_read ? '已读' : '未读'
        }
      )
    }
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    customRender: ({ record }) => {
      return record.created_at ? new Date(record.created_at).toLocaleString() : ''
    }
  },
  {
    title: '操作',
    key: 'actions',
    customRender: ({ record }) => {
      return [
        h(
          Button,
          {
            type: 'link',
            size: 'small',
            onClick: () => markAsRead(record.id)
          },
          '标记为已读'
        ),
        h(
          Popconfirm,
          {
            title: '确定删除该通知？',
            onConfirm: () => deleteNotification(record.id)
          },
          {
            default: () => h(
              Button,
              {
                type: 'link',
                size: 'small',
                danger: true
              },
              '删除'
            )
          }
        )
      ]
    }
  }
]

onMounted(() => {
  fetchNotifications()
  fetchUnreadCount()
})
</script>

<template>
  <page-container title="通知中心">
    <a-card>
      <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
        <div>
          <a-badge :count="unreadCount" :overflow-count="99">
            <a-button type="primary">未读通知</a-button>
          </a-badge>
        </div>
        <div>
          <a-button 
            type="primary" 
            @click="markAllAsRead"
            :disabled="unreadCount === 0"
          >
            全部标记为已读
          </a-button>
        </div>
      </div>
      
      <a-table
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data-source="notifications"
        :pagination="false"
      >
        <template #emptyText>
          <a-empty description="暂无通知" />
        </template>
      </a-table>
      
      <div style="margin-top: 16px; display: flex; justify-content: flex-end;">
        <a-pagination
          v-model:current="pagination.current"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          show-size-changer
          @change="handlePageChange"
        />
      </div>
    </a-card>
  </page-container>
</template>

<style scoped>
</style>