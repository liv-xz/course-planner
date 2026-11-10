import { createRouter, createWebHistory } from 'vue-router'
import { auth } from '../stores/auth'
import MainLayout from '../layouts/MainLayout.vue'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import PlansView from '../views/PlansView.vue'
import CategoriesView from '../views/CategoriesView.vue'
import CoursesView from '../views/CoursesView.vue'
import StudentsView from '../views/StudentsView.vue'
import MajorsView from '../views/MajorsView.vue'
import SemestersView from '../views/SemestersView.vue'
import MyRecordsView from '../views/MyRecordsView.vue'
import AuditView from '../views/AuditView.vue'
import RecommendationsView from '../views/RecommendationsView.vue'
import ChatView from '../views/ChatView.vue'
import StatsView from '../views/StatsView.vue'
import SecurityView from '../views/SecurityView.vue'

const routes = [
  { path: '/login', name: 'login', component: LoginView, meta: { public: true } },
  {
    path: '/',
    component: MainLayout,
    children: [
      { path: '', redirect: () => (auth.isAdmin() ? '/dashboard' : '/my/audit') },
      { path: 'dashboard', name: 'dashboard', component: DashboardView, meta: { role: 'admin' } },
      { path: 'plans', name: 'plans', component: PlansView, meta: { role: 'admin' } },
      { path: 'categories', name: 'categories', component: CategoriesView, meta: { role: 'admin' } },
      { path: 'courses', name: 'courses', component: CoursesView, meta: { role: 'admin' } },
      { path: 'students', name: 'students', component: StudentsView, meta: { role: 'admin' } },
      { path: 'stats', name: 'stats', component: StatsView, meta: { role: 'admin' } },
      { path: 'security', name: 'security', component: SecurityView, meta: { role: 'admin' } },
      { path: 'majors', name: 'majors', component: MajorsView, meta: { role: 'admin' } },
      { path: 'semesters', name: 'semesters', component: SemestersView, meta: { role: 'admin' } },
      { path: 'my/records', name: 'my-records', component: MyRecordsView, meta: { role: 'student' } },
      { path: 'my/audit', name: 'my-audit', component: AuditView, meta: { role: 'student' } },
      { path: 'my/recommendations', name: 'my-recommendations', component: RecommendationsView, meta: { role: 'student' } },
      { path: 'my/chat', name: 'my-chat', component: ChatView, meta: { role: 'student' } },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  if (to.meta.public) {
    if (auth.state.token && auth.state.user) {
      return auth.isAdmin() ? '/dashboard' : '/my/audit'
    }
    return true
  }
  if (!auth.state.token) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (!auth.state.user) {
    try {
      await auth.fetchProfile()
    } catch {
      auth.logout()
      return { path: '/login' }
    }
  }
  if (to.meta.role && auth.state.user?.role !== to.meta.role) {
    return auth.isAdmin() ? '/dashboard' : '/my/audit'
  }
  return true
})

export default router
