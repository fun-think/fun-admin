<script setup>
import { message } from 'ant-design-vue'
import { getDashboardData } from '@/api/dashboard.js'

// 仪表板数据
const dashboardData = ref({})

// 获取当前语言
function getCurrentLanguage() {
  return localStorage.getItem('admin-language') || 'zh-CN'
}

// 加载状态
const loading = ref(false)

// 获取仪表板数据
async function fetchDashboardData() {
  if (loading.value)
    return
  loading.value = true
  try {
    const res = await getDashboardData({ language: getCurrentLanguage() })
    if (res.data.code === 0) {
      dashboardData.value = res.data.data
      processChartData()
    }
  }
  catch (error) {
    console.error('Failed to fetch dashboard data:', error)
    message.error('获取仪表板数据失败')
  }
  finally {
    loading.value = false
  }
}

// 图表配置
const chartOptions = ref({
  // 用户增长图表配置
  userGrowth: {
    tooltip: {
      trigger: 'axis',
    },
    xAxis: {
      type: 'category',
      data: [],
    },
    yAxis: {
      type: 'value',
    },
    series: [
      {
        data: [],
        type: 'line',
        smooth: true,
        itemStyle: {
          color: '#1890ff',
        },
      },
    ],
  },
  // 文章状态图表配置
  postStatus: {
    tooltip: {
      trigger: 'item',
    },
    legend: {
      top: '5%',
      left: 'center',
    },
    series: [
      {
        name: '文章状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2,
        },
        label: {
          show: false,
          position: 'center',
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold',
          },
        },
        labelLine: {
          show: false,
        },
        data: [],
      },
    ],
  },
})

// 处理图表数据
function processChartData() {
  if (!dashboardData.value)
    return

  // 处理用户增长数据
  const recentUsers = dashboardData.value.recent_users || []
  const dates = recentUsers.map(item => item.date)
  const counts = recentUsers.map(item => item.count)

  chartOptions.value.userGrowth.xAxis.data = dates
  chartOptions.value.userGrowth.series[0].data = counts

  // 处理文章状态数据
  const postStats = dashboardData.value.post_stats || {}
  const statusData = [
    {
      value: postStats.published || 0,
      name: postStats.published_label || '已发布',
      itemStyle: { color: '#52c41a' },
    },
    {
      value: postStats.draft || 0,
      name: postStats.draft_label || '草稿',
      itemStyle: { color: '#faad14' },
    },
    {
      value: postStats.archived || 0,
      name: postStats.archived_label || '已归档',
      itemStyle: { color: '#bfbfbf' },
    },
  ]

  chartOptions.value.postStatus.series[0].data = statusData
}

// 监听数据变化
watch(dashboardData, () => {
  processChartData()
})

onMounted(() => {
  fetchDashboardData()
})
</script>

<template>
  <page-container title="仪表板">
    <div class="dashboard">
      <!-- 统计卡片 -->
      <a-row :gutter="[16, 16]" class="mb-4">
        <a-col :span="6">
          <a-card>
            <a-statistic
              :title="dashboardData.user_count_label || '用户总数'"
              :value="dashboardData.user_count || 0"
            >
              <template #prefix>
                <UserOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              :title="dashboardData.department_count_label || '部门数量'"
              :value="dashboardData.department_count || 0"
            >
              <template #prefix>
                <TeamOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              :title="dashboardData.post_count_label || '文章总数'"
              :value="dashboardData.post_count || 0"
            >
              <template #prefix>
                <FileTextOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :span="6">
          <a-card>
            <a-statistic
              :title="dashboardData.role_count_label || '角色数量'"
              :value="dashboardData.role_count || 0"
            >
              <template #prefix>
                <IdcardOutlined />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>

      <!-- 图表区域 -->
      <a-row :gutter="[16, 16]">
        <a-col :span="12">
          <a-card title="用户增长趋势">
            <div style="height: 300px;">
              <v-chart
                v-if="chartOptions.userGrowth.xAxis.data.length > 0"
                :option="chartOptions.userGrowth"
                autoresize
              />
              <a-skeleton v-else />
            </div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="文章状态分布">
            <div style="height: 300px;">
              <v-chart
                v-if="chartOptions.postStatus.series[0].data.length > 0"
                :option="chartOptions.postStatus"
                autoresize
              />
              <a-skeleton v-else />
            </div>
          </a-card>
        </a-col>
      </a-row>

      <!-- 系统信息 -->
      <a-card title="系统信息" class="mb-4">
        <a-descriptions :column="3">
          <a-descriptions-item label="数据表数量">
            <DatabaseOutlined /> {{ dashboardData.system_info?.table_count || 0 }}
          </a-descriptions-item>
          <a-descriptions-item label="协程数量">
            <SyncOutlined /> {{ dashboardData.system_info?.goroutine_count || 0 }}
          </a-descriptions-item>
          <a-descriptions-item label="最后更新">
            {{ dashboardData.system_info?.last_update || '未知' }}
          </a-descriptions-item>
        </a-descriptions>
      </a-card>

      <!-- 最新动态 -->
      <a-card title="最新动态">
        <a-list
          item-layout="horizontal"
          :data-source="[]"
        >
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta
                description="最近的系统操作记录"
              >
                <template #title>
                  <a href="https://www.antdv.com/">{{ item.title }}</a>
                </template>
              </a-list-item-meta>
            </a-list-item>
          </template>
        </a-list>
      </a-card>

      <!-- 语言切换 -->
      <div style="margin-top: 16px; text-align: right;">
        <a-dropdown>
          <template #overlay>
            <a-menu @click="({ key }) => { localStorage.setItem('admin-language', key); location.reload(); }">
              <a-menu-item key="zh-CN">
                中文
              </a-menu-item>
              <a-menu-item key="en">
                English
              </a-menu-item>
            </a-menu>
          </template>
          <a-button>
            语言 <DownOutlined />
          </a-button>
        </a-dropdown>
      </div>
    </div>
  </page-container>
</template>

<style scoped>
.dashboard {
  padding: 16px;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>
