<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import {
  DataBoard, Collection, Grid, Reading, User, School, Calendar, SwitchButton,
  List, TrendCharts, Opportunity, DataAnalysis, Lock, ChatDotRound,
} from '@element-plus/icons-vue'
import { auth } from '../stores/auth'

const route = useRoute()
const router = useRouter()

const isAdmin = computed(() => auth.state.user?.role === 'admin')

const titleMap = {
  dashboard: '工作台',
  plans: '培养方案',
  categories: '课程类别',
  courses: '课程库',
  students: '学生管理',
  majors: '专业管理',
  semesters: '学期管理',
  stats: '学分统计',
  security: '安全日志',
  'my-chat': '选课助手',
  'my-records': '我的修读记录',
  'my-audit': '学分审核',
  'my-recommendations': '选课推荐',
}
const pageTitle = computed(() => titleMap[route.name] || '课程规划系统')

async function logout() {
  try {
    await ElMessageBox.confirm('确定退出登录吗？', '提示', { type: 'warning' })
    auth.logout()
    router.push('/login')
  } catch {
    // cancelled
  }
}
</script>

<template>
  <el-container class="app-shell">
    <el-aside width="236px" class="sidebar">
      <div class="brand">
        <div class="brand-mark">JQ</div>
        <div>
          <div class="brand-title">课程规划系统</div>
          <div class="brand-sub">上海建桥学院</div>
        </div>
      </div>

      <el-menu :default-active="route.path" router class="side-menu">
        <template v-if="isAdmin">
          <el-menu-item index="/dashboard">
            <el-icon><DataBoard /></el-icon><span>工作台</span>
          </el-menu-item>
          <el-menu-item index="/plans">
            <el-icon><Collection /></el-icon><span>培养方案</span>
          </el-menu-item>
          <el-menu-item index="/categories">
            <el-icon><Grid /></el-icon><span>课程类别</span>
          </el-menu-item>
          <el-menu-item index="/courses">
            <el-icon><Reading /></el-icon><span>课程库</span>
          </el-menu-item>
          <el-menu-item index="/students">
            <el-icon><User /></el-icon><span>学生管理</span>
          </el-menu-item>
          <el-menu-item index="/majors">
            <el-icon><School /></el-icon><span>专业管理</span>
          </el-menu-item>
          <el-menu-item index="/semesters">
            <el-icon><Calendar /></el-icon><span>学期管理</span>
          </el-menu-item>
          <el-menu-item index="/stats">
            <el-icon><DataAnalysis /></el-icon><span>学分统计</span>
          </el-menu-item>
          <el-menu-item index="/security">
            <el-icon><Lock /></el-icon><span>安全日志</span>
          </el-menu-item>
        </template>
        <template v-else>
          <el-menu-item index="/my/audit">
            <el-icon><TrendCharts /></el-icon><span>学分审核</span>
          </el-menu-item>
          <el-menu-item index="/my/recommendations">
            <el-icon><Opportunity /></el-icon><span>选课推荐</span>
          </el-menu-item>
          <el-menu-item index="/my/chat">
            <el-icon><ChatDotRound /></el-icon><span>选课助手</span>
          </el-menu-item>
          <el-menu-item index="/my/records">
            <el-icon><List /></el-icon><span>修读记录</span>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="topbar">
        <div>
          <div class="page-title">{{ pageTitle }}</div>
          <div class="page-sub">选课规划智能体 · 基础管理平台</div>
        </div>
        <div class="user-area">
          <el-tag :type="isAdmin ? 'warning' : 'success'" effect="light">
            {{ isAdmin ? '教务管理员' : '学生' }}
          </el-tag>
          <span class="user-name">{{ auth.state.user?.name || auth.state.user?.username }}</span>
          <el-button text :icon="SwitchButton" @click="logout">退出</el-button>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>
