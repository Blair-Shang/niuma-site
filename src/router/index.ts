import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../pages/HomePage.vue'),
      meta: { title: 'NiuMa — 专业开发与运维工具' },
    },
    {
      path: '/products',
      name: 'products',
      component: () => import('../pages/ProductsPage.vue'),
      meta: { title: '产品 — NiuMa' },
    },
    {
      path: '/products/niuma',
      name: 'product-niuma',
      component: () => import('../pages/ProductNiumaPage.vue'),
      meta: { title: 'NiuMa 桌面端 — NiuMa' },
    },
    {
      path: '/download',
      name: 'download',
      component: () => import('../pages/DownloadPage.vue'),
      meta: { title: '下载 — NiuMa' },
    },
    {
      path: '/feedback',
      name: 'feedback',
      component: () => import('../pages/FeedbackPage.vue'),
      meta: { title: '问题反馈 — NiuMa' },
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../pages/AboutPage.vue'),
      meta: { title: '关于 — NiuMa' },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

router.afterEach((to) => {
  const title = typeof to.meta.title === 'string' ? to.meta.title : 'NiuMa'
  document.title = title
})
