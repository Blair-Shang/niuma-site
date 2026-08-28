import { siteConfig } from '../config/site'

export type FeedbackPayload = {
  category: string
  title: string
  body: string
  contact: string
  product?: string
  clientVersion?: string
}

export type FeedbackPublic = {
  id: string
  category: string
  title: string
  body: string
  contact?: string
  status: string
  staffReply?: string
  staffReplyAt?: string | null
  createdAt: string
  updatedAt?: string
}

function cloudBase(): string {
  return siteConfig.cloudApiBase.replace(/\/$/, '')
}

async function readError(res: Response): Promise<string> {
  let code = `http_${res.status}`
  try {
    const data = (await res.json()) as { error?: string }
    if (data.error) code = data.error
  } catch {
    /* ignore */
  }
  return code
}

/** 提交到 niuma-cloud（非本站 API） */
export async function submitSiteFeedback(payload: FeedbackPayload): Promise<{ id: string }> {
  const res = await fetch(`${cloudBase()}/api/v1/feedback`, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      category: payload.category,
      title: payload.title,
      body: payload.body,
      contact: payload.contact,
      product: payload.product || 'niuma',
      clientVersion: payload.clientVersion || '',
    }),
  })
  if (!res.ok) {
    throw new Error(await readError(res))
  }
  const data = (await res.json()) as { id?: string }
  if (!data.id) {
    throw new Error('missing_id')
  }
  return { id: data.id }
}

/** 官网匿名查询回复：工单 ID + 联系邮箱 */
export async function lookupSiteFeedback(id: string, contact: string): Promise<FeedbackPublic> {
  const qs = new URLSearchParams({
    id: id.trim(),
    contact: contact.trim(),
  })
  const res = await fetch(`${cloudBase()}/api/v1/feedback/lookup?${qs}`, {
    headers: { Accept: 'application/json' },
  })
  if (!res.ok) {
    throw new Error(await readError(res))
  }
  return (await res.json()) as FeedbackPublic
}
