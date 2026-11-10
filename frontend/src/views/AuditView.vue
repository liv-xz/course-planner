<script setup>
import { computed, onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import http from '../api/http'

const loading = ref(false)
const audit = ref(null)

const totals = computed(() => {
  const roots = audit.value?.categories || []
  const required = roots.reduce((s, c) => s + (c.required_credits || 0), 0)
  const completed = roots.reduce((s, c) => s + (c.completed_credits || 0), 0)
  return { required, completed, gap: Math.max(required - completed, 0) }
})

const percent = computed(() => {
  if (!totals.value.required) return 0
  return Math.min(100, Math.round((totals.value.completed / totals.value.required) * 100))
})

async function load() {
  loading.value = true
  try {
    audit.value = await http.get('/audit/me')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading">
    <div class="audit-banner">
      <div>
        <div class="hero-kicker">CREDIT AUDIT</div>
        <h2>{{ audit?.plan_name || '学分审核' }}</h2>
        <p>按培养方案实时计算已修学分与缺口，同一门课只计算一次。</p>
      </div>
      <el-button :icon="Refresh" @click="load">刷新</el-button>
    </div>

    <el-row :gutter="16" class="stat-row">
      <el-col :xs="24" :sm="8">
        <div class="stat-card"><div class="stat-label">要求总学分</div><div class="stat-value">{{ totals.required }}</div></div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-card"><div class="stat-label">已完成</div><div class="stat-value">{{ totals.completed }}</div></div>
      </el-col>
      <el-col :xs="24" :sm="8">
        <div class="stat-card"><div class="stat-label">总缺口</div><div class="stat-value danger">{{ totals.gap }}</div></div>
      </el-col>
    </el-row>

    <el-card shadow="never" class="panel-card">
      <div class="progress-row">
        <span class="muted">总体完成度</span>
        <el-progress :percentage="percent" :stroke-width="12" style="flex: 1" />
      </div>
      <el-table v-if="audit" :data="audit.categories" row-key="id" default-expand-all :tree-props="{ children: 'children' }">
        <el-table-column prop="name" label="课程类别" min-width="280" />
        <el-table-column prop="required_credits" label="要求学分" width="120" />
        <el-table-column prop="completed_credits" label="已完成" width="120" />
        <el-table-column label="缺口" width="120">
          <template #default="{ row }">
            <el-tag :type="row.gap_credits > 0 ? 'danger' : 'success'" effect="light">
              {{ row.gap_credits > 0 ? `缺 ${row.gap_credits}` : '已达标' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" min-width="200">
          <template #default="{ row }">
            <el-progress
              :percentage="row.required_credits ? Math.min(100, Math.round((row.completed_credits / row.required_credits) * 100)) : 0"
              :status="row.gap_credits > 0 ? 'warning' : 'success'"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
