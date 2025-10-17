<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getResourceList } from '~/api/admin/resources.js'

const router = useRouter()
const message = useMessage()

// 资源列表
const resources = ref([])
const loading = ref(false)

// 获取所有资源
const fetchResources = async () => {
  loading.value = true
  try {
    const res = await getResourceList({ language: localStorage.getItem('admin-language') || 'zh-CN' })
    if (res.code === 0) resources.value = res.data?.resources || []
  } catch (error) {
    console.error('Failed to fetch resources:', error)
    message.error('获取资源列表失败')
  } finally {
    loading.value = false
  }
}

// 跳转到资源列表页
const goToResource = (slug) => {
  router.push(`/admin/${slug}`)
}

onMounted(() => {
  fetchResources()
})
</script>

<template>
  <page-container title="管理后台">
    <a-card>
      <a-spin :spinning="loading">
        <a-row :gutter="[16, 16]">
          <a-col
            v-for="resource in resources"
            :key="resource.slug"
            :span="6"
          >
            <a-card
              hoverable
              style="cursor: pointer"
              @click="goToResource(resource.slug)"
            >
              <a-card-meta :title="resource.title" />
            </a-card>
          </a-col>
        </a-row>
        
        <a-empty
          v-if="!loading && resources.length === 0"
          description="暂无资源"
        />
      </a-spin>
    </a-card>
  </page-container>
</template>

<style scoped>
:deep(.ant-card) {
  cursor: pointer;
}
</style>