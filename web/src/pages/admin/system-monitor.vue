<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import {
  Card,
  Row,
  Col,
  Statistic,
  Progress,
  Table,
  message,
  Spin,
  Divider,
  Descriptions
} from 'ant-design-vue'

// 系统信息
const systemInfo = ref({})
const latestMonitor = ref({})
const monitorStats = ref([])
const loading = ref(false)

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const res = await axios.get('/api/admin/system/info')
    if (res.data.code === 0) {
      systemInfo.value = res.data.data
    } else {
      message.error(res.data.msg || '获取系统信息失败')
    }
  } catch (error) {
    console.error('Failed to fetch system info:', error)
    message.error('获取系统信息失败')
  }
}

// 获取最新监控数据
const fetchLatestMonitorData = async () => {
  try {
    const res = await axios.get('/api/admin/system/monitor/latest')
    if (res.data.code === 0) {
      latestMonitor.value = res.data.data || {}
    } else {
      message.error(res.data.msg || '获取监控数据失败')
    }
  } catch (error) {
    console.error('Failed to fetch latest monitor data:', error)
    message.error('获取监控数据失败')
  }
}

// 获取监控统计数据
const fetchMonitorStats = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/admin/system/monitor/stats', {
      params: {
        hours: 24
      }
    })
    if (res.data.code === 0) {
      monitorStats.value = res.data.data || []
    } else {
      message.error(res.data.msg || '获取统计数据失败')
    }
  } catch (error) {
    console.error('Failed to fetch monitor stats:', error)
    message.error('获取统计数据失败')
  } finally {
    loading.value = false
  }
}

// 格式化字节大小
const formatBytes = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化时间
const formatUptime = (seconds) => {
  const days = Math.floor(seconds / (24 * 3600))
  const hours = Math.floor((seconds % (24 * 3600)) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = Math.floor(seconds % 60)
  
  return `${days}天 ${hours}小时 ${minutes}分钟 ${secs}秒`
}

// 定时器
let timer = null

onMounted(() => {
  fetchSystemInfo()
  fetchLatestMonitorData()
  fetchMonitorStats()
  
  // 每30秒刷新一次数据
  timer = setInterval(() => {
    fetchSystemInfo()
    fetchLatestMonitorData()
  }, 30000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<template>
  <page-container title="系统监控">
    <a-spin :spinning="loading">
      <a-row :gutter="[16, 16]">
        <!-- CPU使用率 -->
        <a-col :span="24" :md="8">
          <a-card>
            <a-statistic 
              title="CPU使用率" 
              :value="systemInfo.cpu_usage || 0" 
              suffix="%" 
            />
            <a-progress 
              :percent="systemInfo.cpu_usage || 0" 
              :status="systemInfo.cpu_usage > 80 ? 'exception' : 'normal'" 
              style="margin-top: 16px;"
            />
          </a-card>
        </a-col>
        
        <!-- 内存使用率 -->
        <a-col :span="24" :md="8">
          <a-card>
            <a-statistic 
              title="内存使用率" 
              :value="systemInfo.memory_usage || 0" 
              suffix="%" 
            />
            <a-progress 
              :percent="systemInfo.memory_usage || 0" 
              :status="systemInfo.memory_usage > 80 ? 'exception' : 'normal'" 
              style="margin-top: 16px;"
            />
            <div style="margin-top: 8px; font-size: 12px;">
              {{ formatBytes(systemInfo.memory_used) }} / {{ formatBytes(systemInfo.memory_total) }}
            </div>
          </a-card>
        </a-col>
        
        <!-- 磁盘使用率 -->
        <a-col :span="24" :md="8">
          <a-card>
            <a-statistic 
              title="磁盘使用率" 
              :value="systemInfo.disk_usage || 0" 
              suffix="%" 
            />
            <a-progress 
              :percent="systemInfo.disk_usage || 0" 
              :status="systemInfo.disk_usage > 90 ? 'exception' : 'normal'" 
              style="margin-top: 16px;"
            />
            <div style="margin-top: 8px; font-size: 12px;">
              {{ formatBytes(systemInfo.disk_used) }} / {{ formatBytes(systemInfo.disk_total) }}
            </div>
          </a-card>
        </a-col>
      </a-row>
      
      <a-divider />
      
      <a-row :gutter="[16, 16]">
        <!-- 系统信息 -->
        <a-col :span="24" :md="12">
          <a-card title="系统信息">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="主机名">
                {{ systemInfo.hostname || 'N/A' }}
              </a-descriptions-item>
              <a-descriptions-item label="操作系统">
                {{ systemInfo.os || 'N/A' }}
              </a-descriptions-item>
              <a-descriptions-item label="平台">
                {{ systemInfo.platform || 'N/A' }} {{ systemInfo.platform_version || '' }}
              </a-descriptions-item>
              <a-descriptions-item label="内核版本">
                {{ systemInfo.kernel_version || 'N/A' }}
              </a-descriptions-item>
              <a-descriptions-item label="CPU核心数">
                {{ systemInfo.cpu_count || 0 }}
              </a-descriptions-item>
              <a-descriptions-item label="协程数">
                {{ systemInfo.goroutine_count || 0 }}
              </a-descriptions-item>
              <a-descriptions-item label="系统运行时间">
                {{ systemInfo.uptime ? formatUptime(systemInfo.uptime) : 'N/A' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
        
        <!-- 最新监控数据 -->
        <a-col :span="24" :md="12">
          <a-card title="最新监控数据">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="网络流入">
                {{ latestMonitor.network_in ? latestMonitor.network_in + ' MB' : 'N/A' }}
              </a-descriptions-item>
              <a-descriptions-item label="网络流出">
                {{ latestMonitor.network_out ? latestMonitor.network_out + ' MB' : 'N/A' }}
              </a-descriptions-item>
              <a-descriptions-item label="请求次数">
                {{ latestMonitor.request_count || 0 }}
              </a-descriptions-item>
              <a-descriptions-item label="错误次数">
                {{ latestMonitor.error_count || 0 }}
              </a-descriptions-item>
              <a-descriptions-item label="平均响应时间">
                {{ latestMonitor.response_time ? latestMonitor.response_time + ' ms' : 'N/A' }}
              </a-descriptions-item>
              <a-descriptions-item label="在线用户数">
                {{ latestMonitor.online_users || 0 }}
              </a-descriptions-item>
              <a-descriptions-item label="监控时间">
                {{ latestMonitor.created_at ? new Date(latestMonitor.created_at).toLocaleString() : 'N/A' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-card>
        </a-col>
      </a-row>
    </a-spin>
  </page-container>
</template>

<style scoped>
</style>