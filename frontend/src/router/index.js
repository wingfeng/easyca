import Vue from "vue";
import VueRouter from "vue-router";

const Layout = () => import( /* webpackChunkName: "layout" */ "@/layout")
const Home = () => import( /* webpackChunkName: "Layouut" */ "@/views/home/index.vue")

const CaServer = () => import( /* webpackChunkName: 'Newca'*/ "@/views/newca/server")
const CaPersonal = () => import( /* webpackChunkName: 'Newca'*/ "@/views/newca/personal")
const CaSign = () => import( /* webpackChunkName: 'Newca'*/ "@/views/newca/sign")

const User = () => import( /* webpackChunkName: 'User'*/ "@/views/setting/user")

const Error = () => import( /* webpackChunkName: 'Error'*/ "@/views/error/404")

const Help = () => import( /* webpackChunkName: 'Help'*/ "@/views/help")

const System = () => import( /* webpackChunkName: 'System'*/ "@/views/system")

Vue.use(VueRouter);

const routes = [{
    path: '/',
    redirect: '/v1/home'
  },
  {
    path: "/layout",
    name: "Layout",
    component: Layout,
    redirect: 'v1/home',
    children: [{
        path: '/v1/home',
        component: Home,
        name: 'home',
        meta: {
          title: '首页'
        }
      },
      {
        path: '/v1/newca/server',
        component: CaServer,
        meta: {
          title: '服务器证书'
        }
      },
      {
        path: '/v1/newca/personal',
        component: CaPersonal,
        meta: {
          title: '个人证书'
        }
      },
      {
        path: '/v1/newca/sign',
        component: CaSign,
        meta: {
          title: '签名证书'
        }
      },
      {
        path: '/v1/setting/user',
        component: User,
        meta: {
          title: '用户设置'
        }
      },
      {
        path: '/v1/help',
        component: Help,
        meta: {
          title: '使用帮助'
        }
      },
      {
        path: '/v1/system',
        component: System,
        meta: {
          title: '系统信息'
        }
      }
    ]
  },
  {
    path: '*',
    component: Error,
    meta: {
      title: '404页面'
    }
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

// set the doucment title by router navigation guard
router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }

  next()
})

export default router;
