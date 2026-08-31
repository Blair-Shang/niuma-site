<script setup lang="ts">
import { ref } from 'vue'
import { RsButton, RsInput } from '@/ui'
import { lookupSiteFeedback, submitSiteFeedback, type FeedbackPublic } from '../api/feedback'

const category = ref('bug')
const title = ref('')
const body = ref('')
const contact = ref('')
const loading = ref(false)
const error = ref('')
const done = ref(false)
const ticketId = ref('')

const lookupId = ref('')
const lookupContact = ref('')
const lookupLoading = ref(false)
const lookupError = ref('')
const lookupResult = ref<FeedbackPublic | null>(null)

async function submit() {
  error.value = ''
  done.value = false
  ticketId.value = ''
  loading.value = true
  try {
    const res = await submitSiteFeedback({
      category: category.value,
      title: title.value.trim(),
      body: body.value.trim(),
      contact: contact.value.trim(),
      product: 'niuma',
    })
    done.value = true
    ticketId.value = res.id
    lookupId.value = res.id
    lookupContact.value = contact.value.trim()
    title.value = ''
    body.value = ''
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'submit_failed'
  } finally {
    loading.value = false
  }
}

async function lookup() {
  lookupError.value = ''
  lookupResult.value = null
  lookupLoading.value = true
  try {
    lookupResult.value = await lookupSiteFeedback(lookupId.value, lookupContact.value)
  } catch (e) {
    lookupError.value = e instanceof Error ? e.message : 'lookup_failed'
  } finally {
    lookupLoading.value = false
  }
}
</script>

<template>
  <div class="site-section">
    <div class="site-container">
      <header class="site-page-head">
        <p class="site-section__eyebrow">Support</p>
        <h1>问题反馈</h1>
        <p>
          帮助我们改进 NiuMa。提交前可先查阅
          <RouterLink to="/faq">常见问题</RouterLink>；提交后请保存工单号，可在下方查询官方回复（桌面登录用户请在客户端「我的反馈」查看）。
        </p>
      </header>

      <form class="feedback-form" @submit.prevent="submit">
        <label class="feedback-field">
          <span>类型</span>
          <select v-model="category">
            <option value="bug">缺陷报告</option>
            <option value="feature">功能建议</option>
            <option value="other">其他</option>
          </select>
        </label>
        <label class="feedback-field">
          <span>标题</span>
          <RsInput v-model="title" placeholder="简要描述问题" required />
        </label>
        <label class="feedback-field">
          <span>详情</span>
          <textarea
            v-model="body"
            rows="8"
            required
            placeholder="复现步骤、期望结果、系统与 NiuMa 版本（勿粘贴敏感数据）"
          />
        </label>
        <label class="feedback-field">
          <span>联系邮箱</span>
          <RsInput v-model="contact" type="email" placeholder="查询回复时需与此一致" required />
        </label>
        <p v-if="error" class="feedback-error">提交失败：{{ error }}</p>
        <p v-if="done" class="feedback-ok">
          已收到，感谢反馈。
          <template v-if="ticketId">
            请保存工单号：<code>{{ ticketId }}</code>
          </template>
        </p>
        <RsButton variant="primary" radius="md" :disabled="loading" @click="submit">
          {{ loading ? '提交中…' : '提交反馈' }}
        </RsButton>
      </form>

      <section class="feedback-lookup">
        <h2>查询回复</h2>
        <p class="feedback-lookup__hint">使用提交时的工单号与联系邮箱查看处理状态与官方回复。</p>
        <div class="feedback-form">
          <label class="feedback-field">
            <span>工单号</span>
            <RsInput v-model="lookupId" placeholder="如 fb_xxx" />
          </label>
          <label class="feedback-field">
            <span>联系邮箱</span>
            <RsInput v-model="lookupContact" type="email" placeholder="提交时填写的邮箱" />
          </label>
          <p v-if="lookupError" class="feedback-error">查询失败：{{ lookupError }}</p>
          <RsButton variant="secondary" radius="md" :disabled="lookupLoading" @click="lookup">
            {{ lookupLoading ? '查询中…' : '查询' }}
          </RsButton>
          <article v-if="lookupResult" class="feedback-result">
            <h3>{{ lookupResult.title }}</h3>
            <p class="feedback-result__meta">状态：{{ lookupResult.status }} · {{ lookupResult.createdAt }}</p>
            <p class="feedback-result__body">{{ lookupResult.body }}</p>
            <div v-if="lookupResult.staffReply" class="feedback-result__reply">
              <strong>官方回复</strong>
              <p>{{ lookupResult.staffReply }}</p>
              <p v-if="lookupResult.staffReplyAt" class="feedback-result__meta">
                {{ lookupResult.staffReplyAt }}
              </p>
            </div>
            <p v-else class="feedback-result__meta">暂无官方回复，请稍后再查。</p>
          </article>
        </div>
      </section>

      <div class="site-prose" style="margin-top: 2.75rem">
        <h2>提交前建议准备</h2>
        <ul>
          <li>NiuMa 版本号（关于对话框或安装包文件名）</li>
          <li>操作系统版本</li>
          <li>可复现步骤；涉及数据库时说明引擎类型（无需粘贴敏感数据）</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.feedback-form {
  display: grid;
  gap: 1rem;
  max-width: 40rem;
  margin-top: 1.5rem;
}

.feedback-field {
  display: grid;
  gap: 0.45rem;
}

.feedback-field > span {
  font-size: 0.9rem;
  color: var(--rs-muted);
}

.feedback-field select,
.feedback-field textarea {
  width: 100%;
  border-radius: 0.65rem;
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  background: #14171f;
  color: inherit;
  padding: 0.65rem 0.75rem;
  font: inherit;
}

.feedback-error {
  color: #f87171;
  margin: 0;
}

.feedback-ok {
  color: #4ade80;
  margin: 0;
}

.feedback-ok code {
  color: #e2e8f0;
  background: #1e2430;
  padding: 0.1rem 0.35rem;
  border-radius: 0.3rem;
  font-size: 0.9em;
}

.feedback-lookup {
  margin-top: 3rem;
  max-width: 40rem;
}

.feedback-lookup h2 {
  margin: 0 0 0.5rem;
  font-size: 1.25rem;
}

.feedback-lookup__hint {
  margin: 0;
  color: var(--rs-muted);
  font-size: 0.92rem;
}

.site-page-head :deep(a) {
  color: color-mix(in srgb, var(--site-aurora-a) 80%, #fff);
  text-decoration: underline;
  text-underline-offset: 0.15em;
}

.feedback-result {
  border: 1px solid color-mix(in srgb, var(--rs-border) 90%, transparent);
  border-radius: 0.75rem;
  padding: 1rem 1.1rem;
  display: grid;
  gap: 0.55rem;
}

.feedback-result h3 {
  margin: 0;
  font-size: 1.05rem;
}

.feedback-result__meta {
  margin: 0;
  color: var(--rs-muted);
  font-size: 0.85rem;
}

.feedback-result__body {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  color: #cbd5e1;
  font-size: 0.92rem;
}

.feedback-result__reply {
  margin-top: 0.25rem;
  padding: 0.75rem 0.85rem;
  border-radius: 0.55rem;
  background: color-mix(in srgb, #3b82f6 16%, transparent);
  display: grid;
  gap: 0.35rem;
}

.feedback-result__reply p {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
