/**
 * 官网常见问题登记（单一数据源）。
 * 新增条目：在此数组追加一条，填写 category / question / answer / tags；
 * 安装类问题会同步出现在下载页摘要中。
 */

export type FaqCategory = 'install' | 'update' | 'product' | 'account'

export type FaqAnswerBlock =
  | { kind: 'text'; html: string }
  | { kind: 'list'; ordered?: boolean; items: string[] }

export type FaqItem = {
  id: string
  category: FaqCategory
  question: string
  answer: FaqAnswerBlock[]
  /** 搜索关键词（小写）；未填时仅匹配标题与正文 */
  tags?: string[]
}

export const faqCategories: { id: FaqCategory; label: string }[] = [
  { id: 'install', label: '安装与下载' },
  { id: 'update', label: '更新' },
  { id: 'product', label: '产品与系统' },
  { id: 'account', label: '账户与反馈' },
]

/** 按展示顺序排列；同分类内建议把最常见的问题放前面 */
export const faqItems: FaqItem[] = [
  {
    id: 'windows-smartscreen',
    category: 'install',
    question: 'Windows 提示「已保护你的电脑」或「是否继续运行」',
    tags: ['smartscreen', 'defender', '未知发布者', '仍要运行', '签名'],
    answer: [
      {
        kind: 'text',
        html: '这是 Microsoft Defender SmartScreen 对未签名或新发布者的常规拦截，不代表安装包被篡改。',
      },
      {
        kind: 'list',
        ordered: true,
        items: [
          '在 SmartScreen 窗口点击 <strong>更多信息</strong>（或「详细信息」）。',
          '再点 <strong>仍要运行</strong>（Run anyway），即可启动安装程序。',
        ],
      },
      {
        kind: 'text',
        html: '若浏览器下载时提示「可能有害」，请选择 <strong>保留</strong> 或 <strong>仍要下载</strong>，安装包仅来自本站官方渠道。',
      },
    ],
  },
  {
    id: 'no-code-signing',
    category: 'install',
    question: '为什么没有企业签名？可以继续安装吗？',
    tags: ['企业签名', 'ev', 'authenticode', 'notarization', 'gatekeeper', 'sha256', '哈希', '校验'],
    answer: [
      {
        kind: 'text',
        html: '各平台对「未知名发布者」都有类似机制：Windows SmartScreen、macOS Gatekeeper（开发者 ID + 公证）、Linux 则多见于包管理器策略或企业环境限制。我们尚未完成 Apple / Windows 的正式签名与公证流程。',
      },
      {
        kind: 'text',
        html: '在证书到位前，请通过以下方式确认文件可信：',
      },
      {
        kind: 'list',
        items: [
          '只从本站 <a href="/download">下载页</a> 获取安装包，勿使用第三方镜像或网盘转载。',
          '下载完成后，用下载页卡片中的 SHA-256 与本地文件哈希比对（Windows PowerShell：<code>Get-FileHash .\\Setup.exe</code>；macOS / Linux：<code>shasum -a 256 文件名</code>）。',
        ],
      },
      {
        kind: 'text',
        html: '校验一致即可按各平台说明继续安装；正式签名 / 公证上线后，系统拦截会逐步减少。',
      },
    ],
  },
  {
    id: 'windows-uac',
    category: 'install',
    question: '安装时弹出「是否允许此应用对你的设备进行更改」',
    tags: ['uac', '管理员', '权限'],
    answer: [
      {
        kind: 'text',
        html: '这是 Windows 用户账户控制（UAC）的标准提示。NiuMa 安装程序需要写入 Program Files 并创建开始菜单快捷方式，请选择 <strong>是</strong>。',
      },
      {
        kind: 'text',
        html: '若你使用的是标准用户账户且无法授权，请联系系统管理员或使用具备管理员权限的账户。',
      },
    ],
  },
  {
    id: 'system-requirements',
    category: 'install',
    question: 'Windows 安装有哪些系统要求？',
    tags: ['windows 10', 'windows 11', 'x64', '配置'],
    answer: [
      {
        kind: 'list',
        items: [
          '操作系统：Windows 10 或 Windows 11（64 位 / x64）。',
          '架构：当前正式包为 x64；请从下载页获取对应平台安装包。',
          '网络：首次启动与检查更新需要联网；离线使用核心功能请查阅后续版本说明。',
        ],
      },
    ],
  },
  {
    id: 'macos-gatekeeper',
    category: 'install',
    question: 'macOS 提示「无法打开，因为无法验证开发者」或「来自身份不明的开发者」',
    tags: ['macos', 'gatekeeper', '公证', 'notarization', 'quarantine', 'dmg', 'pkg'],
    answer: [
      {
        kind: 'text',
        html: 'macOS Gatekeeper 会拦截未公证或未签名的应用，与 Windows SmartScreen 类似，<strong>不代表文件被篡改</strong>。预览版目前尚未完成 Apple 开发者签名与公证。',
      },
      {
        kind: 'list',
        ordered: true,
        items: [
          '优先尝试：在 Finder 中 <strong>按住 Control 键点击</strong> 安装包或 .app，选择 <strong>打开</strong>，在对话框中再次确认打开。',
          '若仍被拦：打开 <strong>系统设置 → 隐私与安全性</strong>，在页面底部找到被拦应用的 <strong>仍要打开</strong>（Open Anyway）。',
          '浏览器下载的文件可能带隔离标记；终端可执行 <code>xattr -dr com.apple.quarantine /path/to/NiuMa.app</code> 后重试（将路径换成实际位置）。',
        ],
      },
      {
        kind: 'text',
        html: '.pkg 安装包安装时若要求输入密码，属于 macOS 写入「应用程序」文件夹的正常权限提示，请输入本机管理员密码。',
      },
    ],
  },
  {
    id: 'linux-install-permissions',
    category: 'install',
    question: 'Linux 安装需要额外权限吗？',
    tags: ['linux', 'chmod', 'sudo', 'deb', 'run', 'appimage', '权限'],
    answer: [
      {
        kind: 'text',
        html: 'Linux 没有像 Windows SmartScreen 或 macOS Gatekeeper 那样的统一桌面拦截，但安装方式不同，可能遇到<strong>可执行权限</strong>或 <strong>sudo</strong> 要求：',
      },
      {
        kind: 'list',
        items: [
          '<strong>.run 安装脚本</strong>：先 <code>chmod +x NiuMa-*.run</code>，再 <code>./NiuMa-*.run</code>；写入系统目录时安装器会提示输入 sudo 密码。',
          '<strong>.deb 包</strong>：使用 <code>sudo dpkg -i 包名.deb</code> 或图形化软件中心安装，需要管理员权限。',
          '若提示「Permission denied」，检查文件是否有执行权限，或是否在只读目录中运行。',
        ],
      },
      {
        kind: 'text',
        html: '企业环境若启用 SELinux / AppArmor 等策略，可能需管理员放行；此类情况请至 <a href="/feedback">问题反馈</a> 说明发行版与策略环境。',
      },
    ],
  },
  {
    id: 'linux-macos-beta',
    category: 'install',
    question: 'Linux 与 macOS 版本何时提供正式版？',
    tags: ['linux', 'macos', 'beta', '预览'],
    answer: [
      {
        kind: 'text',
        html: 'Linux 与 macOS 安装包目前为 <strong>预览版（Beta）</strong>，功能与稳定性仍在完善中，可在 <a href="/download">下载页</a> 获取最新预览包。',
      },
      {
        kind: 'text',
        html: '正式版发布时间将在更新日志与官网公告中同步，建议关注下载页版本说明。',
      },
    ],
  },
  {
    id: 'in-app-update',
    category: 'update',
    question: '如何在应用内检查更新？',
    tags: ['更新', '升级', 'changelog'],
    answer: [
      {
        kind: 'text',
        html: '打开 NiuMa，进入菜单中的 <strong>关于</strong> 或 <strong>更新日志</strong>，可查看当前版本并检查是否有新版本。',
      },
      {
        kind: 'text',
        html: '若有可用更新，按提示下载并安装；部分版本需重启应用后生效。历史说明见 <a href="/download">下载页更新日志</a>。',
      },
    ],
  },
  {
    id: 'update-failed',
    category: 'update',
    question: '自动更新失败怎么办？',
    tags: ['更新失败', '下载失败', '网络'],
    answer: [
      {
        kind: 'list',
        ordered: true,
        items: [
          '确认网络正常，关闭代理或 VPN 后重试。',
          '到 <a href="/download">下载页</a> 手动下载最新安装包覆盖安装（数据目录通常保留）。',
          '仍失败请至 <a href="/feedback">问题反馈</a> 附上版本号与错误截图。',
        ],
      },
    ],
  },
  {
    id: 'what-is-niuma',
    category: 'product',
    question: 'NiuMa 是什么？适合谁使用？',
    tags: ['介绍', '数据库', '运维', 'sql'],
    answer: [
      {
        kind: 'text',
        html: 'NiuMa 是面向开发与运维人员的专业桌面工作台，整合数据库管理、SQL、监控与 AI 辅助等能力，减少多工具切换成本。',
      },
      {
        kind: 'text',
        html: '详见 <a href="/products/niuma">产品介绍</a> 与 <a href="/about">关于我们</a>。',
      },
    ],
  },
  {
    id: 'offline-data',
    category: 'product',
    question: '连接数据库时数据会经过你们的服务器吗？',
    tags: ['隐私', '安全', '本地', '云端'],
    answer: [
      {
        kind: 'text',
        html: '数据库连接与查询在本地客户端执行；请勿在反馈或对话中粘贴生产环境密码、密钥等敏感信息。',
      },
      {
        kind: 'text',
        html: '使用云端账户、更新检查等功能时会与 NiuMa 云服务通信，具体范围以客户端内说明为准。',
      },
    ],
  },
  {
    id: 'submit-feedback',
    category: 'account',
    question: '如何提交问题或功能建议？',
    tags: ['反馈', 'bug', '工单', 'support'],
    answer: [
      {
        kind: 'text',
        html: '请前往 <a href="/feedback">问题反馈</a> 填写类型、标题与详情，并保存返回的工单号以便查询回复。',
      },
      {
        kind: 'text',
        html: '已登录桌面端的用户也可在客户端「我的反馈」中查看历史工单。',
      },
    ],
  },
  {
    id: 'feedback-info',
    category: 'account',
    question: '反馈时应准备哪些信息？',
    tags: ['日志', '版本', '复现'],
    answer: [
      {
        kind: 'list',
        items: [
          'NiuMa 版本号（关于对话框或安装包文件名）。',
          '操作系统与架构（如 Windows 11 x64）。',
          '可复现步骤；涉及数据库时说明引擎类型，勿粘贴敏感数据。',
        ],
      },
    ],
  },
  {
    id: 'install-still-blocked',
    category: 'install',
    question: '按说明操作后仍无法安装',
    tags: ['安装失败', '无法安装'],
    answer: [
      {
        kind: 'text',
        html: '请至 <a href="/feedback">问题反馈</a> 提交系统拦截截图、操作系统与版本号（如 Windows 11 / macOS 14 / Ubuntu 22.04）及安装包 SHA-256（若已校验），我们会协助排查。',
      },
    ],
  },
]

export function faqItemsByCategory(category: FaqCategory | 'all'): FaqItem[] {
  if (category === 'all') return faqItems
  return faqItems.filter((item) => item.category === category)
}

export function faqCategoryLabel(id: FaqCategory): string {
  return faqCategories.find((c) => c.id === id)?.label ?? id
}

/** 将条目正文展平为可搜索纯文本 */
export function faqSearchText(item: FaqItem): string {
  const parts: string[] = [item.question, ...(item.tags ?? [])]
  for (const block of item.answer) {
    if (block.kind === 'text') {
      parts.push(block.html.replace(/<[^>]+>/g, ' '))
    } else {
      parts.push(...block.items.map((s) => s.replace(/<[^>]+>/g, ' ')))
    }
  }
  return parts.join(' ').toLowerCase()
}
