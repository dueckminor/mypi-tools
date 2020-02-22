import Vue from "vue";
import Router from "vue-router";
import WizardSD from "./views/WizardSD.vue";
import InitializedSD from "./views/InitializeSD.vue";
import Terminal from "./views/Terminal.vue";
import About from "./views/About.vue";

Vue.use(Router);

const router = new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    {
      path: "/",
      name: "home",
      component: WizardSD
    },
    {
      path: "/wizardsd",
      name: "wizardsd",
      component: WizardSD
    },
    {
      path: "/initializesd",
      name: "initializesd",
      component: InitializedSD
    },
    {
      path: "/terminal",
      name: "terminal",
      component: Terminal
    },
    {
      path: "/about",
      name: "about",
      component: About
    },
  ]
});

router.beforeEach((to, from, next) => {
  next();
});

export default router;
