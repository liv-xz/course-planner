<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import http from '../api/http'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref()
const form = reactive({ code: '', name: '' })
const rules = {
  code: [{ required: true, message: '请输入专业代码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入专业名称', trigger: 'blur' }],
}

async function load() {
  loading.value = true
  try {
    const data = await http.get('/majors', { params: { page: page.value, size: size.value } })
    rows.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { code: '', name: '' })
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  Object.assign(form, { code: row.code, name: row.name })
  dialogVisible.value = true
}

async function save() {
  await formRef.value.validate()
  if (editingId.value) {
    await http.put(`/majors/${editingId.value}`, { ...form })
    ElMessage.success('专业已更新')
  } else {
    await http.post('/majors', { ...form })
    ElMessage.success('专业已创建')
  }
  dialogVisible.value = false
  await load()
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确定删除「${row.name}」吗？`, '删除确认', { type: 'warning' })
    await http.delete(`/majors/${row.id}`)
    ElMessage.success('删除成功')
    await load()
  } catch {
    // handled
  }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="toolbar">
      <div class="muted">专业是培养方案的归属单位，删除前需先清空方案和学生。</div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新建专业</el-button>
    </div>
    <el-card shadow="never" class="panel-card">
      <el-table v-loading="loading" :data="rows" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="代码" width="160" />
        <el-table-column prop="name" label="专业名称" min-width="240" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" plain :icon="Delete" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination v-model:current-page="page" :page-size="size" :total="total" layout="total, prev, pager, next" @current-change="load" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑专业' : '新建专业'" width="480px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="代码" prop="code"><el-input v-model="form.code" placeholder="如 CS" /></el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="form.name" placeholder="如 计算机科学与技术" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
