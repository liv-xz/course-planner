<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import http from '../api/http'

const route = useRoute()
const loading = ref(false)
const plans = ref([])
const planId = ref(null)
const tree = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref()

const form = reactive({
  plan_id: null,
  parent_id: null,
  name: '',
  required_credits: 0,
  sort: 0,
})

const rules = {
  plan_id: [{ required: true, message: '请选择培养方案', trigger: 'change' }],
  name: [{ required: true, message: '请输入类别名称', trigger: 'blur' }],
}

const selectedPlan = computed(() => plans.value.find((p) => p.id === planId.value))
const treeSelectData = computed(() => tree.value)

async function loadPlans() {
  const data = await http.get('/plans', { params: { page: 1, size: 100 } })
  plans.value = data.list || []
  const queryPlan = Number(route.query.planId)
  planId.value = plans.value.some((p) => p.id === queryPlan) ? queryPlan : plans.value[0]?.id || null
}

async function loadTree() {
  if (!planId.value) {
    tree.value = []
    return
  }
  loading.value = true
  try {
    tree.value = await http.get(`/plans/${planId.value}/categories`)
  } finally {
    loading.value = false
  }
}

function openCreate(parent = null) {
  if (!planId.value) {
    ElMessage.warning('请先选择培养方案')
    return
  }
  editingId.value = null
  Object.assign(form, {
    plan_id: planId.value,
    parent_id: parent?.id ?? null,
    name: '',
    required_credits: 0,
    sort: 0,
  })
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  Object.assign(form, {
    plan_id: row.plan_id,
    parent_id: row.parent_id,
    name: row.name,
    required_credits: row.required_credits,
    sort: row.sort,
  })
  dialogVisible.value = true
}

async function save() {
  await formRef.value.validate()
  const payload = {
    parent_id: form.parent_id ?? null,
    name: form.name,
    required_credits: Number(form.required_credits) || 0,
    sort: Number(form.sort) || 0,
  }
  if (editingId.value) {
    await http.put(`/categories/${editingId.value}`, payload)
    ElMessage.success('类别已更新')
  } else {
    await http.post('/categories', { ...payload, plan_id: form.plan_id })
    ElMessage.success('类别已创建')
  }
  dialogVisible.value = false
  await loadTree()
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确定删除「${row.name}」吗？`, '删除确认', { type: 'warning' })
    await http.delete(`/categories/${row.id}`)
    ElMessage.success('删除成功')
    await loadTree()
  } catch {
    // handled by interceptor
  }
}

async function changePlan() {
  await loadTree()
}

onMounted(async () => {
  await loadPlans()
  await loadTree()
})
</script>

<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-left">
        <el-select v-model="planId" placeholder="选择培养方案" style="width: 360px" @change="changePlan">
          <el-option v-for="p in plans" :key="p.id" :label="`${p.grade} · ${p.name}`" :value="p.id" />
        </el-select>
        <span v-if="selectedPlan" class="muted">总学分 {{ selectedPlan.total_credits ?? '-' }}</span>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate()">新建类别</el-button>
    </div>

    <el-card shadow="never" class="panel-card">
      <el-table
        v-loading="loading"
        :data="tree"
        row-key="id"
        default-expand-all
        :tree-props="{ children: 'children' }"
      >
        <el-table-column prop="name" label="类别名称" min-width="260" />
        <el-table-column prop="required_credits" label="要求学分" width="120" />
        <el-table-column prop="sort" label="排序" width="90" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Plus" @click="openCreate(row)">加子类</el-button>
            <el-button size="small" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" plain :icon="Delete" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && tree.length === 0" description="该方案还没有类别" />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑类别' : '新建类别'" width="560px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="92px">
        <el-form-item label="所属方案">
          <el-input :model-value="selectedPlan?.name" disabled />
        </el-form-item>
        <el-form-item label="父类别">
          <el-tree-select
            v-model="form.parent_id"
            :data="treeSelectData"
            :props="{ label: 'name', children: 'children' }"
            node-key="id"
            check-strictly
            clearable
            placeholder="不选表示根类别"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如 通识教育课程" />
        </el-form-item>
        <el-form-item label="要求学分">
          <el-input-number v-model="form.required_credits" :min="0" :max="300" :precision="1" style="width: 180px" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="999" style="width: 180px" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
