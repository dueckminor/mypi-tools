import { createRouter, createWebHistory } from 'vue-router'
import HostsView from "../views/HostsView.vue";
import WizardSD from "../views/WizardSD.vue";
import CertificatesView from "../views/CertificatesView.vue";
import InitializedSD from "../views/InitializeSD.vue";
import TerminalView from "../views/TerminalView.vue";
import ActionView from "../views/ActionView.vue";
import AboutView from "../views/AboutView.vue";

const routes = [
  {
    path: "/",
    name: "home",
    component: AboutView,
  },
  {
    path: "/hosts",
    name: "hosts",
    component: HostsView,
  },
  {
    path: "/certificates",
    name: "certificates",
    component: CertificatesView,
  },
  {
    path: "/hosts/:host/wizardsd",
    name: "wizardsd",
    component: WizardSD,
  },
  {
    path: "/initializesd",
    name: "initializesd",
    component: InitializedSD,
  },
  {
    path: "/hosts/:host/terminal",
    name: "terminal",
    component: TerminalView,
  },
  {
    path: "/hosts/:host/actions/:action",
    name: "setup",
    component: ActionView,
  },
  {
    path: "/about",
    name: "about",
    component: AboutView,
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
