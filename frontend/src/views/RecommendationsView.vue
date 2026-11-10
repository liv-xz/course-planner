<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Refresh, Opportunity, ChatDotRound } from '@element-plus/icons-vue'
import http from '../api/http'

const router = useRouter()
const loading = ref(false)
const semesters = ref([])
const semester = ref('')
const rows = ref([])

async function loadSemesters() {
  const data = await http.get('/semesters', { params: { page: 1, size: 100 } })
  semesters.value = data.list || []
}

async function load() {
  loading.value = true
  try {
    rows.value = await http.get('/recommendations', { params: semester.value ? { semester: semester.value } : {} })
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadSemesters()
  await load()
})
</script>

<template>
  <div>
    <div class="recommend-banner">
      <div>
        <div class="hero-kicker">COURSE RECOMMENDATION</div>
        <h2>选课推荐</h2>
        <p>基于学分缺口、已修课程、先修关系和开课学期筛选。只推荐当前课程库中真实存在且已排课的课程。</p>
      </div>
      <el-icon :size="44"><Opportunity /></el-icon>
    </div>

    <div class="toolbar">
      <div class="toolbar-left">
        <el-select v-model="semester" clearable placeholder="全部学期" style="width: 260px" @change="load">
          <el-option v-for="s in semesters" :key="s.code" :label="`${s.code} · ${s.display_name}`" :value="s.code" />
        </el-select>
        <el-button :icon="Refresh" @click="load">刷新推荐</el-button>
        <el-button type="primary" plain :icon="ChatDotRound" @click="router.push('/my/chat')">让 AI 帮我规划</el-button>
      </div>
      <span class="muted">共 {{ rows.length }} 门候选课程</span>
    </div>

    <el-card shadow="never" class="panel-card">
      <el-table v-loading="loading" :data="rows" stripe>
        <el-table-column label="课程" min-width="240">
          <template #default="{ row }">
            <div class="course-name">{{ row.course_name }}</div>
            <div class="muted">{{ row.course_code }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="credits" label="学分" width="90" />
        <el-table-column prop="category_name" label="补足类别" min-width="180" />
        <el-table-column label="开课学期" min-width="220">
          <template #default="{ row }">
            <el-tag v-for="s in row.semesters" :key="s" class="semester-tag" effect="plain">{{ s }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="推荐理由" min-width="220" />
      </el-table>
      <el-empty v-if="!loading && rows.length === 0" description="当前没有满足条件的推荐课程" />
    </el-card>
  </div>
</template>
