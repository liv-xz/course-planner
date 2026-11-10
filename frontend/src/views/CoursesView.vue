<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Link, Connection, Calendar } from '@element-plus/icons-vue'
import http from '../api/http'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const keyword = ref('')
const plans = ref([])
const semesters = ref([])
const allCourses = ref([])

const courseDialog = ref(false)
const editingId = ref(null)
const formRef = ref()
const form = reactive({ code: '', name: '', credits: 3, nature: 1 })
const rules = {
  code: [{ required: true, message: '请输入课程代码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入课程名称', trigger: 'blur' }],
  credits: [{ required: true, message: '请输入学分', trigger: 'blur' }],
}

const mapDialog = ref(false)
const mapCourse = ref(null)
const mapPlanId = ref(null)
const mapCategoryId = ref(null)
const mapCategories = ref([])
const mapCurrent = ref([])

const prereqDialog = ref(false)
const prereqCourse = ref(null)
const prereqList = ref([])
const prereqSelect = ref(null)

const offeringDialog = ref(false)
const offeringCourse = ref(null)
const offeringList = ref([])
const offeringSemester = ref('')

const courseNameMap = computed(() => Object.fromEntries(allCourses.value.map((c) => [c.id, `${c.code} ${c.name}`])))
const natureText = (n) => (n === 1 ? '必修' : '选修')

