<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Grid } from '@element-plus/icons-vue'
import http from '../api/http'

const router = useRouter()
const loading = ref(false)
const rows = ref([])
const total = ref(0)
const majors = ref([])
const page = ref(1)
const size = ref(10)
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref()

const form = reactive({
  major_id: null,
  grade: '',
  name: '',
  total_credits: 165,
})

const rules = {
  major_id: [{ required: true, message: '请选择专业', trigger: 'change' }],
  grade: [
    { required: true, message: '请输入年级', trigger: 'blur' },
    { pattern: /^\d{4}$/, message: '年级为 4 位数字', trigger: 'blur' },
  ],
  name: [{ required: true, message: '请输入方案名称', trigger: 'blur' }],
}

const majorMap = computed(() => Object.fromEntries(majors.value.map((m) => [m.id, m.name])))

async function loadMajors() {
  const data = await http.get('/majors', { params: { page: 1, size: 100 } })
  majors.value = data.list || []
}

async function load() {
  loading.value = true
  try {
    const data = await http.get('/plans', { params: { page: page.value, size: size.value } })
    rows.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { major_id: majors.value[0]?.id || null, grade: '', name: '', total_credits: 165 })
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  Object.assign(form, {
    major_id: row.major_id,
    grade: row.grade,
    name: row.name,
    total_credits: row.total_credits ?? 0,
  })
  dialogVisible.value = true
}

async function save() {
  await formRef.value.validate()
  const payload = { ...form }
  if (!payload.total_credits) delete payload.total_credits
  if (editingId.value) {
    await http.put(`/plans/${editingId.value}`, {
      grade: payload.grade,
      name: payload.name,
      total_credits: payload.total_credits,
    })
    ElMessage.success('培养方案已更新')
  } else {
    await http.post('/plans', payload)
    ElMessage.success('培养方案已创建')
  }
  dialogVisible.value = false
  await load()
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确定删除「${row.name}」吗？`, '删除确认', { type: 'warning' })
    await http.delete(`/plans/${row.id}`)
    ElMessage.success('删除成功')
    await load()
  } catch (e) {
    if (e !== 'cancel') {
      // interceptor already displays business error
    }
  }
}

function goCategories(row) {
  router.push({ path: '/categories', query: { planId: row.id } })
}

onMounted(async () => {
  await loadMajors()
  await load()
})
</script>

<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-left">
        <span class="muted">按专业和年级维护培养方案，类别入口可继续配置课程树。</span>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新建方案</el-button>
    </div>

    <el-card shadow="never" class="panel-card">
      <el-table v-loading="loading" :data="rows" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="专业" min-width="180">
          <template #default="{ row }">{{ majorMap[row.major_id] || row.major_id }}</template>
        </el-table-column>
        <el-table-column prop="grade" label="年级" width="100" />
        <el-table-column prop="name" label="方案名称" min-width="260" />
        <el-table-column label="总学分" width="100">
          <template #default="{ row }">{{ row.total_credits ?? '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Grid" @click="goCategories(row)">类别</el-button>
            <el-button size="small" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" plain :icon="Delete" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination
          v-model:current-page="page"
          :page-size="size"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="load"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑培养方案' : '新建培养方案'" width="560px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="92px">
        <el-form-item label="专业" prop="major_id">
          <el-select v-model="form.major_id" placeholder="选择专业" style="width: 100%">
            <el-option v-for="m in majors" :key="m.id" :label="`${m.code} · ${m.name}`" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="年级" prop="grade">
          <el-input v-model="form.grade" maxlength="4" placeholder="如 2023" />
        </el-form-item>
        <el-form-item label="方案名称" prop="name">
          <el-input v-model="form.name" placeholder="如 2023 级计算机科学与技术培养方案" />
        </el-form-item>
        <el-form-item label="总学分">
          <el-input-number v-model="form.total_credits" :min="1" :max="300" :precision="1" style="width: 180px" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
