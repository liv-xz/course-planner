<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import http from '../api/http'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(20)
const dialogVisible = ref(false)
const editingId = ref(null)
const formRef = ref()
const form = reactive({
  code: '',
  academic_year: '',
  term_no: 1,
  sequence_no: 1,
  display_name: '',
})
const rules = {
  code: [{ required: true, message: '请输入学期代码', trigger: 'blur' }],
  academic_year: [{ required: true, message: '请输入学年', trigger: 'blur' }],
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
}

function deriveAcademicYear(code) {
  const m = String(code).match(/^(\d{4})-(\d{4})-([12])$/)
  if (!m) return ''
  return `${m[1]}-${m[2]}`
}

function onCodeChange() {
  form.academic_year = deriveAcademicYear(form.code)
  const m = String(form.code).match(/-([12])$/)
  if (m) form.term_no = Number(m[1])
  if (!form.display_name && form.academic_year) {
    form.display_name = `${form.academic_year}学年${form.term_no === 1 ? '第一' : '第二'}学期`
  }
}

async function load() {
  loading.value = true
  try {
    const data = await http.get('/semesters', { params: { page: page.value, size: size.value } })
    rows.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { code: '', academic_year: '', term_no: 1, sequence_no: rows.value.length + 1, display_name: '' })
  dialogVisible.value = true
}

function openEdit(row) {
  editingId.value = row.id
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

async function save() {
  await formRef.value.validate()
  const payload = {
    code: form.code,
    academic_year: form.academic_year,
    term_no: Number(form.term_no),
    sequence_no: Number(form.sequence_no),
    display_name: form.display_name,
  }
  if (editingId.value) {
    await http.put(`/semesters/${editingId.value}`, payload)
    ElMessage.success('学期已更新')
  } else {
    await http.post('/semesters', payload)
    ElMessage.success('学期已创建')
  }
  dialogVisible.value = false
  await load()
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确定删除学期「${row.code}」吗？有开课/记录引用时不能删除。`, '删除确认', { type: 'warning' })
    await http.delete(`/semesters/${row.id}`)
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
      <div class="muted">学期使用日历格式，如 2023-2024-1；课程开课和修读记录都引用这里。</div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新建学期</el-button>
    </div>
    <el-card shadow="never" class="panel-card">
      <el-table v-loading="loading" :data="rows" stripe>
        <el-table-column prop="code" label="学期代码" width="150" />
        <el-table-column prop="academic_year" label="学年" width="140" />
        <el-table-column prop="term_no" label="学期" width="80" />
        <el-table-column prop="sequence_no" label="顺序" width="80" />
        <el-table-column prop="display_name" label="显示名称" min-width="220" />
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑学期' : '新建学期'" width="520px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="学期代码" prop="code">
          <el-input v-model="form.code" placeholder="2023-2024-1" @blur="onCodeChange" />
        </el-form-item>
        <el-form-item label="学年" prop="academic_year">
          <el-input v-model="form.academic_year" placeholder="2023-2024" />
        </el-form-item>
        <el-form-item label="第几学期">
          <el-select v-model="form.term_no" style="width: 160px">
            <el-option label="第一学期" :value="1" />
            <el-option label="第二学期" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sequence_no" :min="1" :max="999" style="width: 160px" />
        </el-form-item>
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="form.display_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
