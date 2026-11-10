<script setup>
import { onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import http from '../api/http'

const loading = ref(false)
const rows = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(20)

async function load() {
  loading.value = true
  try {
    const data = await http.get('/security/logs', { params: { page: page.value, size: size.value } })
    rows.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="toolbar">
      <div>
        <div class="panel-title">安全拦截日志</div>
        <div class="muted">记录提示词注入、越权查询等被拦截的输入。</div>
      </div>
      <el-button :icon="Refresh" @click="load">刷新</el-button>
    </div>
    <el-card shadow="never" class="panel-card">
      <el-table v-loading="loading" :data="rows" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="user_id" label="用户" width="90" />
        <el-table-column prop="input" label="输入内容" min-width="300">
          <template #default="{ row }">
            <div class="log-input">{{ row.input }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="拦截原因" min-width="220" />
        <el-table-column prop="created_at" label="时间" width="190" />
      </el-table>
      <el-empty v-if="!loading && rows.length === 0" description="暂无拦截记录" />
      <div class="pager">
        <el-pagination v-model:current-page="page" :page-size="size" :total="total" layout="total, prev, pager, next" @current-change="load" />
      </div>
    </el-card>
  </div>
</template>
