<script setup lang="jsx">
import { ref } from 'vue'
import { globalSearch } from '~/api/admin/resources.js'

const keyword = ref('')
const loading = ref(false)
const results = ref([])

const handleSearch = async () => {
  if (!keyword.value.trim()) return
  loading.value = true
  try {
    const res = await globalSearch(keyword.value, [], 5)
    if (res.code === 0) {
      results.value = res.data?.results || []
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <page-container title="全局搜索">
    <a-card>
      <a-space>
        <a-input v-model:value="keyword" placeholder="请输入关键词" style="width: 320px" @pressEnter="handleSearch" />
        <a-button type="primary" :loading="loading" @click="handleSearch">搜索</a-button>
      </a-space>
      <a-divider />
      <a-empty v-if="!results.length && !loading" description="暂无结果" />
      <div v-else>
        <div v-for="group in results" :key="group.slug" style="margin-bottom: 16px">
          <h3>{{ group.title }}</h3>
          <a-list :data-source="group.items" :renderItem="item => itemRender(item, group.slug)" />
        </div>
      </div>
    </a-card>
  </page-container>
</template>

<script lang="jsx">
export default {
  methods: {
    itemRender(item, slug) {
      return (
        <a-list-item>
          <div style="display:flex; justify-content:space-between; width:100%">
            <div>
              <div><strong>ID:</strong> { item.id }</div>
              {
                Object.keys(item).filter(k => k !== 'id').slice(0,3).map(k => (
                  <div style="color:#666; font-size:12px">{k}: {String(item[k])}</div>
                ))
              }
            </div>
            <a-button type="link" href={`#/admin/${slug}/edit/${item.id}`}>打开</a-button>
          </div>
        </a-list-item>
      )
    }
  }
}
</script>

