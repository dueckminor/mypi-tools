<template>
  <v-app light id="app">
    <app-navigation></app-navigation>
    <!--
    <v-app-bar app
      ><v-app-bar-nav-icon></v-app-bar-nav-icon
      ><v-toolbar-title>MYPI-Setup</v-toolbar-title></v-app-bar
    >
    -->
    <v-main transition="slide-x-transition" width="100%">
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script>
import AppNavigation from "@/components/AppNavigation";
import Vue from "vue";

export default {
  name: "app",
  components: {
    AppNavigation,
  },
  created() {
    // eslint-disable-next-line no-console
    console.log("app-created");
    Vue.prototype.$app = this;
  },
  mounted() {
    // eslint-disable-next-line no-console
    console.log("app-mounted");
  },
  data: function() {
    return {
      showLogin: true,
    };
  },
  methods: {},
};
</script>

<style></style>
