import { createRouter, createWebHistory } from 'vue-router'
import DebugView from "../views/DebugView.vue";

const routes = [
  {
    path: "/",
    name: "home",
    component: DebugView,
  },
  {
    path: "/debug",
    name: "terminal",
    component: DebugView,
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
