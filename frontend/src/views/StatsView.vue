<script setup>
import { computed, onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import http from '../api/http'

const loading = ref(false)
const overview = ref({})
const rows = ref([])

const maxGap = computed(() => Math.max(1, ...rows.value.map((r) => r.total_gap || 0)))

async function load() {
  loading.value = true
  try {
    const [o, g] = await Promise.all([
      http.get('/stats/overview'),
      http.get('/stats/credit-gap'),
    ])
    overview.value = o
    rows.value = g || []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading">
    <div class="toolbar">
      <div>
        <div class="panel-title">学分缺口统计</div>
        <div class="muted">按类别汇总全部学生的缺口人数与平均缺口，用于教学管理参考。</div>
      </div>
      <el-button :icon="Refresh" @click="load">刷新</el-button>
    </div>

    <el-row :gutter="16" class="stat-row">
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card"><div class="stat-label">学生</div><div class="stat-value">{{ overview.students ?? 0 }}</div></div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card"><div class="stat-label">培养方案</div><div class="stat-value">{{ overview.plans ?? 0 }}</div></div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card"><div class="stat-label">课程</div><div class="stat-value">{{ overview.courses ?? 0 }}</div></div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card"><div class="stat-label">类别</div><div class="stat-value">{{ overview.categories ?? 0 }}</div></div>
      </el-col>
      <el-col :xs="12" :sm="8" :lg="4">
        <div class="stat-card"><div class="stat-label">学期</div><div class="stat-value">{{ overview.semesters ?? 0 }}</div></div>
      </el-col>
    </el-row>

    <el-card shadow="never" class="panel-card">
      <template #header><div class="panel-title">类别缺口分布</div></template>
      <el-table :data="rows" stripe>
        <el-table-column prop="plan_name" label="培养方案" min-width="220" />
        <el-table-column prop="category_name" label="类别" min-width="180" />
        <el-table-column prop="required_credits" label="要求" width="90" />
        <el-table-column prop="student_count" label="学生数" width="90" />
        <el-table-column label="有缺口人数" width="110">
          <template #default="{ row }">
            <el-tag :type="row.gap_student_count > 0 ? 'danger' : 'success'" effect="light">{{ row.gap_student_count }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_gap" label="总缺口" width="100" />
        <el-table-column prop="average_gap" label="平均缺口" width="100" />
        <el-table-column label="缺口规模" min-width="180">
          <template #default="{ row }">
            <div class="gap-bar">
              <div class="gap-bar-inner" :style="{ width: Math.round(((row.total_gap || 0) / maxGap) * 100) + '%' }" />
            </div>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && rows.length === 0" description="暂无统计数据" />
    </el-card>
  </div>
</template>
