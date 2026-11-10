<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Document, TrendCharts, Upload } from '@element-plus/icons-vue'
import http from '../api/http'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')
const majors = ref([])
const courses = ref([])
const semesters = ref([])

const studentDialog = ref(false)
const editingStudentId = ref(null)
const studentFormRef = ref()
const studentForm = reactive({
  username: '',
  password: '123456',
  name: '',
  student_no: '',
  major_id: null,
  grade: '2023',
  status: 1,
})
const studentRules = {
  username: [{ required: true, message: '请输入用户名/学号', trigger: 'blur' }],
  student_no: [{ required: true, message: '请输入学号', trigger: 'blur' }],
  major_id: [{ required: true, message: '请选择专业', trigger: 'change' }],
  grade: [{ required: true, message: '请输入年级', trigger: 'blur' }, { pattern: /^\d{4}$/, message: '年级为 4 位数字', trigger: 'blur' }],
}

const recordsDialog = ref(false)
const recordStudent = ref(null)
const records = ref([])
const recordDialog = ref(false)
const editingRecordId = ref(null)
const recordFormRef = ref()
const recordForm = reactive({
  course_id: null,
  score: 80,
  semester: '',
  status: 1,
  credits: null,
})
const recordRules = {
  course_id: [{ required: true, message: '请选择课程', trigger: 'change' }],
  semester: [{ required: true, message: '请选择学期', trigger: 'change' }],
}

const auditDialog = ref(false)
const auditStudent = ref(null)
const audit = ref(null)
const importing = ref(false)
const importDialog = ref(false)
const importResult = ref(null)

const courseMap = computed(() => Object.fromEntries(courses.value.map((c) => [c.id, `${c.code} ${c.name}`])))
const majorMap = computed(() => Object.fromEntries(majors.value.map((m) => [m.id, m.name])))
const statusText = { 0: '在读', 1: '通过', 2: '未通过' }
const statusType = { 0: 'info', 1: 'success', 2: 'danger' }

