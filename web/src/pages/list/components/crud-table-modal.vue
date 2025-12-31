<script setup>
import { cloneDeep } from 'lodash'
import { computed, nextTick, ref, watch, toRefs } from 'vue'

const props = defineProps({
  fields: {
    type: Array,
    default: () => [],
  },
})
const { fields } = toRefs(props)

const emit = defineEmits(['cancel', 'ok'])
const isUpdate = ref(false)
const visible = ref(false)
const title = computed(() => (isUpdate.value ? '编辑' : '新增'))
const formRef = ref()
const formData = ref({})
const labelCol = { style: { width: '100px' } }
const wrapperCol = { span: 24 }

function buildFormValues(record = null) {
  const values = {}
  fields.value.forEach(field => {
    if (field.type === 'boolean') {
      values[field.name] = record?.[field.name] ?? false
    }
    else {
      values[field.name] = record?.[field.name] ?? ''
    }
  })
  return values
}

function resetForm() {
  formData.value = buildFormValues()
  formRef.value?.resetFields()
}

watch(
  fields,
  () => {
    if (!visible.value) {
      formData.value = buildFormValues()
    }
  },
  { immediate: true },
)

async function open(record) {
  isUpdate.value = !!record?.id
  formData.value = buildFormValues(record)
  visible.value = true
  await nextTick()
  formRef.value?.resetFields()
}

async function handleOk() {
  try {
    await formRef.value?.validate()
    emit('ok', cloneDeep(formData.value))
    visible.value = false
  }
  catch (errorInfo) {
    console.log('Form Validate Failed:', errorInfo)
  }
}

function handleCancel() {
  resetForm()
  visible.value = false
  emit('cancel')
}

defineExpose({
  open,
})
</script>

<template>
  <a-modal v-model:open="visible" :title="title" @ok="handleOk" @cancel="handleCancel">
    <a-form ref="formRef" :model="formData" class="w-full" :label-col="labelCol" :wrapper-col="wrapperCol">
      <template v-for="field in fields" :key="field.name">
        <a-form-item
          :name="field.name"
          :label="field.label"
          :rules="field.required ? [{ required: true, message: `请输入${field.label}` }] : []"
        >
          <template v-if="field.type === 'textarea'">
            <a-textarea
              v-model:value="formData[field.name]"
              :maxlength="field.maxlength || 200"
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
          </template>
          <template v-else-if="field.type === 'select'">
            <a-select
              v-model:value="formData[field.name]"
              :placeholder="field.placeholder || `请选择${field.label}`"
            >
              <a-select-option v-for="option in field.options || []" :key="option.value" :value="option.value">
                {{ option.label }}
              </a-select-option>
            </a-select>
          </template>
          <template v-else-if="field.type === 'boolean'">
            <a-switch v-model:checked="formData[field.name]" />
          </template>
          <template v-else-if="field.type === 'number'">
            <a-input-number
              v-model:value="formData[field.name]"
              :placeholder="field.placeholder || `请输入${field.label}`"
              style="width: 100%"
            />
          </template>
          <template v-else-if="field.type === 'date' || field.type === 'datetime'">
            <a-date-picker
              v-model:value="formData[field.name]"
              :placeholder="field.placeholder || `请选择${field.label}`"
              :show-time="field.type === 'datetime'"
              style="width: 100%"
            />
          </template>
          <template v-else>
            <a-input
              v-model:value="formData[field.name]"
              :placeholder="field.placeholder || `请输入${field.label}`"
              :type="field.type === 'email' ? 'email' : 'text'"
              :maxlength="field.maxlength || 255"
            />
          </template>
        </a-form-item>
      </template>
    </a-form>
  </a-modal>
</template>
