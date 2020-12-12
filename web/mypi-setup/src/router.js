import Vue from "vue";
import Router from "vue-router";
import Hosts from "./views/Hosts.vue";
import WizardSD from "./views/WizardSD.vue";
import Certificates from "./views/Certificates.vue";
import InitializedSD from "./views/InitializeSD.vue";
import Terminal from "./views/Terminal.vue";
import Action from "./views/Action.vue";
import About from "./views/About.vue";

Vue.use(Router);

const router = new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    {
      path: "/",
      name: "home",
      component: Hosts,
    },
    {
      path: "/hosts",
      name: "hosts",
      component: Hosts,
    },
    {
      path: "/certificates",
      name: "certificates",
      component: Certificates,
    },
    {
      path: "/wizardsd",
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
      component: Terminal,
    },
    {
      path: "/hosts/:host/actions/:action",
      name: "setup",
      component: Action,
    },
    {
      path: "/about",
      name: "about",
      component: About,
    },
  ],
});

router.beforeEach((to, from, next) => {
  next();
});

export default router;
