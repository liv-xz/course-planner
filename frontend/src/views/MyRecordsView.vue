<script setup>
import { computed, onMounted, ref } from 'vue'
import http from '../api/http'

const loading = ref(false)
const records = ref([])
const statusText = { 0: '在读', 1: '通过', 2: '未通过' }
const statusType = { 0: 'info', 1: 'success', 2: 'danger' }

const summary = computed(() => {
  const passed = records.value.filter((r) => r.status === 1)
  return {
    total: records.value.length,
    passed: passed.length,
    credits: passed.reduce((s, r) => s + (r.credits || 0), 0),
  }
})

async function load() {
  loading.value = true
  try {
    records.value = await http.get('/records')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <el-row :gutter="16" class="stat-row">
      <el-col :xs="24" :sm="8">
        <div class="stat-card">
          <div class="stat-label">记录数</div>
          <div class="stat-value">{{ summary.total }}</div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-card">
          <div class="stat-label">通过课程</div>
          <div class="stat-value">{{ summary.passed }}</div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-card">
          <div class="stat-label">已通过学分</div>
          <div class="stat-value">{{ summary.credits }}</div>
        </div>
      </el-col>
    </el-row>

    <el-card shadow="never" class="panel-card">
      <template #header><div class="panel-title">我的修读记录</div></template>
      <el-table v-loading="loading" :data="records" stripe>
        <el-table-column label="课程" min-width="240">
          <template #default="{ row }">{{ row.course_code }} {{ row.course_name }}</template>
        </el-table-column>
        <el-table-column prop="score" label="成绩" width="90" />
        <el-table-column prop="semester" label="学期" width="140" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType[row.status]" effect="light">{{ statusText[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="credits" label="学分" width="90" />
      </el-table>
      <el-empty v-if="!loading && records.length === 0" description="还没有修读记录" />
    </el-card>
  </div>
</template>
