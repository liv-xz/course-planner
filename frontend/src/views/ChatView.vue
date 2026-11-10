<script setup>
import { nextTick, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Promotion, Plus, Delete, ChatDotRound } from '@element-plus/icons-vue'
import http from '../api/http'

const sessions = ref([])
const currentSession = ref(null)
const messages = ref([])
const input = ref('')
const sending = ref(false)
const mode = ref('')
const toolStatus = ref('')
const scroller = ref()

const examples = [
  '我还差哪些学分？',
  '下学期推荐什么课？',
  '数据结构我能选吗？',
  '重修成绩怎么计算？',
]

async function loadSessions() {
  try {
    sessions.value = await http.get('/chat/sessions')
  } catch {
    // ignore
  }
}

async function loadStatus() {
  try {
    const data = await http.get('/chat/status')
    mode.value = data.mode
  } catch {
    mode.value = ''
  }
}

async function openSession(session) {
  currentSession.value = session
  try {
    const list = await http.get(`/chat/sessions/${session.id}/messages`)
    messages.value = (list || []).map((m) => ({ role: m.role, content: m.content || '' }))
    await scrollBottom()
  } catch {
    messages.value = []
  }
}

function newSession() {
  currentSession.value = null
  messages.value = []
  toolStatus.value = ''
}

async function removeSession(session, event) {
  event?.stopPropagation()
  try {
    await ElMessageBox.confirm('确定删除这个会话吗？', '删除确认', { type: 'warning' })
    await http.delete(`/chat/sessions/${session.id}`)
    if (currentSession.value?.id === session.id) newSession()
    await loadSessions()
  } catch {
    // cancelled or handled
  }
}

async function scrollBottom() {
  await nextTick()
  if (scroller.value) scroller.value.scrollTop = scroller.value.scrollHeight
}

function parseSSEBlock(block, assistant) {
  let event = 'message'
  const dataLines = []
  for (const line of block.split('\n')) {
    if (line.startsWith('event:')) event = line.slice(6).trim()
    else if (line.startsWith('data:')) dataLines.push(line.slice(5).trim())
  }
  if (dataLines.length === 0) return
  let data
  try {
    data = JSON.parse(dataLines.join('\n'))
  } catch {
    data = { text: dataLines.join('\n') }
  }
  if (event === 'meta') {
    if (data.session_id) {
      currentSession.value = { id: data.session_id, title: input.value || '新会话' }
    }
  } else if (event === 'tool_call') {
    toolStatus.value = `正在调用：${data.name}`
    assistant.tools.push(data.name)
  } else if (event === 'tool_result') {
    toolStatus.value = `已完成：${data.name}`
  } else if (event === 'delta') {
    assistant.content += data.text || ''
    scrollBottom()
  } else if (event === 'error') {
    assistant.content += `\n[错误] ${data.message || '对话失败'}`
    sending.value = false
  } else if (event === 'done') {
    sending.value = false
  }
}

async function send(text = input.value) {
  const content = String(text || '').trim()
  if (!content || sending.value) return
  sending.value = true
  toolStatus.value = ''
  messages.value.push({ role: 'user', content })
  const assistant = { role: 'assistant', content: '', tools: [] }
  messages.value.push(assistant)
  input.value = ''
  await scrollBottom()
  try {
    const token = localStorage.getItem('sjqu_token')
    const resp = await fetch('/api/v1/chat/completions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        message: content,
        session_id: currentSession.value?.id || 0,
      }),
    })
    if (!resp.ok || !resp.body) {
      const body = await resp.json().catch(() => ({}))
      throw new Error(body.message || `请求失败 (${resp.status})`)
    }
    const reader = resp.body.getReader()
    const decoder = new TextDecoder('utf-8')
    let buffer = ''
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true })
      const blocks = buffer.split('\n\n')
      buffer = blocks.pop() || ''
      for (const block of blocks) parseSSEBlock(block, assistant)
      await scrollBottom()
    }
    if (buffer.trim()) parseSSEBlock(buffer, assistant)
  } catch (e) {
    assistant.content += `\n[错误] ${e.message || '对话失败'}`
    ElMessage.error(e.message || '对话失败')
  } finally {
    sending.value = false
    toolStatus.value = ''
    await loadSessions()
    await scrollBottom()
  }
}

function onKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}

onMounted(async () => {
  await Promise.all([loadSessions(), loadStatus()])
})
</script>

<template>
  <div class="chat-page">
    <aside class="chat-side">
      <el-button type="primary" :icon="Plus" class="new-chat" @click="newSession">新建会话</el-button>
      <div class="session-list">
        <div
          v-for="s in sessions"
          :key="s.id"
          class="session-item"
          :class="{ active: currentSession?.id === s.id }"
          @click="openSession(s)"
        >
          <el-icon><ChatDotRound /></el-icon>
          <span class="session-title">{{ s.title || '未命名会话' }}</span>
          <el-button text size="small" :icon="Delete" @click="removeSession(s, $event)" />
        </div>
        <el-empty v-if="sessions.length === 0" description="还没有会话" :image-size="70" />
      </div>
    </aside>

    <section class="chat-main">
      <div class="chat-header">
        <div>
          <div class="panel-title">选课规划助手</div>
          <div class="muted">
            {{ mode === 'llm' ? '大模型 + 工具调用' : '规则兜底模式（在 config.yaml 填入 LLM API Key 后可启用大模型）' }}
          </div>
        </div>
        <el-tag :type="mode === 'llm' ? 'success' : 'info'" effect="light">
          {{ mode === 'llm' ? 'LLM' : 'RULE' }}
        </el-tag>
      </div>

      <div ref="scroller" class="chat-scroll">
        <div v-if="messages.length === 0" class="chat-welcome">
          <div class="welcome-icon"><el-icon :size="34"><Promotion /></el-icon></div>
          <h2>我可以帮你规划选课</h2>
          <p>查询学分缺口、筛选候选课程、检查先修关系、解释培养方案规则。</p>
          <div class="example-list">
            <el-button v-for="e in examples" :key="e" plain @click="send(e)">{{ e }}</el-button>
          </div>
        </div>

        <div v-for="(m, idx) in messages" :key="idx" class="msg-row" :class="m.role">
          <div class="msg-avatar">{{ m.role === 'user' ? '我' : 'AI' }}</div>
          <div class="msg-bubble">
            <div v-if="m.tools?.length" class="tool-tags">
              <el-tag v-for="t in m.tools" :key="t" size="small" effect="plain">{{ t }}</el-tag>
            </div>
            <div class="msg-content">{{ m.content || (m.role === 'assistant' && sending ? '正在思考…' : '') }}</div>
          </div>
        </div>
      </div>

      <div v-if="toolStatus" class="tool-status">{{ toolStatus }}</div>

      <div class="chat-input">
        <el-input
          v-model="input"
          type="textarea"
          :autosize="{ minRows: 2, maxRows: 5 }"
          resize="none"
          placeholder="例如：我还差哪些学分？下学期推荐什么课？"
          @keydown="onKeydown"
        />
        <el-button type="primary" :icon="Promotion" :loading="sending" @click="send()">发送</el-button>
      </div>
    </section>
  </div>
</template>
