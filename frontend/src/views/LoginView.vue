<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { auth } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const form = reactive({ username: 'admin', password: '123456' })

async function submit() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const data = await auth.login(form.username, form.password)
    ElMessage.success('登录成功')
    const redirect = route.query.redirect
    if (redirect) {
      router.push(String(redirect))
    } else {
      router.push(data.role === 'admin' ? '/dashboard' : '/my/audit')
    }
  } catch {
    // http interceptor already shows message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-glow login-glow-a" />
    <div class="login-glow login-glow-b" />
    <div class="login-card">
      <div class="login-brand">
        <div class="login-mark">JQ</div>
        <div>
          <h1>课程规划系统</h1>
          <p>上海建桥学院 · 选课规划智能体</p>
        </div>
      </div>

      <el-form :model="form" size="large" @keyup.enter="submit">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名 / 学号" :prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" show-password placeholder="密码" :prefix-icon="Lock" />
        </el-form-item>
        <el-button type="primary" size="large" class="login-btn" :loading="loading" @click="submit">
          登录
        </el-button>
      </el-form>

      <div class="login-tips">
        <div>管理员：admin / 123456</div>
        <div>学生：2323500 / 123456</div>
      </div>
    </div>
  </div>
</template>
