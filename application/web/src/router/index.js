import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/uplink',
    children: [{
      path: 'uplink',
      name: 'Uplink',
      component: () => import('@/views/uplink/index'),
      meta: { title: '报价信息录入', icon: 'el-icon-edit-outline', roles: ['普通用户'] }
    }]
  },

  {
    path: '/',
    component: Layout,
    redirect: '/approve',
    children: [{
      path: 'approve',
      name: 'Approve',
      component: () => import('@/views/approve/index'),
      meta: { title: '用户资质审核', icon: 'el-icon-edit-outline', roles: ['管理员'] }
    }]
  },
  {
    path: '/adtrace',
    component: Layout,
    //redirect: '/adtrace',
    children: [{
      path: 'adtrace',
      name: 'Adtrace',
      component: () => import('@/views/adtrace/index'),
      meta: { title: '管理员溯源信息查询', icon: 'el-icon-search', roles: ['管理员'] }
    }]
  },
  {
    path: '/trace',
    component: Layout,
    //redirect: '/trace',
    children: [{
      path: 'trace',
      name: 'Trace',
      component: () => import('@/views/trace/index'),
      meta: { title: '用户溯源信息查询', icon: 'el-icon-search', roles: ['普通用户'] }
    }]
  },

  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: 'http://119.45.247.29:8080',
        meta: { title: '区块链浏览器', icon: 'el-icon-discover' }
      }
    ]
  },

  { path: '*', redirect: '/404', hidden: true }
]

// 定义不同用户类型的路由
export const userRoutes = [
  
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/uplink',
    children: [{
      path: 'uplink',
      name: 'Uplink',
      component: () => import('@/views/uplink/index'),
      meta: { title: '报价信息录入', icon: 'el-icon-edit-outline', roles: ['普通用户'] }
    }]
  },
  // {
  //   path: '/uplink',
  //   name: 'Uplink',
  //   component: () => import('@/views/uplink/index'),
  //   meta: { title: '溯源信息录入', icon: 'el-icon-edit-outline' }
  // },
  {
    path: '/trace',
    component: Layout,
    //redirect: '/trace',
    children: [{
      path: 'trace',
      name: 'Trace',
      component: () => import('@/views/trace/index'),
      meta: { title: '用户溯源信息查询', icon: 'el-icon-search', roles: ['普通用户'] }
    }]
  },
  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: 'http://192.168.15.132:8080',
        meta: { title: '区块链浏览器', icon: 'el-icon-discover' }
      }
    ]
  },

  { path: '*', redirect: '/404', hidden: true }
]

export const adminRoutes = [
  
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/approve',
    children: [{
      path: 'approve',
      name: 'Approve',
      component: () => import('@/views/approve/index'),
      meta: { title: '用户资质审核', icon: 'el-icon-edit-outline', roles: ['管理员'] }
    }]
  },
  // {
  //   path: '/approve',
  //   name: 'Approve',
  //   component: () => import('@/views/approve/index'),
  //   meta: { title: '用户资质审核', icon: 'el-icon-edit-outline' }
  // },
  {
    path: '/adtrace',
    component: Layout,
    //redirect: '/adtrace',
    children: [{
      path: 'adtrace',
      name: 'Adtrace',
      component: () => import('@/views/adtrace/index'),
      meta: { title: '管理员溯源信息查询', icon: 'el-icon-search', roles: ['管理员'] }
    }]
  },
  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: 'http://192.168.15.132:8080',
        meta: { title: '区块链浏览器', icon: 'el-icon-discover' }
      }
    ]
  },

  { path: '*', redirect: '/404', hidden: true }
]

// 获取当前用户角色，这里假设从 localStorage 中获取
console.log('当前用户角色:', userRole); // 添加调试信息
const userRole = localStorage.getItem('userType') || '普通用户';
// 根据用户角色选择路由配置
let routes = [];
if (userRole === '管理员') {
  routes = adminRoutes;
} else {
  routes = userRoutes;
}
console.log('当前加载的路由配置:', routes); // 添加调试信息
const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: routes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router




// 根据用户类型生成路由
// export function generateRoutes(userType) {
//   let routes = []
//   switch (userType) {
//     case '普通用户':
//       routes = userRoutes
//       break
//     case '管理员':
//       routes = adminRoutes
//       break
//     // 其他用user户类型...
//     default:
//       routes = []
//   }
//   return [
//     {
//       path: '/',
//       component: Layout,
//       redirect: routes.length > 0 ? routes[0].path : '/',
//       children: routes
//     }
//   ]
// }
