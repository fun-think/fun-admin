<script setup>
import { LogoutOutlined, ProfileOutlined, UserOutlined } from '@ant-design/icons-vue'

const message = useMessage()
const userStore = useUserStore()
const multiTabStore = useMultiTab()
const layoutMenuStore = useLayoutMenu()
const router = useRouter()
const { avatar, nickname, userInfo } = storeToRefs(userStore)
async function handleClick({ key }) {
  if (key === 'logout') {
    const hide = message.loading('退出登录...', 0)
    try {
      await userStore.logout()
    }
    finally {
      hide()
      message.success('退出登录成功', 3)
      router.push({
        path: '/login',
      }).then(() => {
        multiTabStore.clear()
        layoutMenuStore.clear()
      })
    }
  }
}
</script>

<template>
  <div class="flex items-center">
    <a-dropdown :trigger="['click']">
      <div class="flex items-center">
        <a-avatar :size="24" :src="userStore.avatar" class="mr-10px" />
        <span class="hidden md:block">{{ userStore.name }}</span>
      </div>
      <template #content>
        <a-doption>
          <a-space @click="toProfile">
            <icon-user />
            <span>{{ $t('menu.account.center') }}</span>
          </a-space>
        </a-doption>
        <a-doption>
          <a-space @click="logout">
            <icon-export />
            <span>{{ $t('menu.account.logout') }}</span>
          </a-space>
        </a-doption>
      </template>
    </a-dropdown>
  </div>
</template>