async function load() {
  loading.value = true
  try {
    const data = await http.get('/courses', {
      params: { page: page.value, size: size.value, keyword: keyword.value },
    })
    rows.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

async function loadBase() {
  const [p, s, c] = await Promise.all([
    http.get('/plans', { params: { page: 1, size: 100 } }),
    http.get('/semesters', { params: { page: 1, size: 100 } }),
    http.get('/courses', { params: { page: 1, size: 100 } }),
  ])
  plans.value = p.list || []
  semesters.value = s.list || []
  allCourses.value = c.list || []
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { code: '', name: '', credits: 3, nature: 1 })
  courseDialog.value = true
}

function openEdit(row) {
  editingId.value = row.id
  Object.assign(form, { code: row.code, name: row.name, credits: row.credits, nature: row.nature })
  courseDialog.value = true
}

async function saveCourse() {
  await formRef.value.validate()
  if (editingId.value) {
    await http.put(`/courses/${editingId.value}`, {
      name: form.name,
      credits: Number(form.credits),
      nature: Number(form.nature),
    })
    ElMessage.success('课程已更新')
  } else {
    await http.post('/courses', {
      code: form.code,
      name: form.name,
      credits: Number(form.credits),
      nature: Number(form.nature),
    })
    ElMessage.success('课程已创建')
  }
  courseDialog.value = false
  await load()
  await loadBase()
}

async function removeCourse(row) {
  try {
    await ElMessageBox.confirm(`删除后课程将不再出现在推荐中，历史记录仍保留。确定删除「${row.name}」吗？`, '删除确认', { type: 'warning' })
    await http.delete(`/courses/${row.id}`)
    ElMessage.success('删除成功')
    await load()
    await loadBase()
  } catch {
    // handled
  }
}

function flattenLeaves(nodes, depth = 0, out = []) {
  for (const n of nodes || []) {
    if (!n.children || n.children.length === 0) {
      out.push({ ...n, label: `${'— '.repeat(depth)}${n.name}` })
    } else {
      flattenLeaves(n.children, depth + 1, out)
    }
  }
  return out
}

async function openMap(row) {
  mapCourse.value = row
  mapCurrent.value = await http.get(`/courses/${row.id}/categories`)
  mapPlanId.value = mapCurrent.value[0]?.plan_id || plans.value[0]?.id || null
  await loadMapCategories()
  mapDialog.value = true
}

async function loadMapCategories() {
  mapCategoryId.value = null
  mapCategories.value = []
  if (!mapPlanId.value) return
  const tree = await http.get(`/plans/${mapPlanId.value}/categories`)
  mapCategories.value = flattenLeaves(tree)
  const current = mapCurrent.value.find((m) => m.plan_id === mapPlanId.value)
  mapCategoryId.value = current?.category_id || null
}

async function saveMap() {
  if (!mapPlanId.value) {
    ElMessage.warning('请选择培养方案')
    return
  }
  await http.put(`/plans/${mapPlanId.value}/courses/${mapCourse.value.id}/category`, {
    category_id: mapCategoryId.value || null,
  })
  ElMessage.success('课程归属已保存')
  mapDialog.value = false
  mapCurrent.value = await http.get(`/courses/${mapCourse.value.id}/categories`)
}

async function openPrereq(row) {
  prereqCourse.value = row
  prereqSelect.value = null
  prereqList.value = await http.get(`/courses/${row.id}/prerequisites`)
  prereqDialog.value = true
}

async function addPrereq() {
  if (!prereqSelect.value) {
    ElMessage.warning('请选择先修课程')
    return
  }
  await http.post(`/courses/${prereqCourse.value.id}/prerequisites`, { prereq_course_id: prereqSelect.value })
  ElMessage.success('先修关系已添加')
  prereqList.value = await http.get(`/courses/${prereqCourse.value.id}/prerequisites`)
  prereqSelect.value = null
}

async function removePrereq(row) {
  await http.delete(`/prerequisites/${row.id}`)
  ElMessage.success('已删除')
  prereqList.value = await http.get(`/courses/${prereqCourse.value.id}/prerequisites`)
}

async function openOffering(row) {
  offeringCourse.value = row
  offeringSemester.value = ''
  offeringList.value = await http.get(`/courses/${row.id}/offerings`)
  offeringDialog.value = true
}

async function addOffering() {
  if (!offeringSemester.value) {
    ElMessage.warning('请选择学期')
    return
  }
  await http.post(`/courses/${offeringCourse.value.id}/offerings`, { semester: offeringSemester.value })
  ElMessage.success('开课记录已添加')
  offeringList.value = await http.get(`/courses/${offeringCourse.value.id}/offerings`)
  offeringSemester.value = ''
}

async function removeOffering(row) {
  await http.delete(`/offerings/${row.id}`)
  ElMessage.success('已删除')
  offeringList.value = await http.get(`/courses/${offeringCourse.value.id}/offerings`)
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
        <el-input v-model="keyword" clearable placeholder="搜索课程代码/名称" style="width: 260px" @keyup.enter="search" />
        <el-button @click="search">查询</el-button>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreate">新建课程</el-button>
    </div>

    <el-card shadow="never" class="panel-card">
      <el-table v-loading="loading" :data="rows" stripe>
        <el-table-column prop="code" label="代码" width="120" />
        <el-table-column prop="name" label="课程名称" min-width="220" />
        <el-table-column prop="credits" label="学分" width="90" />
        <el-table-column label="性质" width="90">
          <template #default="{ row }">
            <el-tag :type="row.nature === 1 ? 'danger' : 'success'" effect="light">{{ natureText(row.nature) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="390" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Link" @click="openMap(row)">归属</el-button>
            <el-button size="small" :icon="Connection" @click="openPrereq(row)">先修</el-button>
            <el-button size="small" :icon="Calendar" @click="openOffering(row)">开课</el-button>
            <el-button size="small" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" plain :icon="Delete" @click="removeCourse(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination v-model:current-page="page" :page-size="size" :total="total" layout="total, prev, pager, next" @current-change="load" />
      </div>
    </el-card>

    <el-dialog v-model="courseDialog" :title="editingId ? '编辑课程' : '新建课程'" width="520px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="课程代码" prop="code">
          <el-input v-model="form.code" :disabled="!!editingId" placeholder="如 CS101" />
        </el-form-item>
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="form.name" placeholder="如 程序设计基础" />
        </el-form-item>
        <el-form-item label="学分" prop="credits">
          <el-input-number v-model="form.credits" :min="0.5" :max="20" :step="0.5" :precision="1" style="width: 160px" />
        </el-form-item>
        <el-form-item label="性质">
          <el-radio-group v-model="form.nature">
            <el-radio :value="1">必修</el-radio>
            <el-radio :value="2">选修</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="courseDialog = false">取消</el-button>
        <el-button type="primary" @click="saveCourse">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="mapDialog" title="课程类别归属" width="600px">
      <div class="dialog-course">{{ mapCourse?.code }} {{ mapCourse?.name }}</div>
      <el-form label-width="90px">
        <el-form-item label="培养方案">
          <el-select v-model="mapPlanId" style="width: 100%" @change="loadMapCategories">
            <el-option v-for="p in plans" :key="p.id" :label="`${p.grade} · ${p.name}`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属类别">
          <el-select v-model="mapCategoryId" clearable placeholder="不选表示解除归属" style="width: 100%">
            <el-option v-for="c in mapCategories" :key="c.id" :label="c.label" :value="c.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <el-alert type="info" :closable="false" show-icon title="同一培养方案内，一门课程只保留一个叶子类别归属。" />
      <div class="current-map">
        <div class="muted">当前归属</div>
        <el-tag v-for="m in mapCurrent" :key="m.plan_id + '-' + m.category_id" class="map-tag">
          {{ m.plan_name }} / {{ m.category_name }}
        </el-tag>
        <span v-if="mapCurrent.length === 0" class="muted">暂无归属</span>
      </div>
      <template #footer>
        <el-button @click="mapDialog = false">取消</el-button>
        <el-button type="primary" @click="saveMap">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="prereqDialog" title="先修关系" width="620px">
      <div class="dialog-course">{{ prereqCourse?.code }} {{ prereqCourse?.name }}</div>
      <div class="inline-form">
        <el-select v-model="prereqSelect" filterable placeholder="选择先修课程" style="flex: 1">
          <el-option
            v-for="c in allCourses.filter((x) => x.id !== prereqCourse?.id)"
            :key="c.id"
            :label="`${c.code} ${c.name}`"
            :value="c.id"
          />
        </el-select>
        <el-button type="primary" @click="addPrereq">添加</el-button>
      </div>
      <el-table :data="prereqList" size="small" class="dialog-table">
        <el-table-column label="先修课程" min-width="240">
          <template #default="{ row }">{{ courseNameMap[row.prereq_course_id] || row.prereq_course_id }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" text @click="removePrereq(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog v-model="offeringDialog" title="开课学期" width="620px">
      <div class="dialog-course">{{ offeringCourse?.code }} {{ offeringCourse?.name }}</div>
      <div class="inline-form">
        <el-select v-model="offeringSemester" placeholder="选择学期" style="flex: 1">
          <el-option v-for="s in semesters" :key="s.code" :label="`${s.code} · ${s.display_name}`" :value="s.code" />
        </el-select>
        <el-button type="primary" @click="addOffering">添加</el-button>
      </div>
      <el-table :data="offeringList" size="small" class="dialog-table">
        <el-table-column prop="semester" label="学期" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" text @click="removeOffering(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
