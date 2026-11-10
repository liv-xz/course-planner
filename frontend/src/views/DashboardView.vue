<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Collection, Reading, User, School, Calendar, Grid } from '@element-plus/icons-vue'
import http from '../api/http'

const router = useRouter()
const loading = ref(true)
const stats = ref({ plans: 0, courses: 0, students: 0, majors: 0, semesters: 0 })

async function load() {
  loading.value = true
  try {
    const [plans, courses, students, majors, semesters] = await Promise.all([
      http.get('/plans', { params: { page: 1, size: 1 } }),
      http.get('/courses', { params: { page: 1, size: 1 } }),
      http.get('/students', { params: { page: 1, size: 1 } }),
      http.get('/majors', { params: { page: 1, size: 1 } }),
      http.get('/semesters', { params: { page: 1, size: 1 } }),
    ])
    stats.value = {
      plans: plans.total || 0,
      courses: courses.total || 0,
      students: students.total || 0,
      majors: majors.total || 0,
      semesters: semesters.total || 0,
    }
  } finally {
    loading.value = false
  }
}

onMounted(load)

const cards = [
  { key: 'plans', label: '培养方案', icon: Collection, color: '#2563eb', path: '/plans' },
  { key: 'courses', label: '课程库', icon: Reading, color: '#0f766e', path: '/courses' },
  { key: 'students', label: '学生', icon: User, color: '#c2410c', path: '/students' },
  { key: 'majors', label: '专业', icon: School, color: '#7c3aed', path: '/majors' },
  { key: 'semesters', label: '学期', icon: Calendar, color: '#be123c', path: '/semesters' },
]
</script>

<template>
  <div v-loading="loading">
    <div class="hero-card">
      <div>
        <div class="hero-kicker">ACADEMIC PLANNER</div>
        <h2>把培养方案、课程与学分审核放在同一个工作台</h2>
        <p>维护基础数据，配置课程归属，管理学生修读记录；学生端可以查看学分缺口和选课推荐。</p>
      </div>
      <el-button type="primary" size="large" @click="router.push('/courses')">开始维护课程</el-button>
    </div>

    <el-row :gutter="18" class="stat-row">
      <el-col v-for="card in cards" :key="card.key" :xs="24" :sm="12" :md="8" :lg="4">
        <div class="stat-card clickable" @click="router.push(card.path)">
          <div class="stat-icon" :style="{ background: card.color + '18', color: card.color }">
            <el-icon :size="22"><component :is="card.icon" /></el-icon>
          </div>
          <div class="stat-value">{{ stats[card.key] }}</div>
          <div class="stat-label">{{ card.label }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="18">
      <el-col :xs="24" :lg="14">
        <el-card shadow="never" class="panel-card">
          <template #header><div class="panel-title">常用操作</div></template>
          <div class="quick-grid">
            <div class="quick-item" @click="router.push('/plans')">
              <el-icon><Collection /></el-icon>
              <div><b>新建培养方案</b><span>按专业和年级配置总学分</span></div>
            </div>
            <div class="quick-item" @click="router.push('/categories')">
              <el-icon><Grid /></el-icon>
              <div><b>维护课程类别</b><span>树形结构，支持父级汇总</span></div>
            </div>
            <div class="quick-item" @click="router.push('/courses')">
              <el-icon><Reading /></el-icon>
              <div><b>课程归属与开课</b><span>挂类别、配先修、设学期</span></div>
            </div>
            <div class="quick-item" @click="router.push('/students')">
              <el-icon><User /></el-icon>
              <div><b>学生修读记录</b><span>录入成绩，实时看审核</span></div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="10">
        <el-card shadow="never" class="panel-card">
          <template #header><div class="panel-title">数据口径</div></template>
          <ul class="rule-list">
            <li>课程只能挂在叶子类别下，同一方案内一门课只保留一个归属。</li>
            <li>学分审核按课程去重，同一门课重修只取最高分。</li>
            <li>推荐只包含有开课记录的课程，先修未通过会被过滤。</li>
            <li>课程删除为逻辑删除，历史记录仍可查询。</li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
