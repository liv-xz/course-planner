import { reactive } from 'vue'
import http from '../api/http'

const state = reactive({
  token: localStorage.getItem('sjqu_token') || '',
  user: null,
})

async function login(username, password) {
  const data = await http.post('/auth/login', { username, password })
  state.token = data.token
  state.user = {
    user_id: data.user_id,
    username: data.username,
    name: data.name,
    role: data.role,
  }
  localStorage.setItem('sjqu_token', data.token)
  return data
}

async function fetchProfile() {
  if (!state.token) return null
  const data = await http.get('/auth/profile')
  state.user = data
  return data
}

function logout() {
  state.token = ''
  state.user = null
  localStorage.removeItem('sjqu_token')
}

function isAdmin() {
  return state.user?.role === 'admin'
}

function isStudent() {
  return state.user?.role === 'student'
}

export const auth = { state, login, fetchProfile, logout, isAdmin, isStudent }