async function load() {
  loading.value = true
  try {
    const data = await http.get('/students', {
      params: { page: page.value, size: size.value, keyword: keyword.value },
    })
    rows.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

async function loadBase() {
  const [m, c, s] = await Promise.all([
    http.get('/majors', { params: { page: 1, size: 100 } }),
    http.get('/courses', { params: { page: 1, size: 100 } }),
    http.get('/semesters', { params: { page: 1, size: 100 } }),
  ])
  majors.value = m.list || []
  courses.value = c.list || []
  semesters.value = s.list || []
}

function openCreate() {
  editingStudentId.value = null
  Object.assign(studentForm, {
    username: '',
    password: '123456',
    name: '',
    student_no: '',
    major_id: majors.value[0]?.id || null,
    grade: '2023',
    status: 1,
  })
  studentDialog.value = true
}

function openEdit(row) {
  editingStudentId.value = row.user_id
  Object.assign(studentForm, {
    username: row.username,
    password: '',
    name: row.name || '',
    student_no: row.student_no,
    major_id: row.major_id,
    grade: row.grade,
    status: row.status,
  })
  studentDialog.value = true
}

async function saveStudent() {
  await studentFormRef.value.validate()
  if (editingStudentId.value) {
    await http.put(`/students/${editingStudentId.value}`, {
      name: studentForm.name,
      student_no: studentForm.student_no,
      major_id: studentForm.major_id,
      grade: studentForm.grade,
      status: Number(studentForm.status),
    })
    ElMessage.success('学生信息已更新')
  } else {
    await http.post('/students', {
      username: studentForm.username,
      password: studentForm.password,
      name: studentForm.name,
      student_no: studentForm.student_no,
      major_id: studentForm.major_id,
      grade: studentForm.grade,
      status: Number(studentForm.status),
    })
    ElMessage.success('学生已创建')
  }
  studentDialog.value = false
  await load()
}

async function removeStudent(row) {
  try {
    await ElMessageBox.confirm(`确定删除学生「${row.name || row.username}」吗？`, '删除确认', { type: 'warning' })
    await http.delete(`/students/${row.user_id}`)
    ElMessage.success('删除成功')
    await load()
  } catch {
    // handled
  }
}

async function openRecords(row) {
  recordStudent.value = row
  records.value = await http.get(`/students/${row.user_id}/records`)
  recordsDialog.value = true
}

function openRecordCreate() {
  editingRecordId.value = null
  Object.assign(recordForm, {
    course_id: courses.value[0]?.id || null,
    score: 80,
    semester: (semesters.value.find((x) => x.code === '2024-2025-1') || semesters.value[0])?.code || '',
    status: 1,
    credits: null,
  })
  recordDialog.value = true
}

function openRecordEdit(row) {
  editingRecordId.value = row.id
  Object.assign(recordForm, {
    course_id: row.course_id,
    score: row.score,
    semester: row.semester || '',
    status: row.status,
    credits: row.credits,
  })
  recordDialog.value = true
}

async function saveRecord() {
  await recordFormRef.value.validate()
  const payload = {
    course_id: recordForm.course_id,
    score: recordForm.score === '' ? null : Number(recordForm.score),
    semester: recordForm.semester || null,
    status: Number(recordForm.status),
    credits: recordForm.credits === '' || recordForm.credits === null ? null : Number(recordForm.credits),
  }
  if (editingRecordId.value) {
    await http.put(`/records/${editingRecordId.value}`, payload)
    ElMessage.success('记录已更新')
  } else {
    await http.post(`/students/${recordStudent.value.user_id}/records`, payload)
    ElMessage.success('记录已添加')
  }
  recordDialog.value = false
  records.value = await http.get(`/students/${recordStudent.value.user_id}/records`)
}

async function removeRecord(row) {
  try {
    await ElMessageBox.confirm('确定删除这条修读记录吗？', '删除确认', { type: 'warning' })
    await http.delete(`/records/${row.id}`)
    ElMessage.success('删除成功')
    records.value = await http.get(`/students/${recordStudent.value.user_id}/records`)
  } catch {
    // handled
  }
}

async function handleImport(file) {
  if (!file?.raw || !recordStudent.value) return
  importing.value = true
  try {
    const fd = new FormData()
    fd.append('file', file.raw)
    importResult.value = await http.post(`/students/${recordStudent.value.user_id}/records/import`, fd)
    if (importResult.value.failed > 0) {
      importDialog.value = true
    } else {
      ElMessage.success(`导入成功，共 ${importResult.value.success} 条`)
    }
    records.value = await http.get(`/students/${recordStudent.value.user_id}/records`)
  } catch {
    // interceptor handles
  } finally {
    importing.value = false
  }
}

async function openAudit(row) {
  auditStudent.value = row
  audit.value = await http.get(`/students/${row.user_id}/audit`)
  auditDialog.value = true
}

function search() {
  page.value = 1
  load()
}

onMounted(async () => {
  await loadBase()
  await load()
})
</script>

<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-left">
        <el-input v-model="keyword" clearable placeholder="搜索姓名/学号/用户名" style="width: 260px" @keyup.enter="search" />
        <el-button @click="search">查询</el-button>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新建学生</el-button>
    </div>

    <el-card shadow="never" class="panel-card">
      <el-table v-loading="loading" :data="rows" stripe>
        <el-table-column prop="student_no" label="学号" width="130" />
        <el-table-column prop="username" label="用户名" width="130" />
        <el-table-column label="姓名" width="120">
          <template #default="{ row }">{{ row.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="专业" min-width="180">
          <template #default="{ row }">{{ majorMap[row.major_id] || row.major_name }}</template>
        </el-table-column>
        <el-table-column prop="grade" label="年级" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" effect="light">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Document" @click="openRecords(row)">修读记录</el-button>
            <el-button size="small" :icon="TrendCharts" @click="openAudit(row)">审核</el-button>
            <el-button size="small" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" plain :icon="Delete" @click="removeStudent(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination v-model:current-page="page" :page-size="size" :total="total" layout="total, prev, pager, next" @current-change="load" />
      </div>
    </el-card>

    <el-dialog v-model="studentDialog" :title="editingStudentId ? '编辑学生' : '新建学生'" width="560px">
      <el-form ref="studentFormRef" :model="studentForm" :rules="studentRules" label-width="92px">
        <el-form-item label="用户名/学号" prop="username">
          <el-input v-model="studentForm.username" :disabled="!!editingStudentId" />
        </el-form-item>
        <el-form-item v-if="!editingStudentId" label="初始密码">
          <el-input v-model="studentForm.password" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="studentForm.name" />
        </el-form-item>
        <el-form-item label="学号" prop="student_no">
          <el-input v-model="studentForm.student_no" />
        </el-form-item>
        <el-form-item label="专业" prop="major_id">
          <el-select v-model="studentForm.major_id" style="width: 100%">
            <el-option v-for="m in majors" :key="m.id" :label="`${m.code} · ${m.name}`" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="年级" prop="grade">
          <el-input v-model="studentForm.grade" maxlength="4" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="studentForm.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="studentDialog = false">取消</el-button>
        <el-button type="primary" @click="saveStudent">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="recordsDialog" :title="`修读记录 · ${recordStudent?.name || recordStudent?.username || ''}`" width="860px">
      <div class="toolbar compact">
        <span class="muted">同一门课多条通过记录时，审核按最高分计算。</span>
        <div class="toolbar-left">
          <el-upload
            :auto-upload="false"
            :show-file-list="false"
            accept=".csv,.xlsx"
            :on-change="handleImport"
          >
            <el-button :icon="Upload" :loading="importing">导入 CSV/Excel</el-button>
          </el-upload>
          <el-button type="primary" :icon="Plus" @click="openRecordCreate">添加记录</el-button>
        </div>
      </div>
      <el-table :data="records" size="small" stripe>
        <el-table-column label="课程" min-width="220">
          <template #default="{ row }">{{ row.course_code }} {{ row.course_name }}</template>
        </el-table-column>
        <el-table-column prop="score" label="成绩" width="80" />
        <el-table-column prop="semester" label="学期" width="130" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]" effect="light">{{ statusText[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="credits" label="学分" width="80" />
        <el-table-column label="操作" width="130">
          <template #default="{ row }">
            <el-button size="small" text @click="openRecordEdit(row)">编辑</el-button>
            <el-button size="small" text type="danger" @click="removeRecord(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="recordDialog" :title="editingRecordId ? '编辑修读记录' : '添加修读记录'" width="520px">
      <el-form ref="recordFormRef" :model="recordForm" :rules="recordRules" label-width="80px">
        <el-form-item label="课程" prop="course_id">
          <el-select v-model="recordForm.course_id" filterable style="width: 100%">
            <el-option v-for="c in courses" :key="c.id" :label="`${c.code} ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="成绩">
          <el-input-number v-model="recordForm.score" :min="0" :max="100" :precision="1" style="width: 160px" />
        </el-form-item>
        <el-form-item label="学期" prop="semester">
          <el-select v-model="recordForm.semester" style="width: 100%">
            <el-option v-for="s in semesters" :key="s.code" :label="`${s.code} · ${s.display_name}`" :value="s.code" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="recordForm.status" style="width: 100%">
            <el-option label="在读" :value="0" />
            <el-option label="通过" :value="1" />
            <el-option label="未通过" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="学分">
          <el-input-number v-model="recordForm.credits" :min="0" :max="20" :step="0.5" :precision="1" style="width: 160px" placeholder="默认课程学分" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="recordDialog = false">取消</el-button>
        <el-button type="primary" @click="saveRecord">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="importDialog" title="导入结果" width="620px">
      <div v-if="importResult" class="muted" style="margin-bottom: 12px">
        共 {{ importResult.total }} 行，成功 {{ importResult.success }} 行，失败 {{ importResult.failed }} 行。
      </div>
      <el-table v-if="importResult" :data="importResult.errors" size="small">
        <el-table-column prop="row" label="行号" width="80" />
        <el-table-column prop="reason" label="原因" />
      </el-table>
      <template #footer>
        <el-button type="primary" @click="importDialog = false">知道了</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="auditDialog" :title="`学分审核 · ${auditStudent?.name || auditStudent?.username || ''}`" width="820px">
      <div v-if="audit" class="audit-head">
        <div><span class="muted">培养方案</span><b>{{ audit.plan_name }}</b></div>
        <div><span class="muted">年级</span><b>{{ audit.grade }}</b></div>
      </div>
      <el-table v-if="audit" :data="audit.categories" row-key="id" default-expand-all :tree-props="{ children: 'children' }">
        <el-table-column prop="name" label="类别" min-width="260" />
        <el-table-column prop="required_credits" label="要求" width="100" />
        <el-table-column prop="completed_credits" label="已完成" width="100" />
        <el-table-column label="缺口" width="100">
          <template #default="{ row }">
            <el-tag :type="row.gap_credits > 0 ? 'danger' : 'success'" effect="light">{{ row.gap_credits }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
